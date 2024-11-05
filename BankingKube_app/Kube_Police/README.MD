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