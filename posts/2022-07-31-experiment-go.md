# Experiment Go

There come a time when you need to experiment new things, new style, new approach.
So this post serves as it is named.

# Design API by trimming down the interface/struct or whatever

Instead of:

```go
type Client interface {
    GetUser()
    AddUser()
    GetAccount()
    RemoveAccount()
}

// c is Client
c.GetUser()
c.RemoveAccount()
```

Try:

```go
type Client struct {
    User ClientUser
    Account ClientAccount
}

type ClientUser interface {
    Get()
    Add()
}

type ClientAccount interface {
    Get()
    Remove()
}

// c is Client
c.User.Get()
c.Account.Remove()
```

The difference is `c.GetUser()` -> `c.User.Get()`.

For example we have client which connect to bank.
There are many functions like `GetUser`, `GetTransaction`, `VerifyAccount`, ...
So split big client to many children, each child handle single aspect, like user or transaction.

My concert is we replace an interface with a struct which contains multiple interfaces aka children.
I don't know if this is the right call.

This pattern is used by [google/go-github](https://github.com/google/go-github).

## Find alternative to [grpc/grpc-go](https://github.com/grpc/grpc-go)

Why?
[See for yourself](https://github.com/grpc/grpc-go/issues?q=is%3Aissue+compatibility+is%3Aclosed).

Also read:

- [A new Go API for Protocol Buffers](https://go.dev/blog/protobuf-apiv2) to know why `v1.20.0` is `v2`.
- [Go Protobuf Plugin Versioning](https://jbrandhorst.com/post/plugin-versioning/).

Currently there are some:

- [bufbuild/connect-go](https://github.com/bufbuild/connect-go). Comming from buf, trust worthy but need time to make it match feature parity with grpc-go.
- [twitchtv/twirp](https://github.com/twitchtv/twirp)
- [storj/drpc](https://github.com/storj/drpc)

# Thanks

- [API Clients for Humans](https://blog.gopheracademy.com/advent-2019/api-clients-humans/)
