server_provider "terraform" {
    config = [
        "../testdata/terraform1.tfstate",
        "../testdata/terraform2.tfstate"
    ]
    user = "root"
    port = 22
}
