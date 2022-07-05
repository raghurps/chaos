#!/usr/bin/env bash
set -o errexit

IMAGE_NAME=${IMAGE_NAME:-"rpsingh/chaosmonkey"}
IMAGE_TAG=${IMAGE_TAG:-$(git rev-parse --verify HEAD)}

echo "Building image ${IMAGE_NAME}:${IMAGE_VERSION} for chaosmonkey"
docker build --pull --rm -f "Dockerfile.chaosmonkey" -t ${IMAGE_NAME}:${IMAGE_TAG} .