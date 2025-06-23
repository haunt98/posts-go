all: tidy format test-color lint gen format-html

tidy:
    go mod tidy

test:
    go test -race -failfast ./...

test-color:
    # go install github.com/haunt98/go-test-color@latest
    go-test-color -race -failfast ./...

coverage:
    go test -coverprofile=coverage.out ./...

coverage-cli: coverage
    go tool cover -func=coverage.out

coverage-html: coverage
    go tool cover -html=coverage.out

lint:
    golangci-lint run --fix ./...
    modernize -fix -test ./...

format:
    # go install github.com/haunt98/gofimports/cmd/gofimports@latest
    # go install mvdan.cc/gofumpt@latest
    gofimports -w --company github.com/make-go-great,github.com/haunt98 .
    gofumpt -w -extra .

gen:
    go run .

format-html:
    bunx prettier --write ./templates ./docs
    bunx prettier --print-width 120 --tab-width 4 --prose-wrap always --write ./posts

srht:
    # https://srht.site/quickstart
    # go install git.sr.ht/~xenrox/hut@latest
    # hut init
    tar -C docs -cvz . > site.tar.gz
    hut pages publish -d youngyoshie.srht.site site.tar.gz
    rm -rf site.tar.gz
