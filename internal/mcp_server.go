package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPServer struct {
	store  *MemoryStore
	server *mcp.Server
}

func NewMCPServer(store *MemoryStore) *MCPServer {
	s := &MCPServer{
		store: store,
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "memory-server"}, nil)

	type addMemoryArgs struct {
		Content    string   `json:"content" jsonschema:"the content of the memory document"`
		Tags       []string `json:"tags,omitempty" jsonschema:"Tags for the document"`
		Favorite   bool     `json:"favorite,omitempty" jsonschema:"Mark as favorite document"`
		Properties map[string]string `json:"properties,omitempty" jsonschema:"Additional key-value properties"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_memory",
		Description: "Add a new memory document to the store",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addMemoryArgs) (*mcp.CallToolResult, any, error) {
		doc := Document{
			ID:         uuid.New().String(),
			Content:    args.Content,
			CreatedAt:  time.Now(),
			Tags:       args.Tags,
			Favorite:   args.Favorite,
			Properties: args.Properties,
		}
		if err := s.store.AddDocument(doc); err != nil {
			return nil, nil, fmt.Errorf("failed to add document: %w", err)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Memory added successfully with ID: %s", doc.ID)},
			},
		}, nil, nil
	})

	type searchMemoriesArgs struct {
		Query     string  `json:"query" jsonschema:"Search query"`
		Limit     int     `json:"limit,omitempty" jsonschema:"Maximum number of results"`
		Threshold float32 `json:"threshold,omitempty" jsonschema:"Similarity threshold (0.0-1.0)"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_memories",
		Description: "Search for memory documents based on query",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchMemoriesArgs) (*mcp.CallToolResult, any, error) {
		if args.Limit == 0 {
			args.Limit = 10
		}
		if args.Threshold == 0 {
			args.Threshold = 0.1
		}
		docs, err := s.store.SearchDocuments(args.Query, args.Limit, args.Threshold)
		if err != nil {
			return nil, nil, fmt.Errorf("search failed: %w", err)
		}

		var results []string
		for i, doc := range docs {
			tags := strings.Join(doc.Tags, ", ")
			favorite := ""
			if doc.Favorite {
				favorite = " ⭐"
			}
			result := fmt.Sprintf("%d. [%s]%s\nContent: %s\nTags: %s\nCreated: %s\n",
				i+1, doc.ID, favorite, doc.Content, tags, doc.CreatedAt.Format("2006-01-02 15:04:05"))
			results = append(results, result)
		}
		responseText := fmt.Sprintf("Found %d memories:\n\n%s", len(docs), strings.Join(results, "\n"))

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	type listMemoriesArgs struct{}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_memories",
		Description: "List all memory documents",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listMemoriesArgs) (*mcp.CallToolResult, any, error) {
		docs, err := s.store.ListDocuments()
		if err != nil {
			return nil, nil, fmt.Errorf("list failed: %w", err)
		}

		var results []string
		for i, doc := range docs {
			tags := strings.Join(doc.Tags, ", ")
			favorite := ""
			if doc.Favorite {
				favorite = " ⭐"
			}
			result := fmt.Sprintf("%d. [%s]%s\nContent: %s\nTags: %s\nCreated: %s\n",
				i+1, doc.ID, favorite, doc.Content, tags, doc.CreatedAt.Format("2006-01-02 15:04:05"))
			results = append(results, result)
		}
		responseText := fmt.Sprintf("Total %d memories:\n\n%s", len(docs), strings.Join(results, "\n"))

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	type deleteMemoryArgs struct {
		ID string `json:"id" jsonschema:"Document ID to delete"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_memory",
		Description: "Delete a memory document by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteMemoryArgs) (*mcp.CallToolResult, any, error) {
		if err := s.store.DeleteDocument(args.ID); err != nil {
			return nil, nil, fmt.Errorf("delete failed: %w", err)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Memory with ID %s deleted successfully", args.ID)},
			},
		}, nil, nil
	})

	s.server = server
	return s
}

func (s *MCPServer) Start() error {
	return s.server.Run(context.Background(), &mcp.StdioTransport{})
}
