#!/bin/bash
set -e

echo "🚀 Starting demo of GoBaeBounty tool"

# Check auth file exists
AUTH_FILE="./configs/auth.example.json"
if [ ! -f "$AUTH_FILE" ]; then
  echo "Auth file missing. Please create $AUTH_FILE based on example."
  exit 1
fi

# Start a local simple test server (if implemented)
echo "Starting local test HTTP server on port 8080 (simulated)"
# For demo: use Python simple HTTP server here or a Go test server in background

# Run GoBaeBounty scan against localhost (assumes test server running)
./bin/gobaebounty --target localhost:8080 --auth-file "$AUTH_FILE" --o ./demo-results --workers 10 --max-rate 50 --depth 2 --v

echo "Demo complete. Check results in ./demo-results"
