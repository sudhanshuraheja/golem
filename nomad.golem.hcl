recipe "tail-nomad" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    commands = [
        "journalctl -f -u nomad.service"
    ]
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
