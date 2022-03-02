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