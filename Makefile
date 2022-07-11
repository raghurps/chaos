BASEDIR := $(shell pwd)

DOCKER := docker
DOCKER_BUILD_FLAGS := build --pull --rm -f
DOCKER_PUSH_FLAGS := push

DOCKERFILE := "Dockerfile.chaosmonkey"


IMAGE_NAME := rpsingh/chaosmonkey
IMAGE_TAG := $(shell git rev-parse --verify HEAD)

HELM := helm
HELM_INSTALL_ARGS := upgrade --install --atomic --timeout 120s
HELM_CHART := ${BASEDIR}/charts/chaos
HELM_VALUES_FILE := ${BASEDIR}/hack/values.yaml
HELM_RELEASE_NAME := chaosmonkey
NAMESPACE := chaosmonkey
KUBECTL_CREATE_DRY_RUN := kubectl create namespace --dry-run=client -oyaml
KUBECTL_APPLY := kubectl apply -f

build:
	echo "Building image ${IMAGE_NAME}:${IMAGE_TAG} for chaosmonkey"
	${DOCKER} ${DOCKER_BUILD_FLAGS} ${DOCKERFILE} -t ${IMAGE_NAME}:${IMAGE_TAG} .

release: build
	echo "Publishing image ${IMAGE_NAME}:${IMAGE_TAG}"
	${DOCKER} ${DOCKER_PUSH_FLAGS} ${IMAGE_NAME}:${IMAGE_TAG}

install:
	echo "Creating namespace ${NAMESPACE}"
	${KUBECTL_CREATE_DRY_RUN} ${NAMESPACE} | ${KUBECTL_APPLY} -

	echo "Installing chaosmonkey..."
	${HELM} ${HELM_INSTALL_ARGS} ${HELM_RELEASE_NAME} \
		--namespace=${NAMESPACE} \
		${HELM_CHART} -f ${HELM_VALUES_FILE} \
		--set image.repository=${IMAGE_NAME} \
		--set image.tag=${IMAGE_TAG}

uninstall:
	echo "Uninstalling chaosmonkey..."
	${HELM} delete --namespace=${NAMESPACE} ${HELM_RELEASE_NAME}
