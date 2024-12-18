#!/bin/bash

set -e

NAMESPACE="default"
SERVICE="admission-controller-service"
SECRET="webhook-tls-secret"
CONFIG_DIR="../k8s"

# Generate the certificate
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes -subj "/CN=${SERVICE}.${NAMESPACE}.svc"

# Create the Kubernetes secret
kubectl create secret tls ${SECRET} --cert=tls.crt --key=tls.key -n ${NAMESPACE} || kubectl delete secret ${SECRET} -n ${NAMESPACE} && kubectl create secret tls ${SECRET} --cert=tls.crt --key=tls.key -n ${NAMESPACE}

# Extract the CA bundle
CA_BUNDLE=$(cat tls.crt | base64 | tr -d '\n')

# Update webhook configuration files
sed -i.bak "s|<CA_BUNDLE>|${CA_BUNDLE}|g" ${CONFIG_DIR}/webhook-config.yaml
sed -i.bak "s|<CA_BUNDLE>|${CA_BUNDLE}|g" ${CONFIG_DIR}/mutation-webhook-config.yaml

echo "Certificates generated, secret created, and webhook configurations updated."
