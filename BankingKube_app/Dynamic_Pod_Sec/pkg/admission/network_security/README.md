## network_security/

### Purpose:
Ensures that network communication within the cluster adheres to defined security policies, preventing unauthorized access and promoting segmentation.

### Key Benefits:
- Enforces network segmentation and isolation policies.
- Validates network policies are in place and applied to critical workloads.
- Blocks pods from using insecure network configurations like hostNetwork.

### Checks Implemented:
- Verification of existing network policies for each namespace.
- Validation of pods’ network configurations, ensuring compliance.