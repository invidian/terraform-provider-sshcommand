# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Terraform parameters
PROVIDER_NAME=terraform-provider-sshcommand
PROVIDER_VERSION=`git describe --tags --always | cut -d"-" -f1`
PROTOCOL_VERSION=x4
BINARY_NAME=$(PROVIDER_NAME)_$(PROVIDER_VERSION)_$(PROTOCOL_VERSION)

# Output parameters
OUTPUT_DIRECTORY=`pwd`
OUTPUT_FILE=$(OUTPUT_DIRECTORY)/$(BINARY_NAME)

# Build parameters
CGO_ENABLED=0
LD_FLAGS="-extldflags '-static'"

all: test build lint

test:
	$(GOTEST) -v ./...

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -o $(OUTPUT_FILE) -v -buildmode=exe -ldflags $(LD_FLAGS)

lint:
	which golangci-lint 2>&1 >/dev/null && golangci-lint run || echo "'golangci-lint' binary not found, skipping linting."

clean:
	$(GOCLEAN)
	rm -f $(OUTPUT_FILE) || true
	rm -f $(OUTPUT_FILE).sig || true

update:
	$(GOGET) -u
	$(GOMOD) tidy

# TODO Add GitHub integration
release: all pack sign

pack:
	 which upx 2>&1 >/dev/null && upx --brute $(OUTPUT_FILE) || echo "'upx' binary not found, skipping packing."

sign:
	gpg --output $(OUTPUT_FILE).sig  --detach-sig $(OUTPUT_FILE)
