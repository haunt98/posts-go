# Bootstrap Go

It is hard to write bootstrap tool to quickly create Go service.
So I write this guide instead.
This is a quick checklist for me every damn time I need to write a Go service from scratch.
Also, this is my personal opinion, so feel free to comment.

## Structure

```txt
main.go
internal
| business_1
| | http
| | | handler.go
| | | service.go
| | | repository.go
| | | models.go
| | grpc
| | | handler.go
| | | service.go
| | | repository.go
| | | models.go
| | service.go
| | repository.go
| | models.go
| business_2
| | grpc
| | | handler.go
| | | service.go
| | | repository.go
| | | models.go
```

All business codes are inside `internal`.
Each business has a different directory (`business_1`, `business_2`).

Inside each business, there are 2 handlers: `http`, `grpc`:

- `http` is for public APIs (Android, iOS,... are clients).
- `grpc` is for internal APIs (other services are clients).

Inside each handler, there are usually 3 layers: `handler`, `service`, `repository`:

- `handler` interacts directly with gRPC or REST using specific codes (cookies,...)
- `service` is where we write business/logic codes, and only business/logic codes is written here.
- `repository` is where we write codes which interacts with database/cache like MySQL, Redis, ...

`handler` must exist inside `grpc`, `http`.
But `service`, `repository`, `models` can exist directly inside `business` if both `grpc`, `http` has same business/logic.

## Do not repeat!

If we have too many services, some of the logic will be overlapped.

For example, service A and service B both need to make POST call API to service C.
If service A and service B both have libs to call service C to do that API, we need to move the libs to some common pkg libs.
So in the future, service D which needs to call C will not need to copy libs to handle service C api but only need to import from common pkg libs.

Another bad practice is adapter service.
No need to write a new service if what we need is just common pkg libs.

## External libs

### Don't use cli libs ([spf13/cobra](https://github.com/spf13/cobra), [urfave/cli](https://github.com/urfave/cli)) just for Go service

What is the point to pass many params (`--abc`, `--xyz`) when what we only need is start service?

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

- Use `MarshalLogObject` when we need to hide some field of object when log (field has long or sensitive value)

- Don't use `Panic`. Use `Fatal` for errors when start service to check dependencies. If you really need panic level, use `DPanic`.

- Use `contextID` or `traceID` in every log lines for easily debug.

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
