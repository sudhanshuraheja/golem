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