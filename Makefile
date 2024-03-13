PROTOC := protoc # Path to the protoc compiler
PROTO_DIR := proto # Directory containing the .proto files
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto) # List of .proto files
PROTO_OUT_DIR := pkg/dna # Output directory for generated Go files

PACKAGES := $(shell go list ./...)

.PHONY: genproto
genproto:
    @if ! which $(PROTOC) > /dev/null; then \
        echo "Error: protoc not found. Please install Protocol Buffers."; \
        exit 1; \
    fi
    mkdir -p $(PROTO_OUT_DIR)
    $(PROTOC) --go_out=$(PROTO_OUT_DIR) $(PROTO_SRC)

.PHONY: build
build: genproto
    go build $(PACKAGES)

.PHONY: test
test:
    go test $(PACKAGES)

.PHONY: clean
clean:
    rm -rf $(PROTO_OUT_DIR) bin

.PHONY: install
install:
    go install $(PACKAGES)
