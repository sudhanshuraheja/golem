vars = {
    APP = "golem"
    REPO = "sudhanshuraheja"
    ENV_PREFIX = "GOLEM_"
    GOTESTVERBOSE = "go test"
    GOTEST = "gotestsum"
}

recipe "dev" {
    type = "local-exec"
    commands = [
        "go run main.go",
    ]
}

recipe "test" {
    type = "local-exec"
    commands = [
        "go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v",
        "go tool cover -html=coverage.out -o coverage.html",
    ]
}

recipe "test_norace" {
    type = "local-exec"
    commands = [
        "MallocNanoZone=0 go test ./... -timeout 15s -cover -coverprofile=coverage.out -v",
        "go tool cover -html=coverage.out -o coverage.html",
    ]
}

recipe "install" {
    type = "local-exec"
    commands = [
        "go install"
    ]
}

recipe "build" {
    type = "local-exec"
    commands = [
        "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/golem-linux-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/golem-darwin-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/golem-windows-amd64-$(git describe --tags).exe main.go",
    ]
}

recipe "tidy" {
    type = "local-exec"
    commands = [
        "go get -u",
        "go mod tidy"
    ]
}
