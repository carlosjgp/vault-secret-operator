ARG KUBECTL_VERSION
FROM carlosjgp/kubectl:${KUBECTL_VERSION}

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]