install:
	go run main.go go-install

test:
	go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html