# Dynamic Admission Webhook

# README for Each Admission Service Folder

## 1. api_restrictions/

**Purpose**: This service enforces API and service access restrictions to ensure that pods and users operate within predefined boundaries, minimizing potential attack vectors.

**Key Benefits**:
- Restricts sensitive Kubernetes API access paths.
- Enforces service account policies to prevent unauthorized actions.
- Blocks insecure service configurations that could expose sensitive data.

**Checks Implemented**:
- Restricts access to specific API endpoints (e.g., /exec, /proxy).
- Limits the usage of default or admin service accounts for running workloads.

## 2. context_capabilities/

**Purpose**: Manages security contexts and container capabilities, ensuring that workloads do not exceed their required privileges, adhering to security best practices.

**Key Benefits**:
- Enforces security context policies (e.g., runAsNonRoot, runAsUser).
- Prevents the use of dangerous Linux capabilities such as CAP_SYS_ADMIN.
- Mitigates risks associated with elevated privileges in containers.

**Checks Implemented**:
- Privileged container checks.
- Capability validation for add and drop settings.
- Security context configurations like read-only root filesystem.

## 3. image_security/

**Purpose**: Validates container image sources and configurations to ensure the security and integrity of container images running in the cluster.

**Key Benefits**:
- Prevents the use of unverified or potentially malicious container images.
- Enforces image signing and verification policies.
- Blocks the use of the latest or untagged images that can lead to unpredictability.

**Checks Implemented**:
- Validates image tags against allowed and disallowed registries.
- Ensures images are signed and verified before deployment.

## 4. network_security/

**Purpose**: Ensures that network communication within the cluster adheres to defined security policies, preventing unauthorized access and promoting segmentation.

**Key Benefits**:
- Enforces network segmentation and isolation policies.
- Validates network policies are in place and applied to critical workloads.
- Blocks pods from using insecure network configurations like hostNetwork.

**Checks Implemented**:
- Verification of existing network policies for each namespace.
- Validation of pods’ network configurations, ensuring compliance.

## 5. resource_limits/

**Purpose**: Monitors and enforces resource allocation and usage policies to prevent resource abuse and ensure optimal cluster performance.

**Key Benefits**:
- Prevents unbounded resource usage that could lead to performance issues.
- Enforces CPU and memory limits to maintain cluster stability.
- Manages resource quotas to ensure fair resource allocation across namespaces.

**Checks Implemented**:
- Validates that all pods have resource limits defined.
- Ensures compliance with resource quota policies for each namespace.

## 6. volume_security/

**Purpose**: Manages the use of persistent volumes and storage configurations to prevent unauthorized access to sensitive host paths and secure storage configurations.

**Key Benefits**:
- Prevents the use of disallowed host paths like /var/run/docker.sock.
- Enforces secure storage class policies.
- Restricts the use of insecure volume types (e.g., emptyDir).

**Checks Implemented**:
- Validates volume configurations against defined security policies.
- Ensures sensitive paths and storage classes are not used improperly.



## Directory Structure

```plaintext
Dynamic_Pod_Sec/
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

```
├── pkg/
│   ├── admission/           # Core logic for admission control
│   │   ├── api_restrictions/ # API access restrictions
│   │   │   ├── README.md
│   │   │   ├── check_api_access.go          # apiRestrictions: restrictedAPIPaths
│   │   │   └── check_service_account.go     # serviceAccountRestrictions: restrictedServiceAccounts
│   │   ├── context_capabilities/ # Context and capabilities checks
│   │   │   ├── README.md
│   │   │   ├── check_capabilities.go        # capabilities: disallowedCapabilities
│   │   │   ├── check_privilege.go           # podSecurityContext: allowPrivilegeEscalation, runAsNonRoot
│   │   │   ├── check_read_only_root.go      # podSecurityContext: readOnlyRootFilesystem
│   │   │   ├── check_run_as_user.go         # (Potentially redundant, covered by check_privilege.go)
│   │   │   └── validate_capability.go       # (Potentially redundant, covered by check_capabilities.go)
│   │   ├── image_security/  # Image security checks
│   │   │   ├── README.md
│   │   │   ├── check_image_registry.go      # (Not directly mapped, but important for image security)
│   │   │   ├── check_image_signing.go       # (Not directly mapped, but important for image security)
│   │   │   └── check_image_tags.go          # (Not directly mapped, but important for image security)
│   │   ├── network_security/ # Network security checks
│   │   │   ├── README.md
│   │   │   ├── check_host_network.go        # networkSecurity: requiredNetworkPolicies
│   │   │   └── check_network_policy.go      # networkSecurity: requiredNetworkPolicies
│   │   ├── resource_limits/ # Resource limits checks
│   │   │   ├── README.md
│   │   │   ├── check_resource_limits.go     # resourceLimits: cpuLimits, memoryLimits, enforceResourceLimits
│   │   │   └── check_resource_requests.go   # (Potentially redundant, covered by check_resource_limits.go)
│   │   ├── volume_security/ # Volume security checks
│   │   │   ├── README.md
│   │   │   └── check_hostpath.go            # volumeSecurity: disallowedHostPaths