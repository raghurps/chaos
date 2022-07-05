#!/usr/bin/env bash
set -o errexit

IMAGE_NAME="${IMAGE_NAME:-rpsingh/chaosmonkey}"
IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --verify HEAD)}"

echo "Publishing image ${IMAGE_NAME}:${IMAGE_TAG}"
docker push ${IMAGE_NAME}:${IMAGE_TAG}