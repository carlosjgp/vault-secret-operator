# Vault sercret Operator

Reads secrets from Vault and creates configmaps or secrets on Kubenetes.
It uses Vault agent to keep the secrets always InSync with between Vault
and Kubernetes

## TODO

### Operator

- Gihub releases
- Start using good commit messages for the Changelog
- Version Docker images
- How to manage Service accounts and Roles
- How to manage VaultSecret deletion and K8s secret. For now we do nothing

### Testing

- Helm values for the tests
- Commands to configure Vault
- Enable K8s auth in Vault
- Link RBAC - ServiceAccount
- Add secret to test to Vault KV


## Containers

### Vault agent

https://learn.hashicorp.com/vault/identity-access-management/vault-agent-k8s

Ideally we will be using Kubernetes auth backend

### Consul Template

Also referenced on he Vault Agent documentation
https://learn.hashicorp.com/vault/identity-access-management/vault-agent-k8s

It can fetch secrets and template them into a more complex strings, multiline,...

### Kubernetes client

To be able to use the generated files by Consul Template and crete ConfigMap or Secret

## ConfigMaps

### Vault Agent

### Consul template

## Build

Created using operator-sdk for GoLang

https://github.com/operator-framework/operator-sdk

### Commands


```bash
operator-sdk generate k8s &&\
operator-sdk generate openapi &&\
operator-sdk build carlosjgp/vault-secret-operator
```

To publish the image to DockerHub (credentials required)
```
docker push carlosjgp/vault-secret-operator
```

For local development with Minikube use
```
minikube cache add carlosjgp/vault-secret-operator
```


## Test

### Minikube

Using Minikube with Kubernetes 1.15 since the deprecated APIs have been removed on 1.16 and there are still loads of people who hasn't migrated to the latest API spec

#### Vault

Deployed `incubator/vault` Helm chart for Vault version `1.2.3`

#### Operator

Then used the `deploy` folder on this repo
