OWNER=carlosjgp
OPERATOR=vault-secret-operator

DOCKER_IMAGE_KUBECTL=${OWNER}/${OPERATOR}-kubectl
DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}-controller

#TODO run operator-sdk inside the docker image


build: build-operator build-kubectl

build-operator:
	operator-sdk build \
		${DOCKER_IMAGE_CONTROLLER}

build-kubectl:
	docker build \
		-t ${DOCKER_IMAGE_KUBECTL} \
		-f kubectl/Dockerfile \
		kubectl

generate-crd:
	operator-sdk generate k8s

generate-openapi:
	operator-sdk generate openapi

docker-push:
	docker push ${DOCKER_IMAGE_CONTROLLER}
	docker push ${DOCKER_IMAGE_KUBECTL}