OWNER=carlosjgp
OPERATOR=vault-secret-operator

# TODO automate and improve
CONTROLLER_VERSION=v0.0.1

DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}:${CONTROLLER_VERSION}

.PHONY: generate-crd generate-openapi build

build: generate-openapi
	operator-sdk build \
		--verbose \
		${DOCKER_IMAGE_CONTROLLER}

generate-crd:
	operator-sdk generate k8s

generate-openapi: generate-crd
	operator-sdk generate openapi

push: build
	docker push ${DOCKER_IMAGE_CONTROLLER}
