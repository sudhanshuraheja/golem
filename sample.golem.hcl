recipe "apt-update" {
    type = "remote-exec"
    match {
        attribute = "tags"
        operator = "not-contains"
        value = "local"
    }
    commands = [
        "apt-get update",
        "nomad version"
    ]
}

recipe "ls-la" {
    type = "local-exec"
    commands = [
        "ls -la",
    ]
}

recipe "test-exec" {
    type = "remote-exec"
    match {
        attribute = "name"
        operator = "like"
        value = "skye-c3"
    }
    artifact {
        source = "https://github.com/sudhanshuraheja/golem/releases/download/v0.1.0/golem-linux-amd64-v0.1.0"
        destination = "golem"
    }
    artifact {
        source = "LICENSE"
        destination = "LICENSE"
    }
    commands = [
        "ls -la L*",
        "ls -la g*"
    ]
}

recipe "apply-security-patch" {
    type = "remote-exec"
    match {
        attribute = "name"
        operator = "="
        value = "skye-s3"
    }
    commands = [
        "apt-get update",
        "apt-get install unattended-upgrades",
        "unattended-upgrade",
    ]
}