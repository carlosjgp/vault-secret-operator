OWNER=carlosjgp
OPERATOR=vault-secret-operator

# TODO automate and improve
CONTROLLER_VERSION=v0.0.1
KUBECTL_VERSION=v1.16-0.0.1

DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}-controller:${CONTROLLER_VERSION}
DOCKER_IMAGE_KUBECTL=${OWNER}/${OPERATOR}-kubectl:${KUBECTL_VERSION}

#TODO run operator-sdk inside the docker image
build: build-operator build-kubectl

build-operator:
	cd operator && \
	operator-sdk build \
		${DOCKER_IMAGE_CONTROLLER}

build-kubectl:
	docker build \
		-t ${DOCKER_IMAGE_KUBECTL} \
		-f kubectl/Dockerfile \
		kubectl

generate-crd:
	cd operator && \
	operator-sdk generate k8s

generate-openapi:
	cd operator && \
	operator-sdk generate openapi

docker-push: docker-push-controller docker-push-kubectl

docker-push-kubectl:
	docker push ${DOCKER_IMAGE_KUBECTL}

docker-push-controller:
	docker push ${DOCKER_IMAGE_CONTROLLER}
