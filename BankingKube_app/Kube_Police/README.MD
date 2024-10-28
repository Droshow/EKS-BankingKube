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
│   │   │   ├── check_privilege.go
│   │   │   ├── check_read_only_root.go
│   │   │   ├── check_run_as_user.go
│   │   │   └── validate_capability.go
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