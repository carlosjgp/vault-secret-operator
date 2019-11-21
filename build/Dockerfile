FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/vault-secret-operator \
    USER_UID=1001 \
    USER_NAME=vault-secret-operator

COPY build/bin /usr/local/bin
COPY templates /templates
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}

# install operator binary
COPY build/_output/bin/vault-secret-operator ${OPERATOR}
