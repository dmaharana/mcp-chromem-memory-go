# Local Memory Layer for Developers

A lightweight, dependency-free vector database that serves as an MCP (Model Context Protocol) server for knowledge management tasks. This implementation provides a solid foundation for storing and retrieving developer memories using statistical text embeddings without requiring external AI models.

## Features

- **File-based Vector Database**: Uses chromem-go for efficient vector storage
- **Statistical Embeddings**: Custom embedding algorithm using statistical methods (no LLM required)
- **Document Management**: Add, search, list, and delete memory documents
- **Tagging System**: Organize documents with tags for easy lookup
- **Favorites**: Mark important documents as favorites for higher search ranking
- **Key-Value Properties**: Store additional metadata with each document
- **MCP Server**: Compatible with Model Context Protocol for IDE integration
- **Zero Dependencies**: Single binary with no external dependencies

## Building

```bash
go build -o memory-server .
```

## Usage as MCP Server

The memory server implements the Model Context Protocol (MCP) and can be used as a stdio server:

```bash
./memory-server
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