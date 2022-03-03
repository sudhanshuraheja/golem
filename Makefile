TAG := $(shell git describe --tags)

install:
	go run main.go go-install