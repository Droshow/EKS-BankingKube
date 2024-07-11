# EKS-BankingKube - All Infra Development
EKS and other resources for Open Banking resources

## To-Do / DONE [11.7.2024]
- [ ] **connection**:

- [ ] **Ingress**: Set up an Ingress controller, such as NGINX or ALB Ingress Controller, to expose services to the outside world. ALB is set and IAM role by terraform, but Ingress needs to be provisioned via K8 Manifest
- [ ] **Storage**: Set up storage solutions, such as EBS volumes or EFS file systems, for applications that need to persist data or share data between pods.
- [ ] **Logging & Monitoring**: Set up logging and monitoring for observability. This could involve setting up CloudWatch logs for the cluster and integrating with a monitoring solution like Prometheus and Grafana.
- [ ] **Autoscaling**: Set up the Kubernetes Cluster Autoscaler and the Kubernetes Metrics Server to automatically scale the cluster based on load. Not applicable at EKS Fargate however we will set Kubernetes metrics server kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
## Functionally DONE on 13.6.2024 with validations for modules

## To-Do / Provision Modules
- [ x] databases
- [ X] networking
- [ x] eks
- [ x] node_groups
- [ X] security
- [ x] storage


## Benefits of Using TerraTest for Infrastructure Testing

- **Automated Testing**: TerraTest allows you to write automated tests for your infrastructure code. This can help catch issues early before they affect your production environment.

- **Validation**: TerraTest can validate that your infrastructure works as expected. For example, it can verify that a server is up and running, that a database is accessible, or that a load balancer is distributing traffic correctly.

- **Refactoring and Evolving Code**: As your infrastructure evolves, TerraTest can ensure that changes or refactoring of your infrastructure code do not break existing functionality.

- **Consistency**: By testing your modules, you can ensure that they behave consistently across different environments. This can be particularly useful if you use the same modules to create infrastructure in different environments (e.g., dev, test, prod).

- **Confidence**: Having a suite of tests can give you confidence that your infrastructure is working correctly. This can be particularly important for critical infrastructure that your application relies on.

## To-Do / DONE [28.6.2024]

Did initial TerraTest for Networking and Security Modules without ACM certificate

[3.7.2024] Next the task to play around with packages - And have some pause from infra
