recipe "install" "local" {
    // Install cfssl
    commands = [
        "go install github.com/cloudflare/cfssl/cmd/cfssl@latest",
        "go install github.com/cloudflare/cfssl/cmd/cfssljson@latest",
    ]
}

recipe "consul-ca" "local" {
    command {
        // consul-agent-ca.pem -> public key
        // consul-agent-ca-key.pem -> private key
        exec = "consul tls ca create"
    }
    command {
        // dc1-server-consul-0.pem -> public key
        // dc1-server-consul-0-key.pem -> private key
        exec = "consul tls cert create -server -dc do1"
    }
    command {
        exec = "consul tls cert create -client -dc do1"
    }
}

recipe "nomad-ca" "local" {
    command {
        exec = <<EOF
echo '{ "signing": { "default": { "expiry": "87600h", "usages": ["signing", "key encipherment", "server auth", "client auth"] } } }' > cfssl.json
EOF
    }
    command {
        // nomad-ca-key.pem -> private key
        // nomad-ca.csr -> certificate signing request
        // nomad-ca.pem -> public key
        exec = "cfssl print-defaults csr | cfssl gencert -initca - | cfssljson -bare nomad-ca"
    }
    command {
        // server-key.pem -> private key
        // server.csr -> certificate signing request
        // server.pem -> public key
        exec = <<EOF
echo '{}' | cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -config=cfssl.json -hostname="server.global.nomad,localhost,127.0.0.1,{{ range $_, $s := (match "tags" "contains" "nomad-server") }}{{ if not ($s).PublicIP }}{{ else }}{{ ($s).PublicIP }},{{ end }}{{ if not ($s).PrivateIP }}{{ else }}{{ ($s).PrivateIP }},{{ end }}{{ end }}" - | cfssljson -bare server
        EOF
    }
    command {
        // client-key.pem -> private key
        // client.csr -> certificate signing request
        // client.pem -> public key
        exec = <<EOF
echo '{}' | cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -config=cfssl.json -hostname="client.global.nomad,localhost,127.0.0.1,{{ range $_, $s := (match "tags" "contains" "nomad-server") }}{{ if not ($s).PublicIP }}{{ else }}{{ ($s).PublicIP }},{{ end }}{{ if not ($s).PrivateIP }}{{ else }}{{ ($s).PrivateIP }},{{ end }}{{ end }}" - | cfssljson -bare client
        EOF
    }
    command {
        exec = "echo '{}' | cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -profile=client - | cfssljson -bare cli"
    }
}

recipe "clean" "local" {
    commands = [
        "rm *.json",
        "rm *.pem",
        "rm *.csr",
    ]
}