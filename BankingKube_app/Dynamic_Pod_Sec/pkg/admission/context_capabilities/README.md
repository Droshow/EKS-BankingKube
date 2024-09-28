## context_capabilities/

### Purpose:
Manages security contexts and container capabilities, ensuring that workloads do not exceed their required privileges, adhering to security best practices.

### Key Benefits:
- Enforces security context policies (e.g., runAsNonRoot, runAsUser).
- Prevents the use of dangerous Linux capabilities such as CAP_SYS_ADMIN.
- Mitigates risks associated with elevated privileges in containers.

### Checks Implemented:
- Privileged container checks.
- Capability validation for add and drop settings.
- Security context configurations like read-only root filesystem.