#!/bin/bash

# Test script to add sample memories via REST API

BASE_URL="http://localhost:8080"

echo "Adding sample memories to test the system..."
echo ""

# Test if server is running
if ! curl -s "$BASE_URL/api/stats" > /dev/null; then
    echo "Error: Server is not running at $BASE_URL"
    echo "Please start the server with: ./memory-server -web"
    exit 1
fi

echo "Server is running. Adding sample memories..."
echo ""

# Add memory 1
echo "Adding memory 1: Go channels..."
curl -X POST "$BASE_URL/api/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Go channels are used for communication between goroutines. Use make(chan Type) to create a channel.",
    "tags": ["golang", "concurrency", "channels"],
    "favorite": true,
    "properties": {"category": "concept", "difficulty": "intermediate"}
  }'
echo -e "\n"

# Add memory 2
echo "Adding memory 2: Null pointer debugging..."
curl -X POST "$BASE_URL/api/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "To fix null pointer exceptions in Go: Always check if pointer is nil before dereferencing. Use if ptr != nil { ... }",
    "tags": ["golang", "debugging", "error-handling"],
    "favorite": true,
    "properties": {"category": "bug-fix", "language": "go"}
  }'
echo -e "\n"

# Add memory 3
echo "Adding memory 3: Docker best practices..."
curl -X POST "$BASE_URL/api/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Docker best practices: Use multi-stage builds, minimize layers, use .dockerignore, and run as non-root user.",
    "tags": ["docker", "devops", "best-practices"],
    "favorite": false,
    "properties": {"category": "tip", "technology": "docker"}
  }'
echo -e "\n"

# Add memory 4
echo "Adding memory 4: Git workflow..."
curl -X POST "$BASE_URL/api/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Git workflow: Create feature branch, make changes, commit with descriptive messages, push, create PR, review, merge.",
    "tags": ["git", "workflow", "version-control"],
    "favorite": false,
    "properties": {"category": "process", "tool": "git"}
  }'
echo -e "\n"

# Add memory 5
echo "Adding memory 5: REST API design..."
curl -X POST "$BASE_URL/api/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "REST API design principles: Use HTTP methods correctly (GET, POST, PUT, DELETE), return appropriate status codes, use consistent naming conventions.",
    "tags": ["api", "rest", "design", "web-development"],
    "favorite": true,
    "properties": {"category": "design-pattern", "domain": "web"}
  }'
echo -e "\n"

echo "Sample memories added successfully!"
echo ""
echo "Test searches:"
echo ""

echo "1. Search for 'golang':"
curl -s "$BASE_URL/api/search?q=golang&limit=3" | jq -r '.[] | "- \(.content[0:80])..."'
echo ""

echo "2. Search for 'debugging':"
curl -s "$BASE_URL/api/search?q=debugging&limit=3" | jq -r '.[] | "- \(.content[0:80])..."'
echo ""

echo "3. Get current stats:"
curl -s "$BASE_URL/api/stats" | jq '.'
echo ""

echo "Visit http://localhost:8080 to see the web interface!"