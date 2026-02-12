# Use [bufbuild/buf](https://github.com/bufbuild/buf)

We need 3 files:

- `build.go`: install some binaries with pin version in `go.mod`
- `buf.yaml`
- `buf.gen.yaml`

FYI, I use:

- [bufbuild/protoc-gen-validate](github.com/bufbuild/protoc-gen-validate)
- [kei2100/protoc-gen-marshal-zap](github.com/kei2100/protoc-gen-marshal-zap)
- [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

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
    - buf.build/envoyproxy/protoc-gen-validate:v1.3.0
    - buf.build/kei2100/protoc-gen-marshal-zap:081f499bbca4486784773e060c1c1418
    - buf.build/haunt98/googleapis
    - buf.build/haunt98/grpc-gateway:v2.27.2
breaking:
    use:
        - WIRE
```

## `buf.gen.yaml`

```yaml
version: v2
plugins:
    - remote: buf.build/grpc/go:v1.6.1
      out: pkg
    - remote: buf.build/protocolbuffers/go:v1.36.11
      out: pkg
    - remote: buf.build/bufbuild/validate-go:v1.3.0
      out: pkg
    - local: protoc-gen-marshal-zap
      out: pkg
    - remote: buf.build/grpc-ecosystem/gateway:v2.27.2
      out: pkg
    - remote: buf.build/grpc-ecosystem/openapiv2:v2.27.2
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

For easily import, you better make a raw version of proto:

```Makefile
raw:
	mkdir -p ./api/raw
	cp ./api/*.proto ./api/raw/
	# https://github.com/chmln/sd
	sd -F 'import "validate/validate.proto";' '' ./api/raw/*.proto
	sd -F 'import "marshal-zap.proto";' '' ./api/raw/*.proto
	sd -f s '\s\[\s*\(.*?];' ';' ./api/raw/*.proto
	buf format --path ./api/raw -w
```

## FAQ

Some extra tips:

- Add `buf lint` to make sure your proto is good.
- Add `buf breaking --against "https://your-grpc-repo-goes-here.git"` to make sure each time you update proto, you don't
  break backward compatibility.

## Tips

Some experience I got after writing proto files for a living:

- Don't split proto into many files, just copy into 1 big proto. Trust me, it saves you from wasting time to debug how
  to import Go after generated. Because proto import and Go import is
  [2 different things](https://github.com/golang/protobuf/issues/895). If someone already have split proto files, you
  should merge yourself.
