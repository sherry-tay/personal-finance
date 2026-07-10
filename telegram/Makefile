bin/telegram:
	go build -o ./bin/telegram main.go

setup:
	go mod download

.PHONY: test
test:
	go test ./... -race

.PHONY: lint
lint:
	golangci-lint run