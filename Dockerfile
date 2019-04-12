FROM golang:1.11-alpine

# Enable go modules
ENV GO111MODULE=on

# Install dependencies
RUN apk add curl git build-base

# Install linter
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $HOME/bin v1.16.0

# Copy go mod files first and install dependencies to cache this layer
ADD ./go.mod ./go.sum /go/src/terraform-provider-sshcommand/
WORKDIR /go/src/terraform-provider-sshcommand
RUN go get

# Add source code
ADD . /go/src/terraform-provider-sshcommand

# Build, test and lint
RUN go build -v && \
    go test && \
    $HOME/bin/golangci-lint run
