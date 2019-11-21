OWNER=carlosjgp
OPERATOR=vault-secret-operator

# TODO automate and improve
CONTROLLER_VERSION=v0.0.1
KUBECTL_VERSION=v1.16.0

DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}-controller:${CONTROLLER_VERSION}
DOCKER_IMAGE_KUBECTL=${OWNER}/${OPERATOR}-kubectl:${KUBECTL_VERSION}-${CONTROLLER_VERSION}


.PHONY: generate-crd generate-openapi build

build: generate-openapi
	operator-sdk build \
		${DOCKER_IMAGE_CONTROLLER}
	docker build \
		--build-arg KUBECTL_VERSION=${KUBECTL_VERSION} \
		-t ${DOCKER_IMAGE_KUBECTL} \
		kubectl

generate-crd:
	operator-sdk generate k8s

generate-openapi: generate-crd
	operator-sdk generate openapi

push: build
	docker push ${DOCKER_IMAGE_CONTROLLER}
	docker push ${DOCKER_IMAGE_KUBECTL}
