# Dockerfile for Go

Each time I start a new Go project, I repeat many steps.
Like set up `.gitignore`, CI configs, Dockerfile, ...

So I decide to have a baseline Dockerfile like this:

```Dockerfile
FROM golang:1.20-bullseye as builder

RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY vendor .
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -o ./app -tags timetzdata -trimpath -ldflags="-s -w" .

FROM gcr.io/distroless/base-debian11

COPY --from=builder /build/app /app

ENTRYPOINT ["/app"]
```

I use [multi-stage build](https://docs.docker.com/develop/develop-images/multistage-build/) to keep my image size small.
First stage is [Go official image](https://hub.docker.com/_/golang),
second stage is [Distroless](https://github.com/GoogleContainerTools/distroless).

Before Distroless, I use [Alpine official image](https://hub.docker.com/_/alpine),
There is a whole discussion on the Internet to choose which is the best base image for Go.
After reading some blogs, I discover Distroless as a small and secure base image.
So I stick with it for a while.

Also, remember to match Distroless Debian version with Go official image Debian version.

```Dockerfile
FROM golang:1.20-bullseye as builder
```

This is Go image I use as a build stage.
This can be official Go image or custom image is required in some companies.

```Dockerfile
RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download
```

This is optional.
In my case, my company is slow to update Go image so I use this trick to install latest Go version.

```Dockerfile
WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY vendor .
COPY . .
```

I use `/build` to emphasize that I am building something in that directory.

The 4 `COPY` lines are familiar if you use Go enough.
First is `go.mod` and `go.sum` because it defines Go modules.
The second is `vendor`, this is optional but I use it because I don't want each time I build Dockerfile, I need to redownload Go modules.

```Dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -o ./app -tags timetzdata -trimpath -ldflags="-s -w" .
```

This is where I build Go program.

`CGO_ENABLED=0` because I don't want to mess with C libraries.
`GOOS=linux GOARCH=amd64` is easy to explain, Linux with x86-64.
`GOAMD64=v3` is new since [Go 1.18](https://go.dev/doc/go1.18#amd64),
I use v3 because I read about AMD64 version in [Arch Linux rfcs](https://gitlab.archlinux.org/archlinux/rfcs/-/blob/master/rfcs/0002-march.rst). TLDR's newer computers are already x86-64-v3.

`-tags timetzdata` to embed timezone database in case base image does not have.
`-trimpath` to support reproduce build.
`-ldflags="-s -w"` to strip debugging information.

```Dockerfile
FROM gcr.io/distroless/base-debian11

COPY --from=builder /build/app /app

ENTRYPOINT ["/app"]
```

Finally, I copy `app` to Distroless base image.

## Thanks

- [How to start a Go project in 2023](https://boyter.org/posts/how-to-start-go-project-2023/)
- [Shrink your Go binaries with this one weird trick](https://words.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/)
