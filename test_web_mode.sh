#!/bin/bash

# Test script for Web Mode functionality

echo "Building memory server..."
go build -o memory-server ./cmd/memory-server

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"
echo ""
echo "Testing Web Mode:"
echo "1. Start web server: ./memory-server -web"
echo "2. Open browser to: http://localhost:8080"
echo "3. Test REST API endpoints:"
echo ""

echo "Example API tests:"
echo ""

echo "# Get initial stats"
echo "curl http://localhost:8080/api/stats"
echo ""

echo "# Add a test memory"
echo 'curl -X POST http://localhost:8080/api/documents \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '"'"'{'
echo '    "content": "Go channels are used for communication between goroutines",'
echo '    "tags": ["golang", "concurrency", "channels"],'
echo '    "favorite": true,'
echo '    "properties": {"category": "concept", "difficulty": "intermediate"}'
echo '  }'"'"
echo ""

echo "# Search for memories"
echo "curl \"http://localhost:8080/api/search?q=golang%20channels&limit=5\""
echo ""

echo "# List all documents"
echo "curl http://localhost:8080/api/documents"
echo ""

echo "To start the web server with custom options:"
echo "./memory-server -web -web-port 9090 -db-path custom.db -open=false"
echo ""

echo "For MCP mode (stdio):"
echo "./memory-server"
echo ""

echo "For MCP over HTTP:"
echo "./memory-server -http-port 3000"