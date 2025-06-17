package network_security

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	egressTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	egressMeter   = otel.Meter("bankingkube/dynamicpodsec")
	egressDenied  metric.Int64Counter
	egressAllowed metric.Int64Counter
)

func init() {
	var err error
	egressDenied, err = egressMeter.Int64Counter("network_egress.denied")
	if err != nil {
		log.Println("Failed to create metric: network_egress.denied")
	}
	egressAllowed, err = egressMeter.Int64Counter("network_egress.allowed")
	if err != nil {
		log.Println("Failed to create metric: network_egress.allowed")
	}
}

// EgressPolicy defines a structure for egress network policies
type EgressPolicy struct {
	AllowedEgressCIDRs []string `yaml:"allowedEgressCIDRs"`
}

// SecurityPoliciesEgress represents the structure of the security-policies.yaml file
type SecurityPolicies struct {
	NetworkSecurity struct {
		EgressPolicy EgressPolicy `yaml:"egressPolicy"`
	} `yaml:"NetworkSecurity"`
}

// CheckEgress validates if a pod has correct egress restrictions
func CheckEgress(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := egressTracer.Start(ctx, "CheckEgress", trace.WithAttributes(
		attribute.String("operation", string(request.Operation)),
		attribute.String("resource", request.Resource.Resource),
	))
	defer span.End()

	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_parse_pod"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	span.SetAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
	)

	egressPolicy, err := getEgressPolicy()
	if err != nil {
		log.Println("Failed to load egress policy:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policy"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check pod annotations for egress IPs
	egressIPs, ok := pod.Annotations["egressIPs"]
	if !ok {
		log.Printf("Pod %s in namespace %s does not have egress IPs specified\n", pod.Name, pod.Namespace)

		egressDenied.Add(ctx, 1, metric.WithAttributes(
			attribute.String("pod", pod.Name),
			attribute.String("namespace", pod.Namespace),
			attribute.String("denial_reason", "missing_egress_ips"),
		))

		span.SetAttributes(
			attribute.String("result", "denied"),
			attribute.String("denial_reason", "missing_egress_ips"),
		)

		return false
	}

	span.SetAttributes(attribute.String("egress_ips", egressIPs))

	// Verify each egress IP against all allowed CIDRs
	ipList := strings.Split(egressIPs, ",")
	for _, ip := range ipList {
		ip = strings.TrimSpace(ip)

		subSpan := trace.SpanFromContext(ctx)
		subSpan.AddEvent("Checking IP", trace.WithAttributes(
			attribute.String("ip", ip),
		))

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			log.Printf("Invalid IP format: %s\n", ip)

			egressDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("ip", ip),
				attribute.String("denial_reason", "invalid_ip_format"),
			))

			span.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("ip", ip),
				attribute.String("denial_reason", "invalid_ip_format"),
			)

			return false
		}

		allowed := false
		for _, egressCIDR := range egressPolicy.AllowedEgressCIDRs {
			_, allowedCIDR, err := net.ParseCIDR(egressCIDR)
			if err != nil {
				log.Printf("Invalid CIDR format: %s\n", egressCIDR)

				egressDenied.Add(ctx, 1, metric.WithAttributes(
					attribute.String("pod", pod.Name),
					attribute.String("namespace", pod.Namespace),
					attribute.String("cidr", egressCIDR),
					attribute.String("denial_reason", "invalid_cidr_format"),
				))

				span.SetAttributes(
					attribute.String("result", "denied"),
					attribute.String("cidr", egressCIDR),
					attribute.String("denial_reason", "invalid_cidr_format"),
				)
				span.RecordError(err)

				return false
			}

			if allowedCIDR.Contains(parsedIP) {
				allowed = true
				subSpan.AddEvent("IP allowed by CIDR", trace.WithAttributes(
					attribute.String("ip", ip),
					attribute.String("cidr", egressCIDR),
				))
				break
			}
		}

		if !allowed {
			log.Printf("IP %s is not within any allowed CIDR range\n", ip)

			egressDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("ip", ip),
				attribute.String("denial_reason", "ip_not_in_allowed_cidrs"),
			))

			span.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("ip", ip),
				attribute.String("denial_reason", "ip_not_in_allowed_cidrs"),
			)

			return false
		}
	}

	// All egress IPs are allowed
	egressAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
		attribute.String("egress_ips", egressIPs),
		attribute.Int("ip_count", len(ipList)),
	))

	span.SetAttributes(
		attribute.String("result", "allowed"),
		attribute.Int("ip_count", len(ipList)),
	)

	return true
}

// getEgressPolicy loads the egress policy from the configuration file
func getEgressPolicy() (*EgressPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPolicies
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.NetworkSecurity.EgressPolicy, nil
}
