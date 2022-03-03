vars = {
    APP = "golem"
}

recipe "go-dev" "local" {
    command {
        exec = "go run main.go"
    }
}

recipe "go-test" "local" {
    command {
        exec = "go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v"
    }
    command {
        exec = "go tool cover -html=coverage.out -o coverage.html"
    }
}

recipe "go-test-norace" "local" {
    command {
        exec = "go test ./... -timeout 15s -cover -coverprofile=coverage.out -v"
    }
    command {
        exec = "go tool cover -html=coverage.out -o coverage.html"
    }
}

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
        exec = "mkdir -p ./builds"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/{{.Vars.APP}}-linux-amd64-$(git describe --tags) main.go"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/{{.Vars.APP}}-darwin-amd64-$(git describe --tags) main.go"
    }
    command {
        exec = "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/{{.Vars.APP}}-windows-amd64-$(git describe --tags).exe main.go"
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