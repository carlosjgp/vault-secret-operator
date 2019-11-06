FROM alpine as downloader

WORKDIR /tmp

RUN apk add --no-cache curl

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl

COPY entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

FROM alpine

WORKDIR /

RUN apk add --no-cache bash

COPY --from=downloader /tmp/kubectl /usr/local/bin/kubectl
COPY --from=downloader /tmp/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]