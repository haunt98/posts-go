# Migrate to `buf` from `prototool`

Why? Because `prototool` is outdated, and can not run on M1 mac.

We need 3 files:

- `build.go`: need to install some binaries with pin version in `go.mod`
- `buf.yaml`
- `buf.gen.yaml`

FYI, I use:

- [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [bufbuild/protoc-gen-validate](github.com/bufbuild/protoc-gen-validate)
- [kei2100/protoc-gen-marshal-zap](github.com/kei2100/protoc-gen-marshal-zap)

## `build.go`

```go
//go:build tools
// +build tools

import (
  _ "github.com/kei2100/protoc-gen-marshal-zap/plugin/protoc-gen-marshal-zap"
)
```

## `buf.yaml`

```yaml
version: v2
deps:
    - buf.build/envoyproxy/protoc-gen-validate:v1.2.1
    - buf.build/kei2100/protoc-gen-marshal-zap:081f499bbca4486784773e060c1c1418
    - buf.build/haunt98/googleapis
    - buf.build/haunt98/grpc-gateway:v2.27.1
breaking:
    use:
        - WIRE
```

## `buf.gen.yaml`

```yaml
version: v2
plugins:
    - remote: buf.build/grpc/go:v1.5.1
      out: pkg
    - remote: buf.build/protocolbuffers/go:v1.36.6
      out: pkg
    - remote: buf.build/bufbuild/validate-go:v1.2.1
      out: pkg
    - local: protoc-gen-marshal-zap
      out: pkg
    - remote: buf.build/grpc-ecosystem/gateway:v2.27.1
      out: pkg
    - remote: buf.build/grpc-ecosystem/openapiv2:v2.27.1
      out: .
      opt:
          - json_names_for_fields=false
          - enums_as_ints=true
          - disable_service_tags=true
          - preserve_rpc_order=true
```

## `Makefile`

```Makefile
gen:
	# go install github.com/kei2100/protoc-gen-marshal-zap/plugin/protoc-gen-marshal-zap
	go install github.com/bufbuild/buf/cmd/buf@latest
	# buf config migrate
	buf dep update
	buf format --path ./api -w
	buf generate --path ./api
```

If using some imports, you better make a raw proto file for other services to integrate:

```Makefile
raw:
	mkdir -p api/raw
	cp api/*.proto api/raw/
	# https://github.com/chmln/sd
	sd -F 'import "validate/validate.proto";' '' api/raw/*.proto
	sd -F 'import "marshal-zap.proto";' '' api/raw/*.proto
	sd -f s '\s\[\s*\(.*?];' ';' api/raw/*.proto
	buf format --path api/raw -w
```

## FAQ

If you use grpc-gateway:

- Replace `third_party/googleapis/google/api/annotations.proto;` with `protoc-gen-openapiv2/options/annotations.proto`
- Replace `grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger` with
  `grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger`
- Delete `security_definitions`, `security`, in `option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)`.

The last step is delete `prototool.yaml`.

If you are not migrate but start from scratch:

- Add `buf lint` to make sure your proto is good.
- Add `buf breaking --against "https://your-grpc-repo-goes-here.git"` to make sure each time you update proto, you don't
  break backward compatibility.

## Tips

Some experience I got after writing proto files for a living:

- Delete `vendor`
- Ignore DRY (Do not Repeat Yourself) when handling proto, don't split proto into many files, just copy into 1 big
  proto. Trust me, it saves you from wasting time to debug how to import Go after generated. Because proto import and Go
  import is [2](https://github.com/golang/protobuf/issues/895) different things. If someone already have split proto
  files, you should use `sd` to fix the damn things.

## Thanks

- [uber/prototool](https://github.com/uber/prototool)
- [bufbuild/buf](https://github.com/bufbuild/buf)
