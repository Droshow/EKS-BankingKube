#!/bin/bash

set -e

CONFIG_DIR="../k8s"

echo "Deploying Kubernetes resources..."

kubectl apply -f ${CONFIG_DIR}/rbac.yaml
kubectl apply -f ${CONFIG_DIR}/webhook-service.yaml
kubectl apply -f ${CONFIG_DIR}/webhook-deployment.yaml
kubectl apply -f ${CONFIG_DIR}/mutation-webhook-config.yaml
kubectl apply -f ${CONFIG_DIR}/validation-webhook-config.yaml

echo "Deployment complete!"
