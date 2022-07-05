# chaos

## ChaosMonkey

ChaosMonkey is losely based on Netflix's [project](https://netflix.github.io/chaosmonkey/)  with same name. Its purpose is to bring chaos in kubernetes eco-system by randomly terminating one of the pods. This ensures that application developers are writing code that is resilient enough to recover from such random mishaps.

### Installation

#### Requirements

In order to deploy chaosmonkey, you'll need [helm](https://helm.sh/) CLI installed. To install helm, run:

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

You'll also need access to a kubernetes cluster. If you don't have one, you can create it locally using [kind](https://kind.sigs.k8s.io/) and [docker](https://www.docker.com/).

```bash
# install kind binary
go install sigs.k8s.io/kind@v0.14.0 
# create kind cluster of 2 nodes using the config file
kind create cluster --config hack/kindConfig.yaml
```
**or**

Use setup script that creates a local registry accessible at localhost:5001 and a kind cluster that can pull images from this local registry

```bash
./hack/setupKindCluster.sh
```

To remove this local registry and kind cluster, run:

```bash
./hack/cleanupKindCluster.sh
```

#### Steps to install

If you have `make` utility installed then run:

```bash
make install
```

**or** 

In case you don't have `make` utility, use shell script [install.sh](./install.sh)

```bash
./install.sh
```

To provide custom image and tag name, you can provider IMAGE_NAME and IMAGE_TAG variables. e.g. if you have pushed image to you local registry:

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" install
```
**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./install.sh
```

Default image tag is the latest git commit id and the default image name is `rpsingh/chaosmonkey`.

The installation happens via helm CLI and this [`values.yaml`](./hack/values.yaml) file.

#### Steps to uninstall

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" uninstall
```

**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./uninstall.sh
```

### How to Contribute

In order to contribute to this repo, you need to setup local dev environment.

#### Dev Environment

To develop and test locally, you must install
- [Docker](https://www.docker.com/) v20.10.16
- [Helm](https://helm.sh/) v3.8.2
- [Golang](https://go.dev/) v1.17
- [Kind](https://kind.sigs.k8s.io/) v0.14.0
- [GNU Make](https://www.gnu.org/software/make/) (Optional)

#### Build Image

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" build
```

**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./build.sh
```
  
#### Publish Image

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" release
```

**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./release.sh
```

#### Install on kubernetes cluster

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" install
```

**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./install.sh
```

#### Uninstall from kubernetes cluster

```bash
# You can provide your own image name and tag
make IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" uninstall
```

**or**

```bash
# You can provide your own image name and tag
IMAGE_NAME="localhost:5001/chaosmonkey" IMAGE_TAG="v0.0.1" ./uninstall.sh
```

#### Build locally outside container

```bash
go build -o chaosmonkey cmd/chaosmonkey/main.go
```

For flag details

```bash
go run cmd/chaosmonkey/main.go  --help
```

#### Run unit tests

```bash
go test -v -cover ./...
```
**or**

```bash
go test -v -cover chaosmonkey.monke/chaos/...
```


### TODO
- [ ] Make Kubernetes client package testable
- [ ] Add more unit tests
- [ ] Refactor function in main package
