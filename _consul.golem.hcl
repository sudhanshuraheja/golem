vars = {
    HASHI_DC = "do1"
    HASHI_PATH = "./nomad/"
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
        exec = "consul tls cert create -server -dc @golem.HASHI_DC"
    }
    command {
        // do1-client-consul-0.pem -> public key
        // do1-client-consul-0-key.pem -> private key
        exec = "consul tls cert create -client -dc @golem.HASHI_DC"
    }
    command {
        exec = "openssl rand 32 | base64 > @golem.HASHI_PATH/certs/consul.key"
    }
    commands = [
        "mv consul-agent-ca.pem @golem.HASHI_PATH/certs/",
        "mv consul-agent-ca-key.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-server-consul-0.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-server-consul-0-key.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-client-consul-0.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-client-consul-0-key.pem @golem.HASHI_PATH/certs/",
    ]
}
