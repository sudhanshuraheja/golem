vars = {
    APP = "golem"
    REPO = "sudhanshuraheja"
    ENV_PREFIX = "GOLEM_"
    GOTESTVERBOSE = "go test"
    GOTEST = "gotestsum"
}

recipe "apt-update" "remote" {
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

recipe "ls-la" "local" {
    commands = [
        "ls -la {{.Vars.APP}}*",
    ]
}

recipe "ls-la-remote" "remote" {
    match {
        attribute = "name"
        operator = "="
        value = "skye-c3"
    }
    commands = [
        "ls -la {{.Vars.APP}}*",
        "ping {{ (matchOne \"name\" \"=\" \"skye-c2\").PrivateIP  }}",
    ]
}

recipe "test-exec" "remote" {
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

recipe "apply-security-patch" "remote" {
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