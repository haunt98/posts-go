# Migrate to `buf` from `prototool`

Why? Because `prototool` is outdated, and can not run on M1 mac.

There are 2 usecase:

- For services only provide GRPC
- For services provide REST and GRPC, which make use of [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

We need 3 files:

- `build.go`: need to install protoc-gen-\* binaries with pin version in `go.mod`
- `buf.yaml`
- `buf.gen.yaml`

FYI, the libs version I use:

- [golang/protobuf v1.5.2](https://github.com/golang/protobuf/releases/tag/v1.5.2)
- [grpc-ecosystem/grpc-gateway v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.16.0)

## GRPC only

`build.go`:

```go
//go:build tools
// +build tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
)
```

`buf.yaml`:

```yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```

`buf.gen.yml`

```yaml
version: v1
plugins:
  - name: go
    out: pkg
    opt:
      - plugins=grpc
```

Update `Makefile`:

```Makefile
gen:
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/bufbuild/buf/cmd/buf@latest
	buf format -w
	buf generate
```

Then run `make gen` to check result.

## REST and GRPC using grpc-gateway

`build.go`:

```go
//go:build tools
// +build tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
)
```

`buf.yaml`

```yaml
version: v1
deps:
  - buf.build/haunt98/googleapis:b38d93f7ade94a698adff9576474ae7c
  - buf.build/haunt98/grpc-gateway:ecf4f0f58aa8496f8a76ed303c6e06c7
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```

`buf.gen.yaml`:

```yaml
version: v1
plugins:
  - name: go
    out: pkg
    opt:
      - plugins=grpc
  - name: grpc-gateway
    out: pkg
    opt:
      - logtostderr=true
  - name: swagger
    out: .
    opt:
      - logtostderr=true
```

Update `Makefile`:

```Makefile
gen:
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go install github.com/bufbuild/buf/cmd/buf@latest
	buf format -w
	buf mod update
	buf generate
```

Run `make gen` to get fun of course.

## FAQ

If use `vendor`:

- Replace `buf generate` with `buf generate --exclude-path vendor`.
- Replace `buf format -w` with `buf format -w --exclude-path vendor`.

The last step is delete `prototool.yaml`.

## Thanks

- [uber/prototool](https://github.com/uber/prototool)
- [bufbuild/buf](https://github.com/bufbuild/buf)
