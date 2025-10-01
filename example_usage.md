# Example Usage

## Web Mode (Browser Interface)

Start the web server:
```bash
./memory-server -web
```

This will:
1. Start the web server on port 8080 (default)
2. Automatically open your browser to http://localhost:8080
3. Show a dashboard with statistics and document management interface

### Web Interface Features:
- **Dashboard**: View document counts and usage statistics
- **Add Memories**: Form to add new documents with tags, favorites, and properties
- **Search & Browse**: Search through memories or view all documents
- **Edit Documents**: Click "Edit" to modify existing memories
- **Favorite Management**: Toggle favorite status with star button
- **Real-time Updates**: Statistics update automatically

### Custom Web Server Options:
```bash
# Custom port and database
./memory-server -web -web-port 9090 -db-path my-memories.db

# Don't open browser automatically
./memory-server -web -open=false
```

## MCP Mode (Command Line)

Start the MCP server:
```bash
./memory-server
```

### 1. Initialize the MCP connection
```json
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}
```

### 2. List available tools
```json
{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}
```

### 3. Add a memory document
```json
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "add_memory", "arguments": {"content": "How to fix null pointer exception in Go: Always check if pointer is nil before dereferencing", "tags": ["golang", "debugging", "error-handling"], "favorite": true, "properties": {"category": "bug-fix", "language": "go"}}}}
```

### 4. Add another memory
```json
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "add_memory", "arguments": {"content": "Use defer statements for cleanup operations in Go functions", "tags": ["golang", "best-practices"], "properties": {"category": "tip"}}}}
```

### 5. Search for memories
```json
{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "search_memories", "arguments": {"query": "golang error handling", "limit": 5, "threshold": 0.1}}}
```

### 6. List all memories
```json
{"jsonrpc": "2.0", "id": 6, "method": "tools/call", "params": {"name": "list_memories", "arguments": {}}}
```

### 7. Delete a memory (use ID from previous responses)
```json
{"jsonrpc": "2.0", "id": 7, "method": "tools/call", "params": {"name": "delete_memory", "arguments": {"id": "your-document-id-here"}}}
```

## Expected Response Format

Each request will return a JSON-RPC 2.0 response with either a `result` or `error` field:

```json
{"jsonrpc": "2.0", "id": 1, "result": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}, "serverInfo": {"name": "memory-server", "version": "1.0.0"}}}
```

## Integration with Kiro IDE

To use with Kiro IDE, add to your MCP configuration:

```json
{
  "mcpServers": {
    "memory-server": {
      "command": "/path/to/memory-server",
      "args": [],
      "disabled": false,
      "autoApprove": ["add_memory", "search_memories", "list_memories"]
    }
  }
}
```