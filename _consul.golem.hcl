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
