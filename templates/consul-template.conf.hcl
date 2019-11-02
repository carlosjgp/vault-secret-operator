vault {
  # Same as the Vault agent
  address = "http://127.0.0.1:8200"
  vault_agent_token_file = "/tmp/vault/agent/token"

}

{{ .ConsulTemplates }}