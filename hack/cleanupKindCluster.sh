#!/usr/bin/env bash

set -o errexit

BASEDIR="$( dirname -- "$0"; )"

echo "Deleting local registry"
docker rm -f kind-registry

echo "Deleting local kind cluster"
kind delete clusters k8s-diy