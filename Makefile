# Build parameters
CGO_ENABLED=0
LD_FLAGS="-extldflags '-static'"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOBUILD=CGO_ENABLED=$(CGO_ENABLED) $(GOCMD) build -v -buildmode=exe -ldflags $(LD_FLAGS)
GO_PACKAGES=./...
GO_TESTS=^.*$

GOLANGCI_LINT_VERSION=v1.30.0
DISABLED_LINTERS=gci,goerr113

BIN_PATH=$$HOME/bin

.PHONY: all
all: build build-test test lint

.PHONY: download
download:
	$(GOMOD) download

.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_PATH) $(GOLANGCI_LINT_VERSION)

.PHONY: install-cc-test-reporter
install-cc-test-reporter:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > $(BIN_PATH)/cc-test-reporter
	chmod +x $(BIN_PATH)/cc-test-reporter

.PHONY: install-ci
install-ci: install-golangci-lint install-cc-test-reporter

.PHONY: build
build:
	$(GOBUILD)

.PHONY: test
test: build-test
	$(GOTEST) -run $(GO_TESTS) $(GO_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run --enable-all --disable=$(DISABLED_LINTERS) --max-same-issues=0 --max-issues-per-linter=0 --build-tags integration --timeout 10m --exclude-use-default=false $(GO_PACKAGES)

.PHONY: build-test
build-test:
	$(GOTEST) -run=nope $(GO_PACKAGES)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(OUTPUT_FILE) || true
	rm -f $(OUTPUT_FILE).sig || true

.PHONY: update
update:
	$(GOGET) -u $(GO_PACKAGES)
	$(GOMOD) tidy

.PHONY: all-cover
all-cover: build build-test test-cover lint

.PHONY: test-cover
test-cover: build-test
	$(GOTEST) -run $(GO_TESTS) -coverprofile=$(PROFILEFILE) $(GO_PACKAGES)

.PHONY: cover-upload
cover-upload: codecov
	# Make codeclimate as command, as we need to run test-cover twice and make deduplicates that.
	# Go test results are cached anyway, so it's fine to run it multiple times.
	make codeclimate

.PHONY: codecov
codecov: PROFILEFILE=coverage.txt
codecov: SHELL=/bin/bash
codecov: test-cover
codecov:
	bash <(curl -s https://codecov.io/bash)
