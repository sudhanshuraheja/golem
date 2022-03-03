vars = {
    HASHI_PATH = "./hashicorp/"
}

recipe "nomad-ls" "local" {
    command {
        exec = "nomad node status -allocs"
    }
    command {
        exec = "nomad job status"
    }
    command {
        exec = "nomad deployment list"
    }
}

recipe "nomad-local-setup" "local" {
    command {
        exec = "mkdir -p {{.Vars.HASHI_PATH}}certs/"
    }
    artifact {
        template {
            data = <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h",
      "usages": ["signing", "key encipherment", "server auth", "client auth"]
    }
  }
}
            EOF
        }
        destination = "{{.Vars.HASHI_PATH}}certs/cfssl.json"
    }
    command {
        // Install cfssl
        exec = "go install github.com/cloudflare/cfssl/cmd/cfssl@latest"
    }
    command {
        // Install cfssljson
        exec = "go install github.com/cloudflare/cfssl/cmd/cfssljson@latest"
    }
    command {
        // HASHI_PATH/certs/nomad-ca-key.pem -> private key
        // HASHI_PATH/certs/nomad-ca.csr -> certificate signing request
        // HASHI_PATH/certs/nomad-ca.pem -> public key
        exec = "cfssl print-defaults csr | cfssl gencert -initca - | cfssljson -bare {{.Vars.HASHI_PATH}}certs/nomad-ca"
    }
    command {
        // server-key.pem -> private key
        // server.csr -> certificate signing request
        // server.pem -> public key
        exec = <<EOF
echo '{}' | cfssl gencert -ca={{.Vars.HASHI_PATH}}certs/nomad-ca.pem -ca-key={{.Vars.HASHI_PATH}}certs/nomad-ca-key.pem -config={{.Vars.HASHI_PATH}}certs/cfssl.json -hostname="server.global.nomad,localhost,127.0.0.1,
{{- range $_, $s := (match "tags" "contains" "nomad-server") -}}
    {{- if not ($s).PublicIP -}}
    {{- else -}}
        {{- ($s).PublicIP -}},
    {{- end -}}
    {{- if not ($s).PrivateIP -}}
    {{- else -}}
        {{- ($s).PrivateIP -}},
    {{- end -}}
{{- end -}}" - | cfssljson -bare {{.Vars.HASHI_PATH}}certs/server
        EOF
    }
    command {
        // client-key.pem -> private key
        // client.csr -> certificate signing request
        // client.pem -> public key
        exec = <<EOF
echo '{}' | cfssl gencert -ca={{.Vars.HASHI_PATH}}certs/nomad-ca.pem -ca-key={{.Vars.HASHI_PATH}}certs/nomad-ca-key.pem -config={{.Vars.HASHI_PATH}}certs/cfssl.json -hostname="client.global.nomad,localhost,127.0.0.1,
{{- range $_, $s := (match "tags" "contains" "nomad-server") -}}
    {{- if not ($s).PublicIP -}}
    {{- else -}}
        {{- ($s).PublicIP -}},
    {{- end -}}
    {{- if not ($s).PrivateIP -}}
    {{- else -}}
        {{- ($s).PrivateIP -}},
    {{- end -}}
{{ end }}" - | cfssljson -bare {{.Vars.HASHI_PATH}}certs/client
        EOF
    }
    command {
        exec = "echo '{}' | cfssl gencert -ca={{.Vars.HASHI_PATH}}certs/nomad-ca.pem -ca-key={{.Vars.HASHI_PATH}}certs/nomad-ca-key.pem -profile=client - | cfssljson -bare {{.Vars.HASHI_PATH}}certs/cli"
    }
    command {
        exec = "openssl rand 32 | base64 > {{.Vars.HASHI_PATH}}certs/nomad.key"
    }
}

recipe "nomad-remote-setup" "remote" {
    match {
        attribute = "name"
        operator = "="
        value = "skye-c2"
    }
    command {
        apt {
            update = true
        }
        apt {
            pgp = "https://download.docker.com/linux/ubuntu/gpg"
            repository {
                url = "https://download.docker.com/linux/ubuntu"
                sources = "stable"
            }
            install = ["docker-ce", "docker-ce-cli", "containerd.io"]
        }
        apt {
            pgp = "https://apt.releases.hashicorp.com/gpg"
            repository {
                url = "https://apt.releases.hashicorp.com"
                sources = "main"
            }
            install_no_upgrade = ["consul", "nomad"]
        }
    }
}

recipe "nomad-server-config-update" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad-server"
    }
    artifact {
        source = "configs/nomad_server.hcl"
        destination = "/etc/nomad.d/nomad.hcl"
    }
    artifact {
        source = "certs/nomad-ca.pem"
        destination = "/etc/nomad.d/nomad-ca.pem"
    }
    artifact {
        source = "certs/server.pem"
        destination = "/etc/nomad.d/server.pem"
    }
    artifact {
        source = "certs/server-key.pem"
        destination = "/etc/nomad.d/server-key.pem"
    }
    commands = [
        "chown nomad:nomad /etc/nomad.d/server-key.pem",
        "systemctl daemon-reload",
        "systemctl stop nomad",
        "systemctl start nomad",
    ]
}

recipe "nomad-client-config-update" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad-client"
    }
    artifact {
        source = "configs/nomad_client.hcl"
        destination = "/etc/nomad.d/nomad.hcl"
    }
    artifact {
        source = "certs/nomad-ca.pem"
        destination = "/etc/nomad.d/nomad-ca.pem"
    }
    artifact {
        source = "certs/client.pem"
        destination = "/etc/nomad.d/client.pem"
    }
    artifact {
        source = "certs/client-key.pem"
        destination = "/etc/nomad.d/client-key.pem"
    }
    commands = [
        "chown nomad:nomad /etc/nomad.d/client-key.pem",
        "mkdir -p /opt/caddy",
        "chown nomad:nomad /opt/caddy",
        "systemctl daemon-reload",
        "systemctl stop nomad",
        "systemctl start nomad",
    ]
}

recipe "nomad-clean" "local" {
    commands = [
        "rm *.json",
        "rm *.pem",
        "rm *.csr",
    ]
}

recipe "nomad-tail" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    command {
        exec = "journalctl -f -u nomad.service"
    }
}

