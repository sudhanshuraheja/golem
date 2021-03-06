install:
	./version.sh
	go install

test:
	./version.sh
	go test ./... -timeout 15s -race -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out -o coverage.func