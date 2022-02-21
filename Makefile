# Variables
APP = golem
REPO = sudhanshuraheja
ENV_PREFIX = GOLEM_

TAG := $(shell git describe --tags)

# Colours
NO_COLOR = \x1b[0m
GRAY_COLOR = \x1b[30;01m
RED_COLOR = \x1b[31;01m
# GREEN_COLOR = \x1b[32;01m
YELLOW_COLOR = \x1b[33;01m
BLUE_COLOR = \x1b[34;01m
# FIVE_COLOR = \x1b[35;01m
GREEN_COLOR = \x1b[36;01m
WHITE_COLOR = \x1b[37;01m

.PHONY: test influx

dev:
	@echo "$(BLUE_COLOR)➤ Running dev$(NO_COLOR)"
	go run main.go --conf golem.hcl

test:
	@echo "$(BLUE_COLOR)➤ Running tests$(NO_COLOR)"
	MallocNanoZone=0 go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v \
		| sed ''/PASS/s//`printf "\033[34;01mPASS\033[0m"`/'' \
		| sed ''/FAIL/s//`printf "\033[31;01mFAIL\033[0m"`/'' \
		| sed ''/RUN/s//`printf "\033[30;01mRUN\033[0m"`/''
	go tool cover -html=coverage.out -o coverage.html

test_norace:
	@echo "$(BLUE_COLOR)➤ Running tests$(NO_COLOR)"
	MallocNanoZone=0 go test ./... -timeout 15s -cover -coverprofile=coverage.out -v \
		| sed ''/PASS/s//`printf "\033[34;01mPASS\033[0m"`/'' \
		| sed ''/FAIL/s//`printf "\033[31;01mFAIL\033[0m"`/'' \
		| sed ''/RUN/s//`printf "\033[30;01mRUN\033[0m"`/''
	go tool cover -html=coverage.out -o coverage.html

install:
	@echo "$(BLUE_COLOR)➤ Installing the binary$(NO_COLOR)"
	go install

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/golem-linux-amd64-$(TAG) main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/golem-darwin-amd64-$(TAG) main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/golem-windows-amd64-$(TAG).exe main.go

brew:
	code ../homebrew-golem