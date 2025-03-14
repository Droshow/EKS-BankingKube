sasafa## CICD Runner or different access for EKS in private subnet to be able to run kubectl
## CICD Runner or Different Access for EKS in Private Subnet to be Able to Run kubectl
<!-- K8s resources in path: /Users/martin.drotar/Devsbridge/EKS-BankingKube/EKS_infra/modules/eks/main.tf -->
<!-- uncomment provider kubernetes in providers.tf -->
<!-- !!! You need a runner in K8 cluster private subnet-->


## Deployment knowledge

<!-- Separate deployment one using public runner that builds just EC2 instance and Networking module perhaps
And the other that uses self-hosted runner on that instance building everything else --> Done

<!-- Run the EKS AWS Provider Terraform 
Uncomment Kubernetes Provider Terraform
hanle the DB secret in EKS_infra/secrets.tf
Watch for certificate if exists in AWS or not
if not fetch_existing_cert in main.tf = false, then deploy will fail, change to true and deploy again
first run if not existing that switch to fetch_certificate = true
Setting up own runner don't forget to change the token in 
https://github.com/Droshow/EKS-BankingKube/settings/actions/runners/new?arch=x64&os=linux 
and consecutively in secrets


Key-Pair aws command
martin.drotar@CVX-1065 modules % aws ec2 create-key-pair --key-name ssh-key-bankingKube \
    --region eu-central-1 \
    --query 'KeyMaterial' --output text > ssh-key-bankingKube.pem  -->



### Tasks to be Done

0. **Fix self-hosted runner infra workflow** DONE
 

1. **Update Kubeconfig:** DONE
   - Run the `aws eks update-kubeconfig` command to generate the kubeconfig file that allows `kubectl` to interact with your EKS cluster.

3. **Automate aws-auth ConfigMap Update:** DONE
   - Automate the update of the `aws-auth` ConfigMap through the CI/CD workflow to ensure any new IAM roles or changes to existing roles are automatically reflected.

4. **Verify Network Connectivity:** DONE
   - Ensure that the EC2 instance has network access to the EKS API server, including proper security group rules and VPC endpoints if necessary.

5. **Test kubectl Access:** DONE
   - Verify that the EC2 instance can successfully run `kubectl` commands against the EKS cluster.
DONE

6. **Modify deployment workflow - deyploy-dynamic-pod-sec for self runner thing:**
7. **Generate TLS Certificates:** DONE but some remedies need to be done
   - Run the `generate-certs.sh` script to generate TLS certificates and create the Kubernetes secret.

8. **Deploy Webhook Manifests:** DONE
   - Run the `deploy-webhook.sh` script to deploy the webhook manifests to the EKS cluster.

9. **Document the Process:** DONE
   - Document the steps and configurations required to set up and maintain access for future reference and development iterations.

10. **Repaird destroy job that is totallly invalid:** 