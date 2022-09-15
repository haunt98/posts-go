# Bootstrap Go

It is hard to write bootstrap tool to quickly create Go service.
So I write this guide instead.
This is a quick checklist for me every damn time I need to write a Go service from scratch.
Also, this is my personal opinion, so feel free to comment.

## Structure

```txt
main.go
internal
| business
| | http
| | | handler.go
| | | service.go
| | | models.go
| | grpc
| | | handler.go
| | | models.go
| | consumer
| | | handler.go
| | | service.go
| | | models.go
| | service.go
| | repository.go
| | models.go
```

All business codes are inside `internal`.
Each business has a different directory `business`.

Inside each business, there are 2 handlers: `http`, `grpc`:

- `http` is for public APIs (Android, iOS, ... are clients).
- `grpc` is for internal APIs (other services are clients).
- `consumer` is for consuming messages from queue (Kafka, RabbitMQ, ...).

For each handler, there are usually 3 layers: `handler`, `service`, `repository`:

- `handler` interacts directly with gRPC, REST or consumer using specific codes (cookies, ...) In case gRPC, there are frameworks outside handle for us so we can write business/logic codes here too. But remember, gRPC only.
- `service` is where we write business/logic codes, and only business/logic codes is written here.
- `repository` is where we write codes which interacts with database/cache like MySQL, Redis, ...
- `models` is where we put all request, response, data models.

Location:

- `handler` must exist inside `grpc`, `http`, `consumer`.
- `service`, `models` can exist directly inside of `business` if both `grpc`, `http`, `consumer` has same business/logic.
- `repository` should be placed directly inside of `business`.

## Do not repeat!

If we have too many services, some of the logic will be overlapped.

For example, service A and service B both need to make POST call API to service C.
If service A and service B both have libs to call service C to do that API, we need to move the libs to some common pkg libs.
So in the future, service D which needs to call C will not need to copy libs to handle service C api but only need to import from common pkg libs.

Another bad practice is adapter service.
No need to write a new service if what we need is just common pkg libs.

## Taste on style guide

### Stop using global var

If I see someone using global var, I swear I shoot twice in the face.

Why?

- Can not write unit test.
- Is not thread safe.

### Use functional options, but don't overuse it!

For simple struct with 1 or 2 fields, no need to use functional options.

[Example](https://go.dev/play/p/0XnOLiHuoz3):

```go
func main() {
	s := NewS(WithA(1), WithB("b"))
	fmt.Printf("%+v\n", s)
}

type S struct {
	fieldA int
	fieldB string
}

type OptionS func(s *S)

func WithA(a int) OptionS {
	return func(s *S) {
		s.fieldA = a
	}
}

func WithB(b string) OptionS {
	return func(s *S) {
		s.fieldB = b
	}
}

func NewS(opts ...OptionS) *S {
	s := &S{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
```

In above example, I construct `s` with `WithA` and `WithB` option.
No need to pass direct field inside `s`.

## External libs

### No need `vendor`

Only need if you need something from `vendor`, to generate mock or something else.

### Use `build.go` to include build tools in go.mod

To easily control version of build tools.

For example `build.go`:

```go
//go:build tools
// +build tools

package main

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
)
```

And then in `Makefile`:

```Makefile
build:
    go install github.com/golang/protobuf/protoc-gen-go
```

We always get the version of build tools in `go.mod` each time we install it.
Future contributors will not cry anymore.

### Don't use cli libs ([spf13/cobra](https://github.com/spf13/cobra), [urfave/cli](https://github.com/urfave/cli)) just for Go service

What is the point to pass many params (`do-it`, `--abc`, `--xyz`) when what we only need is start service?

In my case, service starts with only config, and config should be read from file or environment like [The Twelve Factors](https://12factor.net/) guide.

### Don't use [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

Just don't.

Use [protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go), [grpc/grpc-go](https://github.com/grpc/grpc-go) for gRPC.

Write 1 for both gRPC, REST sounds good, but in the end, it is not worth it.

### Don't use [uber/prototool](https://github.com/uber/prototool), use [bufbuild/buf](https://github.com/bufbuild/buf)

prototool is deprecated, and buf can generate, lint, format as good as prototool.

### Use [gin-gonic/gin](https://github.com/gin-gonic/gin) for REST.

Don't use `gin.Context` when pass context from handler layer to service layer, use `gin.Context.Request.Context()` instead.

### If you want log, just use [uber-go/zap](https://github.com/uber-go/zap)

It is fast!

- Don't overuse `func (*Logger) With`. Because if log line is too long, there is a possibility that we can lost it.

- Use `MarshalLogObject` when we need to hide some field of object when log (field is long or has sensitive value)

- Don't use `Panic`. Use `Fatal` for errors when start service to check dependencies. If you really need panic level, use `DPanic`.

- If doubt, use `zap.Any`.

- Use `contextID` or `traceID` in every log lines for easily debug.

### To read config, use [spf13/viper](https://github.com/spf13/viper)

Only init config in main or cmd layer.
Do not use `viper.Get...` in business layer or inside business layer.

Why?

- Hard to mock and test
- Put all config in single place for easily tracking

### Don't overuse ORM libs, no need to handle another layer above SQL.

Each ORM libs has each different syntax.
To learn and use those libs correctly is time consuming.
So just stick to plain SQL.
It is easier to debug when something is wrong.

But `database/sql` has its own limit.
For example, it is hard to get primary key after insert/update.
So may be you want to use ORM for those cases.

### If you want test, just use [stretchr/testify](https://github.com/stretchr/testify).

It is easy to write a suite test, thanks to testify.
Also, for mocking, there are many options out there.
Pick 1 then sleep peacefully.

### Replace `go fmt`, `goimports` with [mvdan/gofumpt](https://github.com/mvdan/gofumpt).

`gofumpt` provides more rules when format Go codes.

### Use [golangci/golangci-lint](https://github.com/golangci/golangci-lint).

No need to say more.
Lint or get the f out!

If you get `fieldalignment` error, use [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) to fix them.

```sh
# Install
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

# Fix
fieldalignment -fix ./internal/business/*.go
```

## Thanks

- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
