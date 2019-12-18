helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm repo update
helm upgrade --install vault incubator/vault

# Because we are running Vault in dev mode
export VAULT_ADDR='http://127.0.0.1:8200'

vault auth enable kubernetes

vault write auth/kubernetes/config \
    token_reviewer_jwt=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) \
    kubernetes_host=https://kubernetes \
    kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt

vault policy write my-policy -<<EOF
path "auth/token/lookup-self" {
  capabilities = [ "read" ]
}
path "auth/token/create" {
  capabilities = ["create", "read", "update", "list"]
}

path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "secret/config" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "auth/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "auth/kubernetes/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
EOF

vault write auth/kubernetes/role/example \
    bound_service_account_names=example-vaultsecret \
    bound_service_account_namespaces=default \
    policies=my-policy,default \
    ttl=1m

vault kv put secret/config user=carlos password=gomez