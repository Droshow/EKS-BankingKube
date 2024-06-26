### Checking Module: Databases
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/databases

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
Success! The configuration is valid.

### Checking Module: Networking
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/networking

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```
Error: Reference to undeclared resource

  on alb.tf line 64, in resource "aws_lb_target_group" "eks_tg":
  64:   vpc_id   = aws_vpc.eks_vpc.id

A managed resource "aws_vpc" "eks_vpc" has not been declared in the root
module.
```
### Checking Module: Security
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/security

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```
Error: Reference to undeclared resource

  on acm.tf line 14, in resource "aws_acm_certificate_validation" "cert":
  14:   validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]

A managed resource "aws_route53_record" "cert_validation" has not been
declared in the root module.

Error: Reference to undeclared resource

  on main.tf line 44, in resource "aws_security_group" "eks_alb_sg":
  44:   vpc_id      = aws_vpc.eks_vpc.id

A managed resource "aws_vpc" "eks_vpc" has not been declared in the root
module.
```
### Checking Module: EKS
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/eks

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```
Error: Reference to undeclared module

  on main.tf line 3, in resource "aws_eks_cluster" "example":
   3:   role_arn = module.iam.role_arn

No module call named "iam" is declared in the root module.

Error: Reference to undeclared resource

  on outputs.tf line 15, in output "fargate_pod_execution_role_arn":
  15:   value = aws_iam_role.fargate_pod_execution_role.arn

A managed resource "aws_iam_role" "fargate_pod_execution_role" has not been
declared in the root module.
```

### Checking Module: Node Groups
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/node_groups

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.


``` Error: Reference to undeclared resource

  on main.tf line 5, in resource "aws_eks_fargate_profile" "default":
   5:   pod_execution_role_arn = aws_iam_role.fargate_pod_execution_role.arn

A managed resource "aws_iam_role" "fargate_pod_execution_role" has not been
declared in the root module.

```

### Checking Module: Storage
/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules/storage

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.53.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

``` Error: Reference to undeclared resource

  on main.tf line 11, in resource "aws_efs_mount_target" "efs_mount_target":
  11:   security_groups = [aws_security_group.efs_sg.id]

A managed resource "aws_security_group" "efs_sg" has not been declared in the
root module. ```
