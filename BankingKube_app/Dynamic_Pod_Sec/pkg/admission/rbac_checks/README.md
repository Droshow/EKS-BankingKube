## Description of Each File

### `check_role_binding.go`

**Purpose**: Verifies that role bindings follow security best practices. Ensures that sensitive roles (e.g., admin, edit) are not unnecessarily granted to service accounts or users that don’t need them.

**Functionality**: Checks that role bindings only apply to necessary namespaces and that role assignments are minimal for least privilege.

### `check_cluster_role_binding.go`

**Purpose**: Ensures that ClusterRoleBindings are appropriately scoped and not over-permissive.

**Functionality**: Validates that critical roles with cluster-wide access are limited to specific service accounts or user groups. Logs findings if overly permissive cluster roles are detected.

### `check_role_scope.go`

**Purpose**: Ensures roles have the correct scope based on their intended use. This enforces a clear distinction between namespace-scoped roles and cluster-scoped roles.

**Functionality**: Validates that roles are scoped correctly. Logs a warning if a namespace-scoped role is bound at the cluster level, which could signal misconfiguration.

### `check_permission_levels.go`

**Purpose**: Validates permission levels assigned within roles to prevent unnecessary escalations, such as overly permissive rights (e.g., * access on resources).

**Functionality**: Parses role permissions and flags any that grant broad access to sensitive resources (e.g., secrets, configmaps). Ensures read-only where appropriate.