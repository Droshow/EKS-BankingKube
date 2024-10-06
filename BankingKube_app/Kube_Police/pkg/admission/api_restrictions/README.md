api_restrictions/
Purpose:
This service enforces API and service access restrictions to ensure that pods and users operate within predefined boundaries, minimizing potential attack vectors.

Key Benefits:

Restricts sensitive Kubernetes API access paths.
Enforces service account policies to prevent unauthorized actions.
Blocks insecure service configurations that could expose sensitive data.
Checks Implemented:

Restricts access to specific API endpoints (e.g., /exec, /proxy).
Limits the usage of default or admin service accounts for running workloads.