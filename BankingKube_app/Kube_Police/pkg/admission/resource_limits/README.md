## resource_management/

### Purpose:
Monitors and enforces resource allocation and usage policies to prevent resource abuse and ensure optimal cluster performance.

### Key Benefits:
- Prevents unbounded resource usage that could lead to performance issues.
- Enforces CPU and memory limits to maintain cluster stability.
- Manages resource quotas to ensure fair resource allocation across namespaces.

### Checks Implemented:
- Validates that all pods have resource limits defined.
- Ensures compliance with resource quota policies for each namespace.