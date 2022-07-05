#!/usr/bin/env bash
set -o errexit

IMAGE_NAME="${IMAGE_NAME:-rpsingh/chaosmonkey}"
IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --verify HEAD)}"
RELEASE_NAME="${RELEASE_NAME:-chaosmonkey}"

NAMESPACE="${NAMESPACE:-chaosmonkey}"

# create namespace if it does not exist
echo "Creating namespace ${NAMESPACE}"
kubectl create ns ${NAMESPACE} --dry-run=client -oyaml | kubectl apply -f -

echo "Installing chaosmonkey..."
helm upgrade --install --atomic --timeout 120s \
    --namespace=${NAMESPACE} \
    ${RELEASE_NAME} ${PWD}/charts/chaos -f ${PWD}/hack/values.yaml \
		--set image.repository=${IMAGE_NAME} \
		--set image.tag=${IMAGE_TAG}