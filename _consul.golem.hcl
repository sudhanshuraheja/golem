vars = {
    HASHI_DC = "do1"
}

recipe "consul-local-setup" "local" {
    command {
        // consul-agent-ca.pem -> public key
        // consul-agent-ca-key.pem -> private key
        exec = "consul tls ca create"
    }
    command {
        // do1-server-consul-0.pem -> public key
        // do1-server-consul-0-key.pem -> private key
        exec = "consul tls cert create -server -dc {{.Vars.HASHI_DC}}"
    }
    command {
        // do1-client-consul-0.pem -> public key
        // do1-client-consul-0-key.pem -> private key
        exec = "consul tls cert create -client -dc {{.Vars.HASHI_DC}}"
    }
    commands = [
        "mv consul-agent-ca.pem {{.Vars.HASHI_PATH}}certs/",
        "mv consul-agent-ca-key.pem {{.Vars.HASHI_PATH}}certs/",
        "mv {{.Vars.HASHI_DC}}-server-consul-0.pem {{.Vars.HASHI_PATH}}certs/",
        "mv {{.Vars.HASHI_DC}}-server-consul-0-key.pem {{.Vars.HASHI_PATH}}certs/",
        "mv {{.Vars.HASHI_DC}}-client-consul-0.pem {{.Vars.HASHI_PATH}}certs/",
        "mv {{.Vars.HASHI_DC}}-client-consul-0-key.pem {{.Vars.HASHI_PATH}}certs/",
    ]
}
