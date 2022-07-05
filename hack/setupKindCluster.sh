#!/usr/bin/env bash
set -o errexit

BASEDIR="$( dirname -- "$0"; )"
export PATH=$(go env GOPATH)/bin:$PATH

if [ -x kind ]; then
    go install sigs.k8s.io/kind@v0.14.0
    kind version
fi

echo "Creating local registry"
docker run \
  -d --restart=always -p "127.0.0.1:5001:5000" --name "kind-registry" \
  registry:2

echo "Creating kind cluster"
kind create cluster --config ${BASEDIR}/kindConfig.yaml

echo "Connecting registry to kind cluster network"
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' kind-registry)" = 'null' ]; then
  docker network connect "kind" kind-registry
fi

echo "Creating configMap to document local registry"
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:5001"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo "Installing additional workload for chaosmonkey to devour"
kubectl apply -f ${BASEDIR}/workloads.yaml