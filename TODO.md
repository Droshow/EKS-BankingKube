sasafa## CICD Runner or different access for EKS in private subnet to be able to run kubectl
## CICD Runner or Different Access for EKS in Private Subnet to be Able to Run kubectl
<!-- K8s resources in path: /Users/martin.drotar/Devsbridge/EKS-BankingKube/EKS_infra/modules/eks/main.tf -->
<!-- uncomment provider kubernetes in providers.tf -->
<!-- !!! You need a runner in K8 cluster private subnet-->


## Deployment knowledge

<!-- Run the EKS AWS Provider Terraform 
Uncomment Kubernetes Provider Terraform
hanle the DB secret in EKS_infra/secrets.tf
Watch for certificate if exists in AWS or not
Setting up own runner don't forget to change the token in 
https://github.com/Droshow/EKS-BankingKube/settings/actions/runners/new?arch=x64&os=linux 
and consecutively in secrets

Key-Pair aws command
martin.drotar@CVX-1065 modules % aws ec2 create-key-pair --key-name ssh-key-bankingKube \
    --region eu-central-1 \
    --query 'KeyMaterial' --output text > ssh-key-bankingKube.pem  -->


### Tasks to be Done

0. **Fix self-hosted runner infra workflow**
 

1. **Update Kubeconfig:**
   - Run the `aws eks update-kubeconfig` command to generate the kubeconfig file that allows `kubectl` to interact with your EKS cluster.

3. **Automate aws-auth ConfigMap Update:**
   - Automate the update of the `aws-auth` ConfigMap through the CI/CD workflow to ensure any new IAM roles or changes to existing roles are automatically reflected.

4. **Verify Network Connectivity:**
   - Ensure that the EC2 instance has network access to the EKS API server, including proper security group rules and VPC endpoints if necessary.

5. **Test kubectl Access:**
   - Verify that the EC2 instance can successfully run `kubectl` commands against the EKS cluster.

6. **Generate TLS Certificates:**
   - Run the `generate-certs.sh` script to generate TLS certificates and create the Kubernetes secret.

7. **Deploy Webhook Manifests:**
   - Run the `deploy-webhook.sh` script to deploy the webhook manifests to the EKS cluster.

8. **Document the Process:**
   - Document the steps and configurations required to set up and maintain access for future reference and development iterations.