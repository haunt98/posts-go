# Dockerfile for Go

Each time I start new Go project, I repeat many steps.
Like set up `.gitignore`, CI configs, Dockerfile, ...

So I decide to have a baseline Dockerfile like this:

```Dockerfile
FROM golang:1.18beta1-bullseye as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY vendor .
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -o ./app main.go

FROM gcr.io/distroless/base-debian11

COPY --from=builder /build/app /app

ENTRYPOINT ["/app"]
```
