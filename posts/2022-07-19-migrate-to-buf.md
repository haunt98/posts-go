# Migrate to `buf` from `prototool`

Why? Because `prototool` is outdated, and can not run on M1 mac.

We need 3 files:

- `build.go`: need to install protoc-gen-\* binaries with pin version in `go.mod`
- `buf.yaml`
- `buf.gen.yaml`

FYI, the libs version I use:

- [golang/protobuf v1.5.2](https://github.com/golang/protobuf/releases/tag/v1.5.2)
- [grpc-ecosystem/grpc-gateway v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.16.0)
- [envoyproxy/protoc-gen-validate](github.com/envoyproxy/protoc-gen-validate)
- [kei2100/protoc-gen-marshal-zap](github.com/kei2100/protoc-gen-marshal-zap)

`build.go`:

```go
//go:build tools
// +build tools

import (
  _ "github.com/envoyproxy/protoc-gen-validate"
  _ "github.com/golang/protobuf/protoc-gen-go"
  _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
  _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
  _ "github.com/kei2100/protoc-gen-marshal-zap/plugin/protoc-gen-marshal-zap"
)
```

`buf.yaml`

```yaml
version: v1
deps:
  - buf.build/haunt98/googleapis:b38d93f7ade94a698adff9576474ae7c
  - buf.build/haunt98/grpc-gateway:ecf4f0f58aa8496f8a76ed303c6e06c7
  - buf.build/haunt98/protoc-gen-validate:2686264610fc4ad4a9fcc932647e279d
  - buf.build/haunt98/marshal-zap:2a593ca925134680a5820d3f13c1be5a
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
  - name: validate
    out: pkg
    opt:
      - lang=go
  - name: marshal-zap
    out: pkg
```

Update `Makefile`:

```Makefile
gen:
  go install github.com/golang/protobuf/protoc-gen-go
  go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  go install github.com/envoyproxy/protoc-gen-validate
  go install github.com/kei2100/protoc-gen-marshal-zap/plugin/protoc-gen-marshal-zap
  go install github.com/bufbuild/buf/cmd/buf@latest
  buf mod update
  buf format -w
  buf generate
```

Run `make gen` to have fun of course.

## FAQ

Remember `grpc-ecosystem/grpc-gateway`, `envoyproxy/protoc-gen-validate`, `kei2100/protoc-gen-marshal-zap` is optional, so feel free to delete if you don't use theme.

If use `vendor`:

- Replace `buf generate` with `buf generate --exclude-path vendor`.
- Replace `buf format -w` with `buf format -w --exclude-path vendor`.

If you use grpc-gateway:

- Replace `import "third_party/googleapis/google/api/annotations.proto";` with `import "google/api/annotations.proto";`
- Delete `security_definitions`, `security`, in `option (grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger)`.

The last step is delete `prototool.yaml`.

If you are not migrate but start from scratch:

- Add `buf lint` to make sure your proto is good.
- Add `buf breaking --against "https://your-grpc-repo-goes-here.git"` to make sure each time you update proto, you don't break backward compatibility.

## Thanks

- [uber/prototool](https://github.com/uber/prototool)
- [bufbuild/buf](https://github.com/bufbuild/buf)
