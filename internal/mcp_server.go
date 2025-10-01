package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MCPServer struct {
	store *MemoryStore
}

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Tools map[string]interface{} `json:"tools"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

func NewMCPServer(store *MemoryStore) *MCPServer {
	return &MCPServer{
		store: store,
	}
}

func (s *MCPServer) Start() error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		log.Debug().Str("request", line).Msg("Received MCP request")
		
		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Error().Err(err).Msg("Failed to parse request")
			s.sendError(nil, -32700, "Parse error")
			continue
		}
		
		s.handleRequest(req)
	}
	
	return scanner.Err()
}

func (s *MCPServer) handleRequest(req MCPRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	default:
		s.sendError(req.ID, -32601, "Method not found")
	}
}

func (s *MCPServer) handleInitialize(req MCPRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": ServerCapabilities{
			Tools: map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "memory-server",
			"version": "1.0.0",
		},
	}
	
	s.sendResponse(req.ID, result)
}

func (s *MCPServer) handleToolsList(req MCPRequest) {
	tools := []Tool{
		{
			Name:        "add_memory",
			Description: "Add a new memory document to the store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The content of the memory document",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Tags for the document",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"favorite": map[string]interface{}{
						"type":        "boolean",
						"description": "Mark as favorite document",
					},
					"properties": map[string]interface{}{
						"type":        "object",
						"description": "Additional key-value properties",
					},
				},
				Required: []string{"content"},
			},
		},
		{
			Name:        "search_memories",
			Description: "Search for memory documents based on query",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results",
						"default":     10,
					},
					"threshold": map[string]interface{}{
						"type":        "number",
						"description": "Similarity threshold (0.0-1.0)",
						"default":     0.1,
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "list_memories",
			Description: "List all memory documents",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "delete_memory",
			Description: "Delete a memory document by ID",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Document ID to delete",
					},
				},
				Required: []string{"id"},
			},
		},
	}
	
	result := map[string]interface{}{
		"tools": tools,
	}
	
	s.sendResponse(req.ID, result)
}

func (s *MCPServer) handleToolsCall(req MCPRequest) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}
	
	name, ok := params["name"].(string)
	if !ok {
		s.sendError(req.ID, -32602, "Missing tool name")
		return
	}
	
	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}
	
	switch name {
	case "add_memory":
		s.handleAddMemory(req.ID, arguments)
	case "search_memories":
		s.handleSearchMemories(req.ID, arguments)
	case "list_memories":
		s.handleListMemories(req.ID, arguments)
	case "delete_memory":
		s.handleDeleteMemory(req.ID, arguments)
	default:
		s.sendError(req.ID, -32601, "Unknown tool")
	}
}

func (s *MCPServer) handleAddMemory(id interface{}, args map[string]interface{}) {
	content, ok := args["content"].(string)
	if !ok {
		s.sendError(id, -32602, "Missing content")
		return
	}
	
	doc := Document{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: time.Now(),
		Tags:      []string{},
		Properties: make(map[string]string),
	}
	
	if tags, ok := args["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				doc.Tags = append(doc.Tags, tagStr)
			}
		}
	}
	
	if favorite, ok := args["favorite"].(bool); ok {
		doc.Favorite = favorite
	}
	
	if props, ok := args["properties"].(map[string]interface{}); ok {
		for k, v := range props {
			if vStr, ok := v.(string); ok {
				doc.Properties[k] = vStr
			}
		}
	}
	
	if err := s.store.AddDocument(doc); err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Failed to add document: %v", err))
		return
	}
	
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Memory added successfully with ID: %s", doc.ID),
			},
		},
	}
	
	s.sendResponse(id, result)
}

func (s *MCPServer) handleSearchMemories(id interface{}, args map[string]interface{}) {
	query, ok := args["query"].(string)
	if !ok {
		s.sendError(id, -32602, "Missing query")
		return
	}
	
	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	
	threshold := float32(0.1)
	if t, ok := args["threshold"].(float64); ok {
		threshold = float32(t)
	}
	
	docs, err := s.store.SearchDocuments(query, limit, threshold)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Search failed: %v", err))
		return
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
	
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": responseText,
			},
		},
	}
	
	s.sendResponse(id, result)
}

func (s *MCPServer) handleListMemories(id interface{}, args map[string]interface{}) {
	docs, err := s.store.ListDocuments()
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("List failed: %v", err))
		return
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
	
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": responseText,
			},
		},
	}
	
	s.sendResponse(id, result)
}

func (s *MCPServer) handleDeleteMemory(id interface{}, args map[string]interface{}) {
	docID, ok := args["id"].(string)
	if !ok {
		s.sendError(id, -32602, "Missing document ID")
		return
	}
	
	if err := s.store.DeleteDocument(docID); err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Delete failed: %v", err))
		return
	}
	
	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Memory with ID %s deleted successfully", docID),
			},
		},
	}
	
	s.sendResponse(id, result)
}

func (s *MCPServer) sendResponse(id interface{}, result interface{}) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	data, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		return
	}
	
	fmt.Println(string(data))
	log.Debug().Str("response", string(data)).Msg("Sent MCP response")
}

func (s *MCPServer) sendError(id interface{}, code int, message string) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	
	data, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal error response")
		return
	}
	
	fmt.Println(string(data))
	log.Debug().Str("error", string(data)).Msg("Sent MCP error")
}