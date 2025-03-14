# BankingKube: Secure and Monitored Kubernetes Fintech Platform

# Project Overview
The project is a state-of-the-art Kubernetes microservices platform designed specifically for Fintech organizations. It provides a secure, monitored Kubernetes environment with built-in support for integrated payment gateways—whether card-based or cardless (open banking).

## Key Features
- **High-grade security** with comprehensive Kubernetes monitoring and runtime protection.
- **Integrated payment processing APIs** for seamless card and open banking transactions.
- **Scalable microservices architecture** tailored to handle financial workflows efficiently.
- **Compliance-ready infrastructure** with built-in tools for PCI-DSS, PSD2, and other Fintech regulations.

This boilerplate solution is ideal for Fintech organizations looking to deploy a secure, scalable, and fully monitored Kubernetes environment that meets industry standards while simplifying payment integrations.

## Deployment
1. Change the DB secret in `EKS_infra/secrets.tf`.
2. In `main.tf`, set `fetch_cert = false`.
3. Run the public runner workflow: `.github/workflows/deploy-infra-public-runner.yml`.
4. The job will fail. Run it again with `fetch_cert = true`.
5. Set the configs for the runner. You need to create a new one [here](https://github.com/Droshow/EKS-BankingKube/settings/actions/runners/new) and put it into AWS secrets for the runner.
6. The runner may not succeed. You might need to fiddle around with EC2 and run commands manually.
7. Run the private runner workflow and uncomment the K8s stuff in `EKS_infra/modules/eks/main.tf`.