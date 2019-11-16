# TODO create an operator-sdk image on a different repo
FROM golang:1.13-alpine

ENV RELEASE_VERSION=v0.12.0

RUN wget https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu && \
    chmod +x operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu && \
    mkdir -p /usr/local/bin/ && \
    cp operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk && \
    rm operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

COPY operator /operator

RUN operator-sdk generate k8s && \
    operator-sdk generate openapi


FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/vault-secret-operator \
    USER_UID=1001 \
    USER_NAME=vault-secret-operator

# install operator binary
COPY build/_output/bin/vault-secret-operator ${OPERATOR}

COPY build/bin /usr/local/bin
COPY templates /templates
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
