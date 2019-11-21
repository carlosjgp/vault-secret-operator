pid_file = "./pidfile"

vault {
    address = {{ .VaultAddress | quote }}
}

auto_auth {
{{ .AutoAuthMethod | indent 4 }}

    sink "file" {
        config = {
            path = "/tmp/vault/agent/token"
        }
    }
}

cache {
    use_auto_auth_token = true
}

listener "tcp" {
    address = "127.0.0.1:8200"
    tls_disable = true
}