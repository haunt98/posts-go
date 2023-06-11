# Migrate to `buf` from `prototool`

Why? Because `prototool` is outdated, and can not run on M1 mac.

We need 3 files:

- `build.go`: need to install protoc-gen-\* binaries with pin version in `go.mod`
- `buf.yaml`
- `buf.gen.yaml`

FYI, I use:

- [grpc-ecosystem/grpc-gateway v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.16.0)
- [bufbuild/protoc-gen-validate](github.com/bufbuild/protoc-gen-validate)
- [kei2100/protoc-gen-marshal-zap](github.com/kei2100/protoc-gen-marshal-zap)

`build.go`:

```go
//go:build tools
// +build tools

import (
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
  - buf.build/envoyproxy/protoc-gen-validate:6607b10f00ed4a3d98f906807131c44a
  - buf.build/kei2100/protoc-gen-marshal-zap:081f499bbca4486784773e060c1c1418
  - buf.build/haunt98/googleapis:b38d93f7ade94a698adff9576474ae7c
  - buf.build/haunt98/grpc-gateway:ecf4f0f58aa8496f8a76ed303c6e06c7
breaking:
  use:
    - PACKAGE
lint:
  use:
    - DEFAULT
```

`buf.gen.yaml`:

```yaml
version: v1
plugins:
  - plugin: buf.build/grpc/go:v1.3.0
    out: pkg
  - plugin: buf.build/protocolbuffers/go:v1.30.0
    out: pkg
  - plugin: buf.build/bufbuild/validate-go:v1.0.1
    out: pkg
    opt:
      - lang=go
  - plugin: marshal-zap
    out: pkg
  - plugin: grpc-gateway
    out: pkg
    opt:
      - logtostderr=true
  - plugin: swagger
    out: .
    opt:
      - logtostderr=true
```

Update `Makefile`:

```Makefile
gen:
  go install github.com/kei2100/protoc-gen-marshal-zap/plugin/protoc-gen-marshal-zap
  go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  go install github.com/bufbuild/buf/cmd/buf@latest
  buf mod update
  buf format -w
  buf generate
```

Run `make gen` to have fun of course.

If using `bufbuild/protoc-gen-validate`, `kei2100/protoc-gen-marshal-zap`, better make a raw copy of proto file for other services to integrate:

```Makefile
raw:
    mkdir -p ./raw
    cp ./api.proto ./raw/
    sed -i "" -e "s/import \"validate\/validate\.proto\";//g" ./raw/api.proto
    sed -i "" -e "s/\[(validate\.rules)\.string.min_len = 1\]//g" ./raw/api.proto
    sed -i "" -e "s/import \"marshal-zap\.proto\";//g" ./raw/api.proto
    sed -i "" -e "s/\[(marshal_zap\.mask) = true]//g" ./raw/api.proto
```

## FAQ

Remember `bufbuild/protoc-gen-validate`, `kei2100/protoc-gen-marshal-zap`, `grpc-ecosystem/grpc-gateway` is optional, so feel free to delete if you don't use theme.

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

# Tips

Some experience I got after writing proto files for a living:

- Ignore DRY (Do not Repeat Yourself) when handling proto, don't split proto into many files.
  Trust me, it saves you from wasting time to debug how to import Go after generated.
  Because proto import and Go import is [2](https://github.com/golang/protobuf/issues/895) different things.
  If someone already have split proto files, you should use `sed` to fix the damn things.

## Thanks

- [uber/prototool](https://github.com/uber/prototool)
- [bufbuild/buf](https://github.com/bufbuild/buf)
