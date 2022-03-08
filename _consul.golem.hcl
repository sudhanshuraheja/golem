vars = {
    HASHI_DC = "do1"
    HASHI_PATH = "./testdata/nomad/"
}

recipe "consul-local-setup" "local" {
    commands = [
        // consul-agent-ca.pem -> public key
        // consul-agent-ca-key.pem -> private key
        "consul tls ca create",
        "mv consul-agent-ca.pem @golem.HASHI_PATH/certs/",
        "mv consul-agent-ca-key.pem @golem.HASHI_PATH/certs/",

        // do1-server-consul-0.pem -> public key
        // do1-server-consul-0-key.pem -> private key
        "consul tls cert create -server -dc @golem.HASHI_DC",
        "mv @golem.HASHI_DC-server-consul-0.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-server-consul-0-key.pem @golem.HASHI_PATH/certs/",

        // do1-client-consul-0.pem -> public key
        // do1-client-consul-0-key.pem -> private key
        "consul tls cert create -client -dc @golem.HASHI_DC",
        "mv @golem.HASHI_DC-client-consul-0.pem @golem.HASHI_PATH/certs/",
        "mv @golem.HASHI_DC-client-consul-0-key.pem @golem.HASHI_PATH/certs/",

        "openssl rand 32 | base64 > @golem.HASHI_PATH/certs/consul.key",
    ]
}
