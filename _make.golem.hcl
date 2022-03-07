recipe "go-install" "local" {
    command {
        exec = "./version.sh"
    }
    command {
        exec = "go install"
    }
}

recipe "go-build" "local" {
    command {
        exec = "mkdir -p ./bin"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/@golem.APP-linux-amd64-$(git describe --tags) main.go"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/@golem.APP-darwin-amd64-$(git describe --tags) main.go"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/@golem.APP-windows-amd64-$(git describe --tags).exe main.go"
    }
}

recipe "go-tidy" "local" {
    command {
        // download the latest packages
        exec = "go get -u"
    }
    command {
        // run tidy
        exec = "go mod tidy"
    }
}