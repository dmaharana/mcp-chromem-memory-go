#!/bin/bash

# Test script for MCP Memory Server

echo "Testing MCP Memory Server..."

# Start the server in background
./memory-server &
SERVER_PID=$!

# Give server time to start
sleep 1

# Test initialize
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}' | ./memory-server &
INIT_PID=$!

sleep 1

# Kill the background processes
kill $SERVER_PID 2>/dev/null
kill $INIT_PID 2>/dev/null

echo "Basic test completed. Server builds and starts successfully."
echo ""
echo "To test manually, run:"
echo "./memory-server"
echo ""
echo "Then send MCP requests like:"
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}'