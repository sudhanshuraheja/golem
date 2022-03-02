recipe "nomad-setup" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
//     command {
//         // add GPG for docker
//         exec = "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -"
//     }
//     command {
//         // add repo for docker
//         exec = <<EOF
// sudo apt-add-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
//         EOF
//     }
//     command {
//         // add hashicorp GPG for consul and nomad
//         exec = "curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -"
//     }
//     command {
//         // add hashicorp repo for consul and nomad
//         exec = <<EOF
// sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
//         EOF
//     }
    command {
        // update apt
        apt {
            update = true
        }
        // exec = "sudo apt-get update --quiet --assume-yes --show-upgraded"
    }
    command {
        exec = "sudo apt-get install docker-ce docker-ce-cli containerd.io --quiet --assume-yes --show-upgraded"
    }
    // command {
    //     exec = "sudo apt-get install consul nomad --no-upgrade --quiet --assume-yes --show-upgraded"
    // }
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
