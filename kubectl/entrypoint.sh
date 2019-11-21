#!/bin/bash
while true; do
    kubectl create \
        secret generic \
        --from-file ${FOLDER} \
        --namespace ${NAMESPACE} \
        --dry-run \
        ${SECRET} | \
    kubectl apply -f -
    sleep 10
done