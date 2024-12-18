#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Resolve script's absolute path
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Debug output
echo "Script Location: $SCRIPT_DIR"
echo "Project Root: $PROJECT_ROOT"

# Create output directories if they don't exist
PROTO_DIR="$PROJECT_ROOT/api/proto"
GO_OUT_DIR="$PROJECT_ROOT/pkg/pb"

# Check if directory exists and create it if not
if [ ! -d "$GO_OUT_DIR" ]; then
  mkdir -p "$GO_OUT_DIR"
fi

# Check if protoc and protoc-gen-go are installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc (Protocol Buffers compiler) is not installed"
    exit 1
fi


if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not installed"
    exit 1
fi

# Generate Go code for each .proto file
# Use find to support multiple directories and nested proto files
find "$PROTO_DIR" -name "*.proto" | while read -r proto_file; do
    echo "Generating Go code for $proto_file"
    
    # Generate gRPC and Protocol Buffers Go code
    protoc \
        -I"$PROTO_DIR" \
        --go_out="$GO_OUT_DIR" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$GO_OUT_DIR" \
        --go-grpc_opt=paths=source_relative \
        "$proto_file"
done

echo "Proto generation complete!"