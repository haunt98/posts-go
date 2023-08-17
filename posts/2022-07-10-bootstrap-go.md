# Bootstrap Go

It is hard to write bootstrap tool to quickly create Go service. So I write this
guide instead. This is a quick checklist for me every damn time I need to write
a Go service from scratch. Also, this is my personal opinion, so feel free to
comment.

## Structure

```txt
main.go
internal
    business
        http
            handler.go
            service.go
            models.go

        grpc
            handler.go
            models.go
        consumer
            handler.go
            service.go
            models.go
        service.go
        repository.go
        models.go
```

All business codes are inside `internal`. Each business has a different
directory `business`.

Inside each business, there are 2 handlers: `http`, `grpc`:

- `http` is for public APIs (Android, iOS, ... are clients).
- `grpc` is for internal APIs (other services are clients).
- `consumer` is for consuming messages from queue (Kafka, RabbitMQ, ...).

For each handler, there are usually 3 layers: `handler`, `service`,
`repository`:

- `handler` interacts directly with gRPC, REST or consumer using specific codes
  (cookies, ...) In case gRPC, there are frameworks outside handle for us so we
  can write business/logic codes here too. But remember, gRPC only.
- `service` is where we write business/logic codes, and only business/logic
  codes is written here.
- `repository` is where we write codes which interacts with database/cache like
  MySQL, Redis, ...
- `models` is where we put all request, response, data models.

Location:

- `handler` must exist inside `grpc`, `http`, `consumer`.
- `service`, `models` can exist directly inside of `business` if both `grpc`,
  `http`, `consumer` has same business/logic.
- `repository` should be placed directly inside of `business`.

## Do not repeat!

If we have too many services, some of the logic will be overlapped.

For example, service A and service B both need to make POST call API to service
C. If service A and service B both have libs to call service C to do that API,
we need to move the libs to some common pkg libs. So in the future, service D
which needs to call C will not need to copy libs to handle service C api but
only need to import from common pkg libs.

Another bad practice is adapter service. No need to write a new service if what
we need is just common pkg libs.

## Taste on style guide

### Stop using global var

If I see someone using global var, I swear I will shoot them twice in the face.

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

In above example, I construct `s` with `WithA` and `WithB` option. No need to
pass direct field inside `s`.

### Use [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) as much as possible

If business logic involves calling too many APIs, but they are not depend on
each other. We can fire them parallel :)

Personally, I prefer `errgroup` to `WaitGroup`
(https://pkg.go.dev/sync#WaitGroup). Because I always need deal with error. Be
super careful with `egCtx`, should use this instead of parent `ctx` inside
`eg.Go`.

Example:

```go
eg, egCtx := errgroup.WithContext(ctx)

eg.Go(func() error {
	// Do some thing
	return nil
})

eg.Go(func() error {
	// Do other thing
	return nil
})

if err := eg.Wait(); err != nil {
	// Handle error
}
```

### Use [semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore) when need to implement WorkerPool

Please don't use external libs for WorkerPool, I don't want to deal with
dependency hell.

### Use [sync.Pool](https://pkg.go.dev/sync#Pool) when need to reuse object, mainly for `bytes.Buffer`

Example:

```go
var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func MarshalWithoutEscapeHTML(v any) ([]byte, error) {
	b, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		return nil, ErrBufPoolNotBytesBuffer
	}

	b.Reset()
	defer bufPool.Put(b)

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	result := make([]byte, b.Len())
	copy(result, b.Bytes())
	return result, nil
}
```

## External libs

### No need `vendor`

Only need if you need something from `vendor`, to generate mock or something
else.

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

What is the point to pass many params (`do-it`, `--abc`, `--xyz`) when what we
only need is start service?

In my case, service starts with only config, and config should be read from file
or environment like [The Twelve Factors](https://12factor.net/) guide.

### Don't use [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

Just don't.

Use
[protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go),
[grpc/grpc-go](https://github.com/grpc/grpc-go) for gRPC.

Write 1 for both gRPC, REST sounds good, but in the end, it is not worth it.

### Don't use [uber/prototool](https://github.com/uber/prototool), use [bufbuild/buf](https://github.com/bufbuild/buf)

prototool is deprecated, and buf can generate, lint, format as good as
prototool.

### Use [gin-gonic/gin](https://github.com/gin-gonic/gin) for REST.

With `c *gin.Context`:

- Don't use `c` when passing context, use `c.Request.Context()` instead.
- Don't use `c.Request.URL.Path`, use `c.FullPath()` instead.

Remember to free resources after parse multipart form:

```go
defer func() {
    if err := c.Request.MultipartForm.RemoveAll(); err != nil {
        fmt.Println(err)
    }
}()
```

### If you want log, just use [uber-go/zap](https://github.com/uber-go/zap)

It is fast!

- Don't overuse `func (*Logger) With`. Because if log line is too long, there is
  a possibility that we can lost it.
- Use `MarshalLogObject` when we need to hide some field of object when log
  (field is long or has sensitive value)
- Don't use `Panic`. Use `Fatal` for errors when start service to check
  dependencies. If you really need panic level, use `DPanic`.
- If doubt, use `zap.Any`.
- Use `contextID` or `traceID` in every log lines for easily debug.

### To read config, use [spf13/viper](https://github.com/spf13/viper)

Only init config in main or cmd layer. Do not use `viper.Get...` in business
layer or inside business layer.

Why?

- Hard to mock and test
- Put all config in single place for easily tracking

Also, be careful if config value is empty. You should decide to continue or stop
the service if there is empty config.

### Don't overuse ORM libs, no need to handle another layer above SQL.

Each ORM libs has each different syntax. To learn and use those libs correctly
is time consuming. So just stick to plain SQL. It is easier to debug when
something is wrong.

Also please use
[prepared statement](https://go.dev/doc/database/prepared-statements) as much as
possible. Idealy, we should init all prepared statement when we init database
connection to cached it, not create it every time we need it.

But `database/sql` has its own limit. For example, it is hard to get primary key
after insert/update. So may be you want to use ORM for those cases. I hear that
[go-gorm/gorm](https://github.com/go-gorm/gorm),
[ent/ent](https://github.com/ent/ent) is good.

### Connect Redis with [redis/go-redis](https://github.com/redis/go-redis)

Be careful when use [HGETALL](https://redis.io/commands/hgetall/). If key not
found, empty data will be returned not nil error. See
[redis/go-redis/issues/1668](https://github.com/redis/go-redis/issues/1668)

Use [Pipelines](https://redis.uptrace.dev/guide/go-redis-pipelines.html) for:

- HSET and EXPIRE in 1 command.
- Multiple GET in 1 command.

Prefer to use `Pipelined` instead of `Pipeline`. Inside `Pipelined`, please
return `redis.Cmder` for each command.

Example:

```go
func (c *client) HSetWithExpire(ctx context.Context, key string, values []any, expired time.Duration) error {
	cmds := make([]redis.Cmder, 2)

	if _, err := c.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		cmds[0] = pipe.HSet(ctx, key, values...)

		if expired > 0 {
			cmds[1] = pipe.Expire(ctx, key, expired)
		}

		return nil
	}); err != nil {
		return err
	}

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		if err := cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}
```

Remember to config:

- `ReadTimeout`, `WriteTimeout`

### Connect MySQL with [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

Remember to config:

- `SetConnMaxLifetime`
- `SetMaxOpenConns`
- `SetMaxIdleConns`
- `ParseTime` to true.
- `Loc` to `time.UTC`.
- `CheckConnLiveness` to true.
- `ReadTimeout`, `WriteTimeout`

### Connect SQLite with [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)

Remember to config:

- Write-Ahead Logging: `PRAGMA journal_mode=WAL`
- Disable connections pool with `SetMaxOpenConns` sets to 1

Don't use [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3), it's required
`CGO_ENABLED`.

### Connect Kafka with [Shopify/sarama](https://github.com/Shopify/sarama)

Don't use
[confluentinc/confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go),
it's required `CGO_ENABLED`.

### If you want test, just use [stretchr/testify](https://github.com/stretchr/testify).

It is easy to write a suite test, thanks to testify. Also, for mocking, there
are many options out there. Pick 1 then sleep peacefully.

### If need to mock, choose [matryer/moq](https://github.com/matryer/moq) or [uber/mock](https://github.com/uber/mock)

The first is easy to use but not powerful as the later. If you want to make sure
mock func is called with correct times, use the later.

Example with `matryer/moq`:

```go
// Only gen mock if source code file is newer than mock file
// https://jonwillia.ms/2019/12/22/conditional-gomock-mockgen
//go:generate sh -c "test service_mock_generated.go -nt $GOFILE && exit 0; moq -rm -out service_mock_generated.go . Service"
```

### Be careful with [spf13/cast](https://github.com/spf13/cast)

Don't cast proto enum:

```go
// Bad
a := cast.ToInt32(servicev1.ReasonCode_ABC)

// Good
a := int32(servicev1.ReasonCode_ABC)
```

### Use [stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) if you want your type enum can be print as string

```go
type Drink int

const (
	Beer Drink = iota
	Water
	OrangeJuice
)
```

```sh
go install golang.org/x/tools/cmd/stringer@latest

# Run inside directory which contains Drink
stringer -type=Drink
```

### Don't waste your time rewrite rate limiter if your use case is simple, use [rate](https://pkg.go.dev/golang.org/x/time/rate) or [go-redis/redis_rate](https://github.com/go-redis/redis_rate)

**rate** if you want rate limiter locally in your single instance of service.
**redis_rate** if you want rate limiter distributed across all your instances of
service.

### Replace `go fmt`, `goimports` with [mvdan/gofumpt](https://github.com/mvdan/gofumpt).

`gofumpt` provides more rules when format Go codes.

### Use [golangci/golangci-lint](https://github.com/golangci/golangci-lint).

No need to say more. Lint or get the f out!

If you get `fieldalignment` error, use
[fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment)
to fix them.

```sh
# Install
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

# Fix
fieldalignment -fix ./internal/business/*.go
```

## Snippet/script

Change import:

```sh
gofmt -w -r '"github.com/Sirupsen/logrus" -> "github.com/sirupsen/logrus"' *.go
```

Cleanup if storage is full:

```sh
go clean -cache -testcache -modcache -fuzzcache -x
```

## Thanks

- [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [Designing Go Libraries: The Talk: The Article](https://abhinavg.net/2022/12/06/designing-go-libraries/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style](https://google.github.io/styleguide/go/index)
- [Three bugs in the Go MySQL Driver](https://github.blog/2020-05-20-three-bugs-in-the-go-mysql-driver/)
- [Fixing Memory Exhaustion Bugs in My Golang Web App](https://mtlynch.io/notes/picoshare-perf/)
- [Prevent Logging Secrets in Go by Using Custom Types](https://www.commonfate.io/blog/prevent-logging-secrets-in-go-by-using-custom-types)
- [Speed Up GoMock with Conditional Generation](https://jonwillia.ms/2019/12/22/conditional-gomock-mockgen)
- [Making SQLite faster in Go](https://turriate.com/articles/making-sqlite-faster-in-go)
