FROM docker.io/library/golang:1.19.2 as builder
LABEL org.opencontainers.image.authors=moulickaggarwal

WORKDIR /app
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download -x

# Copy the go source
COPY cmd/ cmd/
COPY validate/ validate/
COPY main.go main.go

# Build
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o k8s-cr-validator main.go

FROM alpine:latest
RUN apk add git --no-cache
WORKDIR /
COPY --from=builder /app/k8s-cr-validator .
