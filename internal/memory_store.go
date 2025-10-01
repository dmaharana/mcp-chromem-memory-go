package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/philippgille/chromem-go"
	"github.com/rs/zerolog/log"
)

type Document struct {
	ID         string            `json:"id"`
	Content    string            `json:"content"`
	Tags       []string          `json:"tags"`
	Properties map[string]string `json:"properties"`
	Favorite   bool              `json:"favorite"`
	CreatedAt  time.Time         `json:"created_at"`
}

type MemoryStore struct {
	db *chromem.DB
}

func NewMemoryStore(path string) (*MemoryStore, error) {
	log.Info().Str("path", path).Msg("Initializing memory store")
	
	db := chromem.NewDB()
	
	// Create collection with custom embedding function
	embedder := NewStatisticalEmbedder()
	collection, err := db.CreateCollection("memories", nil, embedder)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}
	
	log.Info().Str("collection", collection.Name).Msg("Memory store initialized")
	
	return &MemoryStore{
		db: db,
	}, nil
}

func (ms *MemoryStore) AddDocument(doc Document) error {
	log.Info().Str("id", doc.ID).Msg("Adding document to memory store")
	
	collection := ms.db.GetCollection("memories", nil)
	if collection == nil {
		return fmt.Errorf("collection not found")
	}
	
	// Prepare metadata
	metadata := make(map[string]string)
	metadata["tags"] = strings.Join(doc.Tags, ",")
	if doc.Favorite {
		metadata["favorite"] = "true"
	} else {
		metadata["favorite"] = "false"
	}
	metadata["created_at"] = doc.CreatedAt.Format(time.RFC3339)
	
	// Add properties to metadata
	for k, v := range doc.Properties {
		metadata["prop_"+k] = v
	}
	
	err := collection.AddDocument(context.Background(), chromem.Document{
		ID:       doc.ID,
		Content:  doc.Content,
		Metadata: metadata,
	})
	
	if err != nil {
		log.Error().Err(err).Str("id", doc.ID).Msg("Failed to add document")
		return fmt.Errorf("failed to add document: %w", err)
	}
	
	log.Info().Str("id", doc.ID).Msg("Document added successfully")
	return nil
}

func (ms *MemoryStore) SearchDocuments(query string, limit int, threshold float32) ([]Document, error) {
	log.Info().Str("query", query).Int("limit", limit).Float32("threshold", threshold).Msg("Searching documents")
	
	collection := ms.db.GetCollection("memories", nil)
	if collection == nil {
		return nil, fmt.Errorf("collection not found")
	}
	
	results, err := collection.Query(context.Background(), query, limit, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search documents")
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}
	
	var documents []Document
	for _, result := range results {
		// Filter by similarity threshold
		if result.Similarity < threshold {
			continue
		}
		
		doc := Document{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: time.Now(), // Default value
		}
		
		// Parse metadata
		if tagsStr, ok := result.Metadata["tags"]; ok && tagsStr != "" {
			doc.Tags = strings.Split(tagsStr, ",")
		}
		if favoriteStr, ok := result.Metadata["favorite"]; ok {
			doc.Favorite = favoriteStr == "true"
		}
		if createdAt, ok := result.Metadata["created_at"]; ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				doc.CreatedAt = t
			}
		}
		
		// Parse properties
		doc.Properties = make(map[string]string)
		for k, v := range result.Metadata {
			if strings.HasPrefix(k, "prop_") {
				propKey := strings.TrimPrefix(k, "prop_")
				doc.Properties[propKey] = v
			}
		}
		
		// Boost favorite documents
		if doc.Favorite {
			result.Similarity *= 1.2 // Boost favorite documents
		}
		
		documents = append(documents, doc)
	}
	
	log.Info().Int("count", len(documents)).Msg("Search completed")
	return documents, nil
}

func (ms *MemoryStore) DeleteDocument(id string) error {
	log.Info().Str("id", id).Msg("Deleting document")
	
	collection := ms.db.GetCollection("memories", nil)
	if collection == nil {
		return fmt.Errorf("collection not found")
	}
	
	err := collection.Delete(context.Background(), nil, map[string]string{"id": id})
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete document")
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	log.Info().Str("id", id).Msg("Document deleted successfully")
	return nil
}

func (ms *MemoryStore) ListDocuments() ([]Document, error) {
	log.Info().Msg("Listing all documents")
	
	collection := ms.db.GetCollection("memories", nil)
	if collection == nil {
		return nil, fmt.Errorf("collection not found")
	}
	
	// Get all documents by querying with empty string and high limit
	results, err := collection.Query(context.Background(), "", 1000, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list documents")
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	
	var documents []Document
	for _, result := range results {
		doc := Document{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: time.Now(),
		}
		
		// Parse metadata (same as in SearchDocuments)
		if tagsStr, ok := result.Metadata["tags"]; ok && tagsStr != "" {
			doc.Tags = strings.Split(tagsStr, ",")
		}
		if favoriteStr, ok := result.Metadata["favorite"]; ok {
			doc.Favorite = favoriteStr == "true"
		}
		if createdAt, ok := result.Metadata["created_at"]; ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				doc.CreatedAt = t
			}
		}
		
		doc.Properties = make(map[string]string)
		for k, v := range result.Metadata {
			if strings.HasPrefix(k, "prop_") {
				propKey := strings.TrimPrefix(k, "prop_")
				doc.Properties[propKey] = v
			}
		}
		
		documents = append(documents, doc)
	}
	
	log.Info().Int("count", len(documents)).Msg("Listed all documents")
	return documents, nil
}

func (ms *MemoryStore) Close() error {
	log.Info().Msg("Closing memory store")
	return nil
}