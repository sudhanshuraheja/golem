server_provider "terraform" {
    config = []
    user = "sudhanshu"
    port = "22"
}

server "test-server" {
    public_ip = "127.0.0.1"
    private_ip = "127.0.0.1"
    hostname = ["local"]
    user = "sudhanshu"
    port = 22
    tags = ["test-tag"]
}

loglevel = 3

max_parallel_processes = 5

vars = {
    APP = "golem"
}

recipe "recipe_test" "remote" {
    match {
        attribute = "tags"
        operator = "contains"
        value = "test-tag"
    }
    kv {
        path = "test.testing_value"
        value = "rand32"
    }
    artifact {
        source = "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/LICENSE"
        destination = "./LICENCE"
    }
    artifact {
        source = "config"
        destination = "config2"
    }
    command {
        apt {
            update = true
        }
        apt {
            install = ["a", "b"]
        }
    }
}