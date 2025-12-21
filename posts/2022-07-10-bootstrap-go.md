# Bootstrap Go

It is hard to write bootstrap tool to quickly create Go service. So I write this guide instead. This is a quick
checklist for me every damn time I need to write a Go service from scratch. Also, this is my personal opinion, so feel
free to comment.

## Structure

```txt
main.go
internal/
    http/
        handler.go
        service.go
        models.go
    grpc/
        handler.go
        models.go
    consumer/
        handler.go
        service.go
        models.go
    service.go
    repository.go
    models.go
```

All codes are inside `internal`. Because `internal` is magic keyword in Go, you can not import pkg inside `internal`.

There are 3 common handlers:

- `http/` is for public APIs (Android, iOS, ... are clients).
- `grpc/` is for internal APIs (other services are clients).
- `consumer/` is for consuming messages from queue (Kafka, RabbitMQ, ...).

For each handler, there are usually 3 layers: `handler`, `service`, `repository`:

- `handler` interacts directly with gRPC, REST or consumer using specific codes (cookies, ...) In case gRPC, there are
  frameworks outside handle for us so we can write logic codes here too. But remember, gRPC only.
- `service` is where we write logic codes, and only logic codes is written here.
- `repository` is where we write codes which interacts with database/cache like MySQL, Redis, ...
- `models` is where we put all request, response, data models.

Location:

- `handler` must exist inside `grpc`, `http`, `consumer`.
- `service`, `repository` can exist directly inside of `internal` if both `grpc`, `http`, `consumer` has same logic.

## Good taste

### Stop using global var

If I see someone using global var, I swear I will shoot them twice in the face.

Why?

- Can be edited on runtime
- Not thread safe
- Can not write unit test

### Avoid unsigned type

Just a var can not have negative value, doesn't mean it should use `uint`. Just use `int` and do not care about
boundary.

### Use functional options to init complex struct

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

In above example, I construct `s` with `WithA` and `WithB` option. No need to pass direct field inside `s`.

### Use [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup), [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)

If logic involves calling too many APIs, but they are not depend on each other. We can fire them parallel :)

Use `errgroup` if you need to cancel all tasks if first error is caught. Be super careful with `egCtx`, should use this
instead of parent `ctx` inside `eg.Go`.

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

### [Generics](https://go.dev/doc/tutorial/generics) with some tricks

Take value then return pointer, useful with database struct full of pointers:

```go
// Ptr takes in non-pointer and returns a pointer
func Ptr[T any](v T) *T {
	return &v
}
```

Return zero value:

```go
func Zero[T any]() T {
  var zero T
  return zero
}
```

### Modernize style

Sync go 1.26:

- Use `a := new(1)`

Sync go 1.24:

- Use `os.Root` instead of `os.Open`
- Use `omitzero` instead of `omitempty`

Sync go 1.22:

- Use `cmp.Or` as fallback mechanism

Since go 1.21:

- Use `slices.SortFunc` instead of `sort.Slice`.
- Use `ctx.WithoutCancel` to disconnect context from parent.
- Use `clear(m)` to clear map entirely.

Since go 1.20:

- Use `errors.Join` for multiple errors.

Since go 1.18:

- Use `any` instead of `interface{}`.

Use [gopls/modernize](https://pkg.go.dev/golang.org/x/tools/gopls/internal/analysis/modernize) to modernize code.

```sh
modernize -fix -test ./...
```

## External libs

### No need `go mod vendor`

Save storage space.

### Use [bufbuild/buf](https://github.com/bufbuild/buf) for proto related

Please don't use [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway):

- Can not upload files

### Use [gin-gonic/gin](https://github.com/gin-gonic/gin) for HTTP framework

With `c *gin.Context`:

- Don't use `c` when passing context, use `c.Request.Context()` instead.
- Don't use `c.Request.URL.Path`, use `c.FullPath()` instead.

Remember to free resources after parse multipart form:

```go
defer func() {
    if err := c.Request.MultipartForm.RemoveAll(); err != nil {
        // Handle error
    }
}()
```

Combine with [go-playground/validator](https://github.com/go-playground/validator) to validate request structs.

### Log with [uber-go/zap](https://github.com/uber-go/zap)

It is fast!

- Use `MarshalLogObject` when we need to hide some field of object when log (field is long or has sensitive value)
- Don't use `Panic`. Use `Fatal` for errors when start service to check dependencies. If you really need panic level,
  use `DPanic`.
- If doubt, use `zap.Any`.
- Add some kind of `trace_id` in every log lines for easily debug.

### Read config with [spf13/viper](https://github.com/spf13/viper)

Only init config in main or cmd layer. Do not use `viper.Get...` in inside layer.

Why?

- Hard to mock and test
- Put all config in single place for easily tracking

Also, be careful if config value is empty. You should decide to continue or stop the service if there is empty config.

### Connect database with [go-gorm/gorm](https://github.com/go-gorm/gorm)

Make sure to test your code (ORM or not) with [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock).

### Connect Redis with [redis/go-redis](https://github.com/redis/go-redis)

Be careful when use [HGETALL](https://redis.io/docs/latest/commands/hgetall//). If key not found, empty data will be
returned not nil error. See [redis/go-redis/issues/1668](https://github.com/redis/go-redis/issues/1668)

Use [Pipelines](https://redis.uptrace.dev/guide/go-redis-pipelines.html) for:

- HSET and EXPIRE in 1 command

Prefer to use `Pipelined` instead of `Pipeline`.

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

- `ReadTimeout`
- `WriteTimeout`
- `ContextTimeoutEnabled` to true
- `Protocol` to `3`
- `DisableIdentity` to true

### Connect MySQL with [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

Remember to config:

- `SetConnMaxLifetime`
- `SetMaxOpenConns`
- `SetMaxIdleConns`
- `ParseTime` to true.
- `Loc` to `time.UTC`.
- `CheckConnLiveness` to true.
- `ReadTimeout`, `WriteTimeout`

### Connect Kafka with [IBM/sarama](https://github.com/IBM/sarama)

Don't use [confluentinc/confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go), because it's required
`CGO_ENABLED`.

### Test with [stretchr/testify](https://github.com/stretchr/testify)

It is easy to write a suite test, thanks to testify. Also, for mocking, there are many options out there. Pick 1 then
sleep peacefully.

### Mock with [matryer/moq](https://github.com/matryer/moq) or [uber/mock](https://github.com/uber/mock)

The first is easy to use but not powerful as the later. If you want to make sure mock func is called with correct times,
use the later.

Example with `matryer/moq`:

```go
// Only gen mock if source code file is newer than mock file
// https://jonwillia.ms/2019/12/22/conditional-gomock-mockgen
//go:generate sh -c "test service_mock_generated.go -nt $GOFILE && exit 0; moq -rm -out service_mock_generated.go . Service"
```

### Replace `go fmt` with [mvdan/gofumpt](https://github.com/mvdan/gofumpt)

`gofumpt` provides more strict rules when format Go codes.

```sh
gofumpt -w -extra .
```

### Use [golangci/golangci-lint](https://github.com/golangci/golangci-lint)

No need to say more. Lint is the way!

My heuristic for fieldalignment (not work all the time): pointer -> string -> []byte -> int64 -> int32.

```sh
golangci-lint run --fix --no-config ./...
```

## Scripts

Change import:

```sh
gofmt -w -r '"github.com/Sirupsen/logrus" -> "github.com/sirupsen/logrus"' *.go
```

Cleanup if storage is full:

```sh
go clean -cache -testcache -modcache -fuzzcache -x
```

## Thanks

- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [TigerStyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md)
