#!/usr/bin/env bash
set -o errexit

IMAGE_NAME="${IMAGE_NAME:-rpsingh/chaosmonkey}"
IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --verify HEAD)}"

RELEASE_NAME="${RELEASE_NAME:-chaosmonkey}"
NAMESPACE="${NAMESPACE:-chaosmonkey}"

echo "Uninstalling chaosmonkey..."
helm delete --namespace=${NAMESPACE} ${RELEASE_NAME}