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
        "nomad version",
        "consul version",
    ]
}

recipe "tail-nomad" {
    type = "remote-exec"
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    commands = [
        "journalctl -f -u nomad.service"
    ]
}
