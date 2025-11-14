#!/bin/bash
set -e

echo "🚧 Building GoBaeBounty..."

# Check Go version
GO_VERSION=$(go version)
echo "Using $GO_VERSION"

# Create bin directory
mkdir -p ./bin

# Build main binary
go build -o ./bin/gobaebounty ./cmd/bbgolang

echo "✅ Build successful: ./bin/gobaebounty"
echo "You can execute it using: ./bin/gobaebounty --help"
