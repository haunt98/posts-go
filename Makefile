.PHONY: all test test-color coverage coverage-cli coverage-html lint format gen format-html

all: test-color lint format gen format-html
	go mod tidy

test:
	go test -race -failfast ./...

test-color:
	go install github.com/haunt98/go-test-color@latest
	go-test-color -race -failfast ./...

coverage:
	go test -coverprofile=coverage.out ./...

coverage-cli: coverage
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

format:
	go install github.com/haunt98/gofimports/cmd/gofimports@latest
	go install mvdan.cc/gofumpt@latest
	gofimports -w -company github.com/make-go-great .
	gofumpt -w -extra .

gen:
	go run .

format-html:
	yarn prettier --write .
