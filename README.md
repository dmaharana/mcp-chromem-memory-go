# Local Memory Layer for Developers

A lightweight, dependency-free vector database that serves as an MCP (Model Context Protocol) server for knowledge management tasks. This implementation provides a solid foundation for storing and retrieving developer memories using statistical text embeddings without requiring external AI models.

## Features

- **File-based Vector Database**: Uses chromem-go for efficient vector storage
- **Statistical Embeddings**: Custom embedding algorithm using statistical methods (no LLM required)
- **Document Management**: Add, search, list, and delete memory documents
- **Tagging System**: Organize documents with tags for easy lookup
- **Favorites**: Mark important documents as favorites for higher search ranking
- **Key-Value Properties**: Store additional metadata with each document
- **Web Interface**: Browser-based dashboard for managing memories
- **REST API**: Full REST endpoints for integration
- **Usage Statistics**: Track tool usage and document counts
- **MCP Server**: Compatible with Model Context Protocol for IDE integration
- **Zero Dependencies**: Single binary with no external dependencies

## Building

```bash
go build -o memory-server ./cmd/memory-server
```

## Usage Modes

### Web Mode (Browser Interface)

Start the server with a web interface that automatically opens in your browser:

```bash
./memory-server -web
```

Options:
- `-web-port 8080`: Set web server port (default: 8080)
- `-db-path memory.db`: Set database file path
- `-open=false`: Disable automatic browser opening

The web interface provides:
- **Dashboard**: View statistics and document counts
- **Add Memories**: Form to add new documents with tags, favorites, and properties
- **Search & Browse**: Search through memories or view all documents
- **Edit Documents**: Modify existing memories with automatic re-embedding
- **Favorite Management**: Mark/unmark documents as favorites
- **Real-time Stats**: Track usage of each operation

### MCP Server Mode

The memory server implements the Model Context Protocol (MCP) and can be used as a stdio server:

```bash
./memory-server
```

For MCP over HTTP:
```bash
./memory-server -http-port 3000
```

### Available MCP Tools

1. **add_memory**: Add a new memory document
   - `content` (required): The content of the memory document
   - `tags` (optional): Array of tags for the document
   - `favorite` (optional): Mark as favorite document
   - `properties` (optional): Additional key-value properties

2. **search_memories**: Search for memory documents
   - `query` (required): Search query string
   - `limit` (optional): Maximum number of results (default: 10)
   - `threshold` (optional): Similarity threshold 0.0-1.0 (default: 0.1)

3. **list_memories**: List all memory documents

4. **delete_memory**: Delete a memory document
   - `id` (required): Document ID to delete

## Statistical Embedding Algorithm

The custom embedding algorithm uses various statistical features:

- **Basic Statistics**: Text length, word count, lexical diversity
- **Linguistic Features**: Average word length, text entropy, readability scores
- **Character N-grams**: 2-gram and 3-gram frequency analysis
- **Word Features**: Word frequency and positional information
- **Normalization**: Vector normalization for consistent similarity matching

## Configuration

The server uses zerolog for structured logging with filename and line number information. Logs are output to stderr while MCP communication happens over stdout/stdin.

## Integration with IDEs

This server can be integrated with IDEs that support MCP, allowing developers to:

- Store code snippets, bug fixes, and solutions
- Search through past decisions and discussions
- Maintain context across projects and teams
- Build a personal knowledge base that grows with the codebase

## REST API Endpoints

When running in web mode, the following REST endpoints are available:

### Statistics
- `GET /api/stats` - Get server statistics and document counts

### Documents
- `GET /api/documents` - List all documents
- `POST /api/documents` - Add a new document
- `GET /api/documents/{id}` - Get a specific document
- `PUT /api/documents/{id}` - Update a document (triggers re-embedding)
- `DELETE /api/documents/{id}` - Delete a document
- `PUT /api/documents/{id}/favorite` - Toggle favorite status

### Search
- `GET /api/search?q={query}&limit={limit}&threshold={threshold}` - Search documents

### Example API Usage

```bash
# Add a new memory
curl -X POST http://localhost:8080/api/documents \
  -H "Content-Type: application/json" \
  -d '{
    "content": "How to fix null pointer in Go: always check if pointer is nil",
    "tags": ["golang", "debugging"],
    "favorite": true,
    "properties": {"category": "tip"}
  }'

# Search memories
curl "http://localhost:8080/api/search?q=golang%20debugging&limit=5"

# Get statistics
curl http://localhost:8080/api/stats
```

## Example MCP Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "memory-server": {
      "command": "./memory-server",
      "args": []
    }
  }
}
```

## Architecture

- `main.go`: Entry point and server initialization
- `memory_store.go`: Core memory storage and retrieval logic
- `embedder.go`: Statistical text embedding implementation
- `mcp_server.go`: MCP protocol implementation and tool handlers

The system stores documents in a chromem-go vector database with metadata including tags, favorites, creation dates, and custom properties. The statistical embedder creates meaningful similarity matching without requiring external AI models.