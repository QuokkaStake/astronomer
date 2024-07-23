VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
LDFLAGS = -X main.version=${VERSION}

build:
	go build -ldflags '$(LDFLAGS)' cmd/astronomer.go

install:
	go install -ldflags '$(LDFLAGS)' cmd/astronomer.go

lint:
	golangci-lint run --fix ./...

test:
	go test -coverprofile cover.out -v ./...

coverage:
	go tool cover -html=cover.out
