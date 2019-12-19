OWNER=carlosjgp
OPERATOR=vault-secret-operator

# TODO automate and improve
CONTROLLER_VERSION=v0.0.1

DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}:${CONTROLLER_VERSION}


.PHONY: generate-crd build push

generate-crd:
	operator-sdk generate k8s
	operator-sdk generate crds

build: generate-crd
	operator-sdk build \
		${DOCKER_IMAGE_CONTROLLER}

push: build
	docker push ${DOCKER_IMAGE_CONTROLLER}
