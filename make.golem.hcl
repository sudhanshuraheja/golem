vars = {
    APP = "golem"
    REPO = "sudhanshuraheja"
    ENV_PREFIX = "GOLEM_"
    GOTESTVERBOSE = "go test"
    GOTEST = "gotestsum"
}

recipe "dev" "local" {
    commands = [
        "go run main.go",
    ]
}

recipe "test" "local" {
    commands = [
        "go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v",
        "go tool cover -html=coverage.out -o coverage.html",
    ]
}

recipe "test_norace" "local" {
    commands = [
        "MallocNanoZone=0 go test ./... -timeout 15s -cover -coverprofile=coverage.out -v",
        "go tool cover -html=coverage.out -o coverage.html",
    ]
}

recipe "install" "local" {
    commands = [
        "./version.sh",
        "go install"
    ]
}

recipe "build" "local" {
    commands = [
        "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/golem-linux-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/golem-darwin-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/golem-windows-amd64-$(git describe --tags).exe main.go",
    ]
}

recipe "tidy" "local" {
    commands = [
        "go get -u",
        "go mod tidy"
    ]
}
