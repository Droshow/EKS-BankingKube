## image_security/

### Purpose:
Validates container image sources and configurations to ensure the security and integrity of container images running in the cluster.

### Key Benefits:
- Prevents the use of unverified or potentially malicious container images.
- Enforces image signing and verification policies.
- Blocks the use of the latest or untagged images that can lead to unpredictability.

### Checks Implemented:
- Validates image tags against allowed and disallowed registries.
- Ensures images are signed and verified before deployment.