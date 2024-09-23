# Documentation of Development for Dynamic_Pod_Sec

## admission-controller/
│
├── main.go                         # Main entry point
├── webhook/
│   ├── webhook.go                  # Admission webhook setup and configuration
│
├── context_capabilities/           # Folder for security context and capabilities checks
│   ├── check_privilege.go          # Checks for privileged containers
│   ├── check_read_only_root.go     # Ensures read-only root filesystem
│   ├── check_run_as_user.go        # Verifies runAsUser and runAsGroup settings
│   ├── check_capabilities.go       # Ensures appropriate Linux capabilities
│   ├── validate_capability.go      # Validates dropped and allowed capabilities
│
├── volume_security/                # Folder for volume-related security checks
│   ├── check_hostpath.go           # Ensures no disallowed host paths are used
│
├── network_security/               # Folder for network-related security checks (new)
│   ├── check_network_policy.go     # Checks for existence and correctness of network policies
│   ├── check_host_network.go       # Ensures host network mode is not used
│
├── api_restrictions/               # Folder for API-related security checks (new)
│   ├── check_service_account.go    # Enforces service account policies
│   ├── check_api_access.go         # Restricts certain Kubernetes API access
│
├── resource_limits/                # Folder for resource quota and limit checks (new)
│   ├── check_resource_limits.go    # Ensures CPU and memory limits are set
│   ├── check_resource_requests.go  # Validates resource requests for containers
│
├── image_security/                 # Folder for image-related security checks (new)
│   ├── check_image_tags.go         # Ensures no unverified or 'latest' tags are used
│   ├── check_image_registry.go     # Blocks the use of certain image registries
│   ├── check_image_signing.go      # Enforces image signing and verification
│
└── configs/                        # Configuration files for the admission controller
    ├── webhook-config.yaml         # Configuration for the admission webhook
    ├── security-policies.yaml      # Custom security policies for the webhook




### `main.go` Recap:
- **Responsibility:** This is the entry point of the webhook service. It initializes the HTTP server and ties everything together. Here's what it does:
  - Logs the start of the webhook server.
  - Retrieves the TLS certificate and key from AWS Secrets Manager using the `getSecret` function.
  - Registers the `/validate` endpoint to handle incoming admission requests using the logic in the `admission` package.
  - Creates an HTTP server with TLS configuration from the `server` package.
  - Starts the server and gracefully handles any failure.

### **Admission Package (`admission/webhook.go` and `admission/validation.go`):**
- **Responsibility:** This package is the core logic of your admission webhook, where you implement the actual validation of Pods.
  - **`webhook.go`:** Contains the `HandleAdmissionRequest` function, which processes incoming AdmissionReview requests. This is where the webhook interacts with the Kubernetes API.
    - Decodes the AdmissionReview request.
    - Calls the `validatePod` function to determine whether the pod should be allowed.
    - Constructs an AdmissionReview response with the decision.
  - **`validation.go`:** Implements the specific validation logic in `validatePod`, such as checking for privileged containers, enforcing read-only root file systems, or other PodSecurityPolicy (PSP)-like checks.

### **Server Package (`server/server.go`):**
- **Responsibility:** This package is responsible for setting up the HTTP server with various middleware and configurations.
  - **TLS:** The server uses TLS for secure communication, which is mandatory for Kubernetes admission webhooks.
  - **Logging Middleware:** Logs incoming requests and their responses, which is useful for tracking and monitoring (e.g., sending logs to monitoring tools like Datadog or Prometheus).
  - **Health Check Endpoint:** Provides a `/healthz` endpoint to check if the webhook service is running and responsive.
  - **Graceful Shutdown:** Ensures the server can be gracefully stopped, handling ongoing requests before shutting down.

### **Key Accomplishments:**
1. **Webhook Server Initialization:** The server is initialized with proper logging, TLS, and health-check mechanisms.
2. **Admission Handling Logic:** The core admission control logic is set up to validate incoming pod requests based on security policies (Pod Security Policies, etc.).
3. **Runtime Pod Security Enforcement:** The webhook enforces runtime policies for Pod creation and modification, rejecting or allowing requests based on the validation results.
4. **Logging Middleware:** Provides a way to log all incoming requests and the corresponding responses for auditing and monitoring.
5. **Deployment Pipeline Preparedness:** We've discussed the necessary Kubernetes deployment YAML files (`deployment.yaml`, `service.yaml`, `ValidatingWebhookConfiguration.yaml`) to deploy the webhook service into a Kubernetes cluster.

### Next Steps:
1. **Expand Validation Logic:** You can further customize `validatePod` in `validation.go` to check more security properties (e.g., no host networking, disallowing certain image registries, etc.).
2. **Create Deployment Resources:** Write Kubernetes manifests for deploying the webhook service (`deployment.yaml`, `service.yaml`, `ValidatingWebhookConfiguration.yaml`).
3. **Testing:** Test the webhook in audit mode first (non-blocking) to validate that it correctly intercepts and logs requests. Then switch to enforce mode for runtime pod security.

### Architecture Overview:
- **main.go:** Initializes the server and ties the admission logic with the server.
- **admission/webhook.go:** Handles the HTTP logic for processing incoming admission requests.
- **admission/validation.go:** Contains the actual security logic for validating incoming pods.
- **server/server.go:** Sets up the server with TLS, logging, health checks, and graceful shutdown.

You're well-positioned to complete the webhook service and deploy it for runtime pod security enforcement in Kubernetes!

