vars = {
    APP = "golem"
    URL = "http://random"
}

server "test1" {
    public_ip = "1.2.3.4"
    private_ip = "10.11.12.13"
    hostname = ["localhost", "localhost2"]
    user = "sudhanshu"
    port = 22
    tags = ["golem"]
}

recipe "test2" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "@golem.APP"
    }
    kv {
        path = "test.testing_value"
        value = "rand32"
    }
    kv {
        path = "test.testing_value2"
        value = "rand32"
    }
    artifact {
        template {
            path = "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/nomad/nomad_client.template.hcl"
        }
        destination = "@golem.kv.test.testing_value"
    }
    artifact {
        template {
            path = "config"
        }
        destination = "config2"
    }
    artifact {
        template {
            data = <<EOF
testing_value:@golem.kv.test.testing_value
            EOF
        }
        destination = "config3"
    }
    artifact {
        source = "config"
        destination = "config4"
    }
    // commands = [
    //     "ls -la",
    //     "ls -ls m*",
    // ]
    script {
        apt {
            update = true
            pgp = "https://pgp-path.com"
            repository {
                url =  "https://repository.path.com"
                sources = "stable"
            }
            install = ["a", "b"]
            install_no_upgrade = ["c", "d"]
            purge = ["e", "f"]
        }
        command = "ls -la"
        commands = [
            "ls -la",
            "ls -la M*",
        ]
    }
}

recipe "test3" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "@golem.APP"
    }
}

recipe "test4" "remote" {
    kv {
        path = "test.testing_value"
        value = "rand32"
    }
}

recipe "test5" "remote" {
    artifact {
        source = "config"
        destination = "config4"
    }
}

recipe "test6" "remote" {
    commands = [
        "ls -la",
        "ls -ls m*",
    ]
}

recipe "test7" "remote" {
    script {
        apt {
            update = true
        }
    }
}

recipe "test8" "local" {
    kv {
        path = "test.local_testing_value"
        value = "rand32"
    }
    artifact {
        template {
            path = "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/LICENSE"
        }
        destination = "./LICENCE"
    }
    artifact {
        template {
            path = "config"
        }
        destination = "config2"
    }
    artifact {
        template {
            data = <<EOF
some random data
            EOF
        }
        destination = "config2"
    }
    artifact {
        source = "config"
        destination = "config2"
    }
    commands = [
        "ls -la",
        "ls -ls m*",
    ]
    script {
        apt {
            update = true
            pgp = "https://pgp-path.com"
            repository {
                url =  "https://repository.path.com"
                sources = "stable"
            }
            install = ["a", "b"]
            install_no_upgrade = ["c", "d"]
            purge = ["e", "f"]
        }
        command = "ls -la"
        commands = [
            "ls -la",
            "ls -la M*",
        ]
    }
}
