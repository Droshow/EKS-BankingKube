# EKS-BankingKube
EKS and other resources for Open Banking resources

## To-Do / DONE [13.6.2024]
- [ ] **Ingress**: Set up an Ingress controller, such as NGINX or ALB Ingress Controller, to expose services to the outside world. ALB is set and IAM role by terraform, but Ingress needs to be provisioned via K8 Manifest
- [ ] **Storage**: Set up storage solutions, such as EBS volumes or EFS file systems, for applications that need to persist data or share data between pods.
- [ ] **Logging & Monitoring**: Set up logging and monitoring for observability. This could involve setting up CloudWatch logs for the cluster and integrating with a monitoring solution like Prometheus and Grafana.
- [ ] **Autoscaling**: Set up the Kubernetes Cluster Autoscaler and the Kubernetes Metrics Server to automatically scale the cluster based on load. Not applicable at EKS Fargate however we will set Kubernetes metrics server kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
## Functionally DONE on 13.6 with validations for modules

## To-Do / Write TerraTest for Modules
- [ ] databases
- [ ] networking
- [ ] eks
- [ ] node_groups
- [ ] security
- [ ] storage