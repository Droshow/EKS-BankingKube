# Dynamic Admission Webhook

## Directory Structure

```plaintext
dynamic-admission-webhook/
│
├── cmd/                     # Entry point
│   └── main.go              # Main application logic
│
├── pkg/
│   ├── admission/           # Core logic for admission control
│   │   ├── webhook.go       # Webhook server and handling logic
│   │   └── validation.go    # Pod security and validation logic
│   └── server/              # Server setup
│       └── server.go        # Web server setup
│
├── configs/                 # Webhook configurations
│   └── webhook.yaml         # Kubernetes ValidatingWebhookConfiguration
│
├── scripts/                 # Utility scripts for running or building
│   └── install.sh           # Script to install and apply configurations
│
├── Dockerfile               # Dockerfile for the webhook
├── go.mod                   # Go module
├── README.md                # Project documentation