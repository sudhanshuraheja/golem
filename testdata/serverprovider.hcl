loglevel = 5
max_parallel_processes = 4

vars = {
    SERVER_PROVIDERS = 1
}

server_provider "terraform" {
    config = [
        "../testdata/terraform1.tfstate",
        "../testdata/terraform2.tfstate"
    ]
    user = "root"
    port = 22
}
