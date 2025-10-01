#!/bin/bash

set -euo pipefail

DB_PATH="/tmp/test_memory.db"
HTTP_PORT=8080

rm -rf "${DB_PATH}"

# Kill any process already listening on the port
lsof -ti :${HTTP_PORT} | xargs -r kill -9

# Start the memory server in the background
./memory-server --db-path "${DB_PATH}" --http-port "${HTTP_PORT}" < /dev/null &
SERVER_PID=$!

# Give the server a moment to start
sleep 2

# Function to send MCP requests
send_mcp_request() {
  local method="$1"
  local params="$2"
  local id="$3"
  
  curl -X POST -H "Content-Type: application/json" \
    -d "{\"jsonrpc\": \"2.0\", \"method\": \"$method\", \"params\": $params, \"id\": \"$id\"}" \
    http://localhost:"${HTTP_PORT}"
}

# Test cases

# Add Albert Einstein
RESPONSE=$(send_mcp_request "tools/call" '{"name": "add_memory", "arguments": {"content": "Albert Einstein proposed the theory of relativity, which transformed our understanding of time, space, and gravity.", "tags": ["physics", "relativity", "einstein"], "properties": {"scientist": "Albert Einstein", "field": "physics"}}}' "1")
echo "Response: $RESPONSE"
if ! echo "$RESPONSE" | grep -q "Memory added successfully"; then
  echo "Test failed: Failed to add Albert Einstein"
  kill "${SERVER_PID}"
  exit 1
fi

# Add Marie Curie
RESPONSE=$(send_mcp_request "tools/call" '{"name": "add_memory", "arguments": {"content": "Marie Curie was a physicist and chemist who conducted pioneering research on radioactivity and won two Nobel Prizes.", "tags": ["physics", "chemistry", "radioactivity", "curie"], "properties": {"scientist": "Marie Curie", "field": "physics, chemistry"}}}' "2")
if ! echo "$RESPONSE" | grep -q "Memory added successfully"; then
  echo "Test failed: Failed to add Marie Curie"
  kill "${SERVER_PID}"
  exit 1
fi

# Add Isaac Newton
RESPONSE=$(send_mcp_request "tools/call" '{"name": "add_memory", "arguments": {"content": "Isaac Newton formulated the laws of motion and universal gravitation, laying the foundation for classical mechanics.", "tags": ["physics", "mechanics", "newton"], "properties": {"scientist": "Isaac Newton", "field": "physics"}}}' "3")
if ! echo "$RESPONSE" | grep -q "Memory added successfully"; then
  echo "Test failed: Failed to add Isaac Newton"
  kill "${SERVER_PID}"
  exit 1
fi

# Add Charles Darwin
RESPONSE=$(send_mcp_request "tools/call" '{"name": "add_memory", "arguments": {"content": "Charles Darwin introduced the theory of evolution by natural selection in his book 'On the Origin of Species'.", "tags": ["biology", "evolution", "darwin"], "properties": {"scientist": "Charles Darwin", "field": "biology", "book": "On the Origin of Species"}}}' "4")
if ! echo "$RESPONSE" | grep -q "Memory added successfully"; then
  echo "Test failed: Failed to add Charles Darwin"
  kill "${SERVER_PID}"
  exit 1
fi

# Add Ada Lovelace
RESPONSE=$(send_mcp_request "tools/call" '{"name": "add_memory", "arguments": {"content": "Ada Lovelace is regarded as the first computer programmer for her work on Charles Babbage's early mechanical computer, the Analytical Engine.", "tags": ["computer science", "programming", "lovelace"], "properties": {"scientist": "Ada Lovelace", "field": "computer science", "invention": "Analytical Engine"}}}' "5")
if ! echo "$RESPONSE" | grep -q "Memory added successfully"; then
  echo "Test failed: Failed to add Ada Lovelace"
  kill "${SERVER_PID}"
  exit 1
fi

# List all memories and verify count
RESPONSE=$(send_mcp_request "tools/call" '{"name": "list_memories", "arguments": {}}' "6")
if ! echo "$RESPONSE" | grep -q "Total 5 memories"; then
  echo "Test failed: Expected 5 memories, got different count"
  echo "Response: $RESPONSE"
  kill "${SERVER_PID}"
  exit 1
fi

echo "All add_memory tests passed!"

# Clean up
kill "${SERVER_PID}"
rm -f "${DB_PATH}"

exit 0