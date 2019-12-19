OWNER=carlosjgp
OPERATOR=vault-secret-operator

# TODO automate and improve
CONTROLLER_VERSION=v0.0.1
KUBECTL_VERSION=v1.16.0

DOCKER_IMAGE_CONTROLLER=${OWNER}/${OPERATOR}-controller:${CONTROLLER_VERSION}
DOCKER_IMAGE_KUBECTL=${OWNER}/${OPERATOR}-kubectl:${KUBECTL_VERSION}-${CONTROLLER_VERSION}


.PHONY: generate-crd build-kubectl build-controller build push-kubectl push-controller push kubectl controller

build: build-kubectl build-controller
push: push-kubectl push-controller

generate-crd:
	operator-sdk generate k8s
	operator-sdk generate crds


build-controller: generate-crd
	operator-sdk build \
		${DOCKER_IMAGE_CONTROLLER}

push-controller: build-controller
	docker push ${DOCKER_IMAGE_CONTROLLER}

controller: build-controller push-controller

build-kubectl:
	docker build \
		--build-arg KUBECTL_VERSION=${KUBECTL_VERSION} \
		-t ${DOCKER_IMAGE_KUBECTL} \
		kubectl

push-kubectl: build-kubectl
	docker push ${DOCKER_IMAGE_KUBECTL}

kubectl: build-kubectl push-kubectl
