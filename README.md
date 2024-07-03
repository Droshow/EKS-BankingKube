# EKS-BankingKube - All Infra Development
EKS and other resources for Open Banking resources

## To-Do / DONE [13.6.2024]
- [ ] **Ingress**: Set up an Ingress controller, such as NGINX or ALB Ingress Controller, to expose services to the outside world. ALB is set and IAM role by terraform, but Ingress needs to be provisioned via K8 Manifest
- [ ] **Storage**: Set up storage solutions, such as EBS volumes or EFS file systems, for applications that need to persist data or share data between pods.
- [ ] **Logging & Monitoring**: Set up logging and monitoring for observability. This could involve setting up CloudWatch logs for the cluster and integrating with a monitoring solution like Prometheus and Grafana.
- [ ] **Autoscaling**: Set up the Kubernetes Cluster Autoscaler and the Kubernetes Metrics Server to automatically scale the cluster based on load. Not applicable at EKS Fargate however we will set Kubernetes metrics server kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
## Functionally DONE on 13.6.2024 with validations for modules

## To-Do / Write TerraTest for Modules
- [ ] databases
- [ X] networking
- [ ] eks
- [ ] node_groups
- [ X] security
- [ ] storage


## Benefits of Using TerraTest for Infrastructure Testing

- **Automated Testing**: TerraTest allows you to write automated tests for your infrastructure code. This can help catch issues early before they affect your production environment.

- **Validation**: TerraTest can validate that your infrastructure works as expected. For example, it can verify that a server is up and running, that a database is accessible, or that a load balancer is distributing traffic correctly.

- **Refactoring and Evolving Code**: As your infrastructure evolves, TerraTest can ensure that changes or refactoring of your infrastructure code do not break existing functionality.

- **Consistency**: By testing your modules, you can ensure that they behave consistently across different environments. This can be particularly useful if you use the same modules to create infrastructure in different environments (e.g., dev, test, prod).

- **Confidence**: Having a suite of tests can give you confidence that your infrastructure is working correctly. This can be particularly important for critical infrastructure that your application relies on.

## To-Do / DONE [28.6.2024]

Did initial TerraTest for Networking and Security Modules without ACM certificate

[3.7.2024] Next the task to play around with packages - And have some pause from infra

# BankingKube API Development Guide - App Development

## Overview
This guide outlines the steps for developing the BankingKube API, which facilitates open banking transactions. The API will support operations such as initiating payments, retrieving transaction status, and managing user consent.

## Step 1: Define API Specifications

### Objective
Establish the functional requirements and endpoints of your API.

### Tasks
- **Identify Key Operations**: Determine necessary operations like initiating payments, retrieving transaction status, and managing user consent.
- **Design RESTful Endpoints**:
  - `POST /payments` - Initiate a payment.
  - `GET /payments/{id}` - Check payment status.
  - `POST /consents` - Manage user consents.

## Step 2: Set Up the Golang Environment

### Objective
Prepare your development environment for Golang.

### Tasks
- **Install Go**: Ensure the latest version of Go is installed on your development machine.
- **Choose an IDE**: Select an IDE or editor that supports Go, such as Visual Studio Code or GoLand.
- **Set Up Version Control**: Initialize a Git repository to manage version control.

## Step 3: Create a Basic HTTP Server

### Objective
Build a foundational HTTP server in Golang to serve as the backbone for your API.

### Tasks
- **Use `net/http` Package**: Start by creating a simple server using Golang's built-in `net/http` package.
- **Routing**: Implement basic routing to handle requests to your defined endpoints.

## Step 4: Integrate with Open Banking APIs

### Objective
Develop functionality to communicate with bank APIs.

### Tasks
- **Authentication**: Implement the authentication mechanism required by the bank’s API (e.g., OAuth 2.0).
- **API Calls**: Use `net/http` or `github.com/go-resty/resty/v2` for making HTTP calls to the bank APIs.
- **Error Handling**: Handle API responses robustly, focusing on error responses from the bank APIs.

## Step 5: Implement Logging and Error Handling

### Objective
Ensure robust logging and error handling for troubleshooting and compliance.

### Tasks
- **Logging**: Use a logging package such as `logrus` or `zap` for structured logging.
- **Error Handling**: Implement comprehensive error handling across your API.

## Step 6: Testing

### Objective
Write tests to ensure your API functions as expected.

### Tasks
- **Unit Tests**: Write unit tests for individual functions, particularly those interfacing with bank APIs.
- **Integration Tests**: Develop integration tests that run against a test version of the bank APIs (if available).

## Conclusion
Begin with setting up a simple server and gradually integrate more complex functionalities like communicating with bank APIs, handling authentication, and robust error handling. Maintain a focus on writing clean, maintainable code and establishing a solid foundation with good testing practices from the start.

