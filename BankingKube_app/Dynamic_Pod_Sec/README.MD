# Dynamic Admission Webhook

## Directory Structure

**KubePolicy** is a comprehensive security solution for Kubernetes environments. The tool will scan, audit, and provide automated remediation for admissions, RBAC configurations, and network policies, ensuring best practices and compliance with security frameworks such as PSD2, PCI-DSS, and more.

By leveraging tools like **Kyverno** for policy enforcement, the system will not only detect misconfigurations and vulnerabilities but also apply corrective actions to ensure clusters remain secure and compliant.

## TODO
28.10.2024 Admissions finished 
-  EKS-BankingKube/BankingKube_app/Kube_Police/pkg/admission/image_security/check_image_signing.go - func isImageSigned()


Run admissions in Kubernetes Cluster (preferrably customer one) and provide audit , remediation possilibity through guidance through Kyverno. 

## TASK
Refine admissions code to use security-policies.yaml  

```plaintext
KubePolice/
│
├── cmd/                     # Entry point
│   ├── main.go              # Main application logic
│   └── secrets.go           # Secrets management logic
│
├── configs/                 # Webhook configurations
│   ├── security-policies.yaml # Security policies configuration
│   └── webhook-config.yaml  # Kubernetes ValidatingWebhookConfiguration
│
├── pkg/
│   ├── admission/           # Core logic for admission control
│   │   ├── api_restrictions/ # API access restrictions
│   │   │   ├── README.md
│   │   │   ├── check_api_access.go
│   │   │   └── check_service_account.go
│   │   ├── context_capabilities/ # Context and capabilities checks
│   │   │   ├── README.md
│   │   │   ├── check_capabilities.go
│   │   │   ├── check_pod_security_context.go
│   │   ├── image_security/  # Image security checks
│   │   │   ├── README.md
│   │   │   ├── check_image_registry.go
│   │   │   ├── check_image_signing.go
│   │   │   └── check_image_tags.go
│   │   ├── network_security/ # Network security checks
│   │   │   ├── README.md
│   │   │   ├── check_host_network.go
│   │   │   └── check_network_policy.go
│   │   ├── resource_limits/ # Resource limits checks
│   │   │   ├── README.md
│   │   │   ├── check_resource_limits.go
│   │   │   └── check_resource_requests.go
│   │   ├── volume_security/ # Volume security checks
│   │   │   ├── README.md
│   │   │   └── check_hostpath.go
│   │   └── webhook.go       # Webhook server and handling logic
│   ├── diary.md             # Development diary
│   └── server/              # Server setup
│       └── server.go        # Web server setup
│
├── scripts/                 # Utility scripts for running or building
│   └── install.sh           # Script to install and apply configurations
│
├── Dockerfile               # Dockerfile for the webhook
├── go.mod                   # Go module
├── go.sum                   # Go dependencies
├── README.md                # Project documentation

## 5.11.2024

## Immediate Development Focus

### Phase 1: Complete Admission Control Core

#### Finish Key Admission Checks:
- Ensure each category (RBAC, network security, context capabilities, etc.) has robust admission checks.
- Standardize each check to include appropriate logging and error handling to capture critical security events.

#### Implement Centralized Policy Configuration:
- Finalize `security-policies.yaml` to house all necessary policy definitions (RBAC rules, network policies, capabilities, etc.).
- Load policies dynamically from `security-policies.yaml` in each relevant admission check for consistency and ease of customization.

#### Audit-Ready Logging Enhancement:
- Refine logging mechanisms across all admission checks.
- Include details like namespace, pod name, reason for check failure, and the user initiating the action.

### Phase 2: Scanning and Reporting

#### Develop a Scanning Module:
- Implement a periodic scanning process to validate that pods and namespaces comply with configured security policies.
- Ensure that each scan checks for non-compliant RBAC roles, missing or misconfigured network policies, and excessive permissions in context capabilities.

#### Generate Compliance Reports:
- Build basic reporting capabilities to summarize scan results, such as non-compliant resources and the types of violations found.
- Include both a high-level summary and details to help administrators address specific issues.

### Phase 3: Remediation Strategy

#### Design Remediation Actions:
- Start with basic remediation suggestions (e.g., “remove privileged container,” “apply network policy”) as part of the compliance report.
- Define automated remediation logic for issues that can be handled safely, like enforcing `readOnlyRootFilesystem` or applying default-deny network policies.

#### Integrate with Kyverno for Extended Policy Enforcement:
- Integrate Kyverno as a remediation option, allowing it to enforce policies dynamically on detected violations.
- Set up templates for Kyverno policies that users can apply as needed.

#### Develop a Dashboard for Visibility:
- Create a simple dashboard to view audit logs, compliance reports, and active policy violations.
- Include filters for different namespaces, resource types, and violation types, making it easier to spot critical issues.

## Development Milestones

### Admission Control:
- Finalize all key checks and ensure each can dynamically read from `security-policies.yaml`.
- Test admission controls against various scenarios to confirm correct operation.

### Scanning & Reporting:
- Deploy a functional scanning module that runs in the background, periodically validating compliance.
- Implement initial reporting capabilities to summarize scan findings and detail any policy violations.

### Remediation:
- Define remediation strategies and begin integrating Kyverno for automated enforcement options.
- Test automated remediation actions in a controlled environment to avoid disruptions.

## Immediate Deliverables
- Completed admission checks and logging.
- A functional scanner that generates a compliance report.
- Basic dashboard layout and integration with logs and compliance report data.

This plan gives a clear path to achieving a robust Kubernetes admission control and policy auditing solution, with a foundation for future remediation and advanced monitoring.

## 2.12.2024

Using Kubernetes for containerized applications offers great flexibility and scalability, but it also brings complexity and new security challenges. For a company, an automated compliance and hardening framework—like the combination of MutatingWebhookConfiguration, ValidatingWebhookConfiguration, and Kyverno—adds significant business value by making sure that security, governance, and best practices are applied consistently and automatically. Here’s why this approach matters from a business perspective:

1. **Enhanced Security with Minimal Overhead**
   - MutatingWebhookConfiguration instantly adjusts any new resource to meet essential security baselines, such as requiring containers to run without root privileges or setting default resource limits to prevent accidental overuse.
   - This reduces the chances of misconfigured applications, avoiding risks before they even reach production.
   - Automatic enforcement keeps workloads secure without manual intervention, reducing the chance of human error—critical for data protection and customer trust.

2. **Regulatory Compliance Made Easy**
   - ValidatingWebhookConfiguration performs real-time checks to ensure that each application complies with security and compliance standards before it is admitted into the environment.
   - It can block non-compliant resources from deploying, helping the company meet industry standards and regulatory requirements like PCI, SOC 2, or HIPAA without relying on manual compliance checks.
   - With Kyverno handling periodic post-deployment scans, the company can continuously validate that workloads remain compliant, minimizing audit preparation time and demonstrating due diligence to regulators.

3. **Reduced Operational Costs and Risks**
   - By preventing security misconfigurations from entering production, this automated framework minimizes the potential for costly incidents, such as breaches or downtime due to resource conflicts.
   - Since Kyverno remediates non-compliant resources automatically, the operations team can focus on innovation and proactive improvements instead of manually enforcing policies and fixing configuration issues.
   - This approach also allows a smaller operations team to manage security and compliance effectively, even as the environment scales.

4. **Increased Agility with Built-in Governance**
   - Developers get the flexibility to deploy applications while meeting security policies by default, without needing constant input from security teams.
   - The framework enforces guardrails that support DevOps speed without sacrificing compliance, allowing the team to iterate faster and innovate safely.
   - For multi-tenant environments or large teams, automated policy enforcement ensures consistency, so teams can onboard applications without concerns over inconsistent configurations or accidental missteps.

5. **Future-Proofing the Kubernetes Environment**
   - As the business scales, an automated compliance and security framework supports growth by reducing the need for reactive changes or emergency fixes.
   - With dynamic, automated policies, the company can adapt to new security requirements, regulatory demands, or operational best practices by adjusting the policy configuration—not the code.
   - Kyverno’s ability to detect and remediate policy drift keeps the environment in line with changing policies, helping the company stay ahead of potential compliance issues.

   ### more granular:

   2. **Add CA Bundle to Webhook Configuration**
   - Ensure the `<CA_BUNDLE>` placeholder is replaced with a valid base64-encoded CA bundle of the admission controller's certificate.
   - Automate this in your deployment workflow if possible.

3. **Deploy the Admission Controller**
   - Build and deploy the admission-controller-service as a pod in your cluster. Ensure it:
     - Exposes the `/validate/context`, `/validate/volumes`, and `/validate/network` endpoints.
     - Has proper RBAC permissions to access the necessary Kubernetes API objects for validation.
   - Test it locally before deploying to production.

4. **Add MutatingWebhookConfiguration (Optional)**
   - Introduce a MutatingWebhookConfiguration to handle defaulting policies. For example:
     - Add default labels or annotations.
     - Enforce resource requests/limits if not specified.
     - Add `runAsNonRoot` or other security defaults.

5. **Test ValidatingWebhook Endpoints**
   - Use `kubectl` to create pods or other resources to ensure each validation logic is triggered correctly.
   - Simulate failure cases to validate the rejection messages.

6. **Write Unit and Integration Tests**
   - Ensure all admission checks are covered with unit tests.
   - Write integration tests to validate the webhook behavior in a real cluster environment.

7. **Enhance Logging and Observability**
   - Add detailed logs for each validation step.
   - Integrate with tools like Prometheus or Datadog for monitoring webhook performance.