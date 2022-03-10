loglevel = 5

vars = {
    SERVER_PROVIDERS = 1
}

server_provider "terraform" {
    config = [
        "../testdata/bad.tfstate",
    ]
    user = "root"
    port = 22
}
