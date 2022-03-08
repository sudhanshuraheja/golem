recipe "go-install" "local" {
    commands = [
        "./version.sh",
        "go install",
    ]
}

recipe "go-build" "local" {
    commands = [
        "mkdir -p ./bin",
        "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/@golem.APP-linux-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/@golem.APP-darwin-amd64-$(git describe --tags) main.go",
        "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/@golem.APP-windows-amd64-$(git describe --tags).exe main.go",
    ]
}

recipe "go-tidy" "local" {
    commands = [
        // download the latest packages
        "go get -u",

        // run tidy
        "go mod tidy"
    ]
}