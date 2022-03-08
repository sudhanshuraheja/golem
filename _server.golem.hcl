recipe "server-all" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    script {
        command = "sudo dpkg --configure -a"
    }
}

recipe "server-apt-update" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    script {
        apt {
            update = true
        }
    }
}

recipe "server-reboot" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "nomad"
    }
    commands = [
        "sudo reboot"
    ]
}

recipe "server-apply-security-patch" "remote" {
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