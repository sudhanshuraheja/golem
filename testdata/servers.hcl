loglevel = 4

vars = {
    SERVERS = 4
}

server "test1" {
    hostname = ["localhost"]
    user = "sudhanshu"
    port = 22
    tags = ["test-tag"]
}

server "test2" {
    user = "sudhanshu"
    port = 22
}

server "test3" {
    hostname = []
    user = "sudhanshu"
    port = 22
    tags = []
}

server "test4" {
    public_ip = "1.2.3.4"
    private_ip = "10.11.12.13"
    hostname = ["localhost", "localhost2"]
    user = "sudhanshu"
    port = 22
    tags = ["test-tag", "tag_1", "tag*3", "tag4"]
}