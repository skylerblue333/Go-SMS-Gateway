.PHONY: fmt vet test race build lint check

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

race:
	go test -race ./...

build:
	CGO_ENABLED=0 go build -trimpath -o bin/sky-sms-envelope .

lint:
	golangci-lint run

check: vet test build
