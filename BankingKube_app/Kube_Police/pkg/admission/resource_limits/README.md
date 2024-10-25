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

## Validation Differences Between `enforceResourceLimits` and `enforceResourceRequests`

### `enforceResourceLimits: true`

**Purpose**: Ensures that resource limits are defined for each container in a pod and that they fall within specified ranges.

**Validation Checks**:
- **Resource Limits Defined**: Checks if each container in the pod has resource limits defined.
- **CPU Limits Range**: Ensures that the CPU limits fall within the specified range (min and max).
- **Memory Limits Range**: Ensures that the memory limits fall within the specified range (min and max).

### `enforceResourceRequests: true`

**Purpose**: Ensures that resource requests are defined for each container in a pod.

**Validation Checks**:
- **Resource Requests Defined**: Checks if each container in the pod has resource requests defined.

### Summary of Differences

| Aspect            | `enforceResourceLimits`                                  | `enforceResourceRequests`                       |
|-------------------|----------------------------------------------------------|-------------------------------------------------|
| **Purpose**       | Ensures resource limits are defined and within specified ranges | Ensures resource requests are defined           |
| **Validation Checks** | - Resource limits defined<br>- CPU limits within range<br>- Memory limits within range | - Resource requests defined                     |