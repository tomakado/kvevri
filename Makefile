PROJECT_BIN = $(shell pwd)/bin

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
