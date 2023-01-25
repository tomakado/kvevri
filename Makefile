PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint
GOLANGCI_LINT_VERSION = v1.50.1

# === Lint ===
.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) $(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

# === Build ===
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PROJECT_BIN)/kvevri ./cmd/

# === Run === 

.PHONY: run
run:
	go run cmd/main.go

# === Code generation from Protocol buffers ===
PROTOC = $(PROJECT_BIN)/protoc/bin/protoc
PROTOC_VERSION=3.15.8
PROTOC_OS=osx
PROTOC_URL=https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip

.PHONY: .install-protoc
.install-protoc:
	[ -f $(PROTOC) ] || rm -rf /tmp/protoc-$(PROTOC_VERSION) /tmp/protoc.zip && \
	wget -O /tmp/protoc.zip $(PROTOC_URL) && \
	mv /tmp/protoc.zip $(PROJECT_BIN) && \
	unzip -o $(PROJECT_BIN)/protoc.zip -d $(PROJECT_BIN)/protoc/
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: generate-protobuf
generate-protobuf: .install-protoc
	$(PROTOC) -I ./api/ --go_out=./internal/pb/ --go_opt=paths=source_relative \
    --go-grpc_out=./internal/pb/ --go-grpc_opt=paths=source_relative \
    api/store.proto


# === Test ===
.PHONY: test
test:
	go test -v --race --timeout=1m ./...

.PHONY: test-coverage
test-coverage:
	go test -v --timeout=5m --covermode=count --coverprofile=coverage.out ./...
	go tool cover --func=coverage.out
