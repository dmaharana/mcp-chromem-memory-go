package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type WebServer struct {
	store     *MemoryStore
	stats     *ServerStats
	templates *template.Template
}

type ServerStats struct {
	GetDocumentCount    int64 `json:"get_document_count"`
	GetAllDocuments     int64 `json:"get_all_documents"`
	DeleteDocumentCount int64 `json:"delete_document_count"`
	AddDocumentCount    int64 `json:"add_document_count"`
	SearchCount         int64 `json:"search_count"`
}

type WebDocument struct {
	Document
	TagsString string `json:"tags_string"`
}

func NewWebServer(store *MemoryStore) *WebServer {
	ws := &WebServer{
		store: store,
		stats: &ServerStats{},
	}
	
	// Parse HTML templates
	ws.loadTemplates()
	
	return ws
}

func (ws *WebServer) loadTemplates() {
	indexHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Memory Server</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { border-bottom: 2px solid #007bff; padding-bottom: 10px; margin-bottom: 20px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin-bottom: 30px; }
        .stat-card { background: #007bff; color: white; padding: 15px; border-radius: 5px; text-align: center; }
        .stat-number { font-size: 24px; font-weight: bold; }
        .stat-label { font-size: 14px; margin-top: 5px; }
        .section { margin-bottom: 30px; }
        .section h2 { color: #333; border-bottom: 1px solid #ddd; padding-bottom: 5px; }
        .form-group { margin-bottom: 15px; }
        .form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
        .form-group input, .form-group textarea, .form-group select { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        .form-group textarea { height: 100px; resize: vertical; }
        .btn { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; margin-right: 10px; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #218838; }
        .document { border: 1px solid #ddd; padding: 15px; margin-bottom: 15px; border-radius: 5px; background: #f9f9f9; }
        .document.favorite { border-left: 4px solid #ffc107; }
        .document-header { display: flex; justify-content: between; align-items: center; margin-bottom: 10px; }
        .document-id { font-family: monospace; color: #666; font-size: 12px; }
        .document-content { margin: 10px 0; }
        .document-meta { font-size: 12px; color: #666; margin-top: 10px; }
        .tags { margin: 5px 0; }
        .tag { background: #e9ecef; padding: 2px 6px; border-radius: 3px; font-size: 11px; margin-right: 5px; }
        .favorite-star { color: #ffc107; font-size: 18px; }
        .search-box { width: 100%; padding: 10px; margin-bottom: 20px; border: 1px solid #ddd; border-radius: 4px; }
        .hidden { display: none; }
        .edit-form { background: #fff; border: 2px solid #007bff; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Memory Server Dashboard</h1>
            <p>Local Memory Layer for Developers</p>
        </div>

        <div class="stats">
            <div class="stat-card">
                <div class="stat-number" id="total-docs">0</div>
                <div class="stat-label">Total Documents</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="add-count">0</div>
                <div class="stat-label">Documents Added</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="search-count">0</div>
                <div class="stat-label">Searches Performed</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="delete-count">0</div>
                <div class="stat-label">Documents Deleted</div>
            </div>
        </div>

        <div class="section">
            <h2>Add New Memory</h2>
            <form id="add-form">
                <div class="form-group">
                    <label for="content">Content:</label>
                    <textarea id="content" name="content" required placeholder="Enter the memory content..."></textarea>
                </div>
                <div class="form-group">
                    <label for="tags">Tags (comma-separated):</label>
                    <input type="text" id="tags" name="tags" placeholder="golang, debugging, tips">
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="favorite" name="favorite"> Mark as Favorite
                    </label>
                </div>
                <div class="form-group">
                    <label for="properties">Properties (JSON format):</label>
                    <textarea id="properties" name="properties" placeholder='{"category": "tip", "language": "go"}'></textarea>
                </div>
                <button type="submit" class="btn">Add Memory</button>
            </form>
        </div>

        <div class="section">
            <h2>Search & Browse Memories</h2>
            <input type="text" id="search-input" class="search-box" placeholder="Search memories...">
            <button onclick="searchDocuments()" class="btn">Search</button>
            <button onclick="loadAllDocuments()" class="btn">Show All</button>
            
            <div id="documents-container">
                <!-- Documents will be loaded here -->
            </div>
        </div>

        <!-- Edit Modal -->
        <div id="edit-modal" class="hidden">
            <div class="edit-form">
                <h3>Edit Memory</h3>
                <form id="edit-form">
                    <input type="hidden" id="edit-id">
                    <div class="form-group">
                        <label for="edit-content">Content:</label>
                        <textarea id="edit-content" name="content" required></textarea>
                    </div>
                    <div class="form-group">
                        <label for="edit-tags">Tags (comma-separated):</label>
                        <input type="text" id="edit-tags" name="tags">
                    </div>
                    <div class="form-group">
                        <label>
                            <input type="checkbox" id="edit-favorite" name="favorite"> Mark as Favorite
                        </label>
                    </div>
                    <div class="form-group">
                        <label for="edit-properties">Properties (JSON format):</label>
                        <textarea id="edit-properties" name="properties"></textarea>
                    </div>
                    <button type="submit" class="btn btn-success">Save Changes</button>
                    <button type="button" onclick="cancelEdit()" class="btn">Cancel</button>
                </form>
            </div>
        </div>
    </div>

    <script>
        // Load initial data
        loadStats();
        loadAllDocuments();

        // Add form submission
        document.getElementById('add-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = {
                content: formData.get('content'),
                tags: formData.get('tags').split(',').map(t => t.trim()).filter(t => t),
                favorite: formData.get('favorite') === 'on',
                properties: {}
            };
            
            const propsText = formData.get('properties');
            if (propsText) {
                try {
                    data.properties = JSON.parse(propsText);
                } catch (e) {
                    alert('Invalid JSON in properties field');
                    return;
                }
            }

            try {
                const response = await fetch('/api/documents', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });
                
                if (response.ok) {
                    e.target.reset();
                    loadStats();
                    loadAllDocuments();
                    alert('Memory added successfully!');
                } else {
                    alert('Failed to add memory');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        });

        // Edit form submission
        document.getElementById('edit-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const id = document.getElementById('edit-id').value;
            const data = {
                content: formData.get('content'),
                tags: formData.get('tags').split(',').map(t => t.trim()).filter(t => t),
                favorite: formData.get('favorite') === 'on',
                properties: {}
            };
            
            const propsText = formData.get('properties');
            if (propsText) {
                try {
                    data.properties = JSON.parse(propsText);
                } catch (e) {
                    alert('Invalid JSON in properties field');
                    return;
                }
            }

            try {
                const response = await fetch('/api/documents/' + id, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });
                
                if (response.ok) {
                    cancelEdit();
                    loadStats();
                    loadAllDocuments();
                    alert('Memory updated successfully!');
                } else {
                    alert('Failed to update memory');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        });

        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
                const stats = await response.json();
                document.getElementById('total-docs').textContent = stats.total_documents;
                document.getElementById('add-count').textContent = stats.add_document_count;
                document.getElementById('search-count').textContent = stats.search_count;
                document.getElementById('delete-count').textContent = stats.delete_document_count;
            } catch (error) {
                console.error('Failed to load stats:', error);
            }
        }

        async function loadAllDocuments() {
            try {
                const response = await fetch('/api/documents');
                const documents = await response.json();
                displayDocuments(documents);
            } catch (error) {
                console.error('Failed to load documents:', error);
            }
        }

        async function searchDocuments() {
            const query = document.getElementById('search-input').value;
            if (!query.trim()) {
                loadAllDocuments();
                return;
            }

            try {
                const response = await fetch('/api/search?q=' + encodeURIComponent(query));
                const documents = await response.json();
                displayDocuments(documents);
            } catch (error) {
                console.error('Failed to search documents:', error);
            }
        }

        function displayDocuments(documents) {
            const container = document.getElementById('documents-container');
            if (documents.length === 0) {
                container.innerHTML = '<p>No documents found.</p>';
                return;
            }

            container.innerHTML = documents.map(doc => {
                const tags = doc.tags ? doc.tags.map(tag => '<span class="tag">' + tag + '</span>').join('') : '';
                const favorite = doc.favorite ? '<span class="favorite-star">⭐</span>' : '';
                const createdAt = new Date(doc.created_at).toLocaleString();
                
                return '<div class="document' + (doc.favorite ? ' favorite' : '') + '">' +
                    '<div class="document-header">' +
                        '<div class="document-id">ID: ' + doc.id + '</div>' +
                        '<div>' + favorite + '</div>' +
                    '</div>' +
                    '<div class="document-content">' + doc.content + '</div>' +
                    '<div class="tags">' + tags + '</div>' +
                    '<div class="document-meta">Created: ' + createdAt + '</div>' +
                    '<div style="margin-top: 10px;">' +
                        '<button onclick="editDocument(\'' + doc.id + '\')" class="btn">Edit</button>' +
                        '<button onclick="toggleFavorite(\'' + doc.id + '\', ' + !doc.favorite + ')" class="btn">' + 
                            (doc.favorite ? 'Remove Favorite' : 'Add Favorite') + '</button>' +
                        '<button onclick="deleteDocument(\'' + doc.id + '\')" class="btn btn-danger">Delete</button>' +
                    '</div>' +
                '</div>';
            }).join('');
        }

        async function editDocument(id) {
            try {
                const response = await fetch('/api/documents/' + id);
                const doc = await response.json();
                
                document.getElementById('edit-id').value = doc.id;
                document.getElementById('edit-content').value = doc.content;
                document.getElementById('edit-tags').value = doc.tags ? doc.tags.join(', ') : '';
                document.getElementById('edit-favorite').checked = doc.favorite;
                document.getElementById('edit-properties').value = JSON.stringify(doc.properties || {}, null, 2);
                
                document.getElementById('edit-modal').classList.remove('hidden');
            } catch (error) {
                alert('Failed to load document for editing');
            }
        }

        function cancelEdit() {
            document.getElementById('edit-modal').classList.add('hidden');
        }

        async function toggleFavorite(id, favorite) {
            try {
                const response = await fetch('/api/documents/' + id + '/favorite', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ favorite: favorite })
                });
                
                if (response.ok) {
                    loadAllDocuments();
                } else {
                    alert('Failed to update favorite status');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        }

        async function deleteDocument(id) {
            if (!confirm('Are you sure you want to delete this memory?')) {
                return;
            }

            try {
                const response = await fetch('/api/documents/' + id, {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    loadStats();
                    loadAllDocuments();
                    alert('Memory deleted successfully!');
                } else {
                    alert('Failed to delete memory');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        }

        // Search on Enter key
        document.getElementById('search-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                searchDocuments();
            }
        });
    </script>
</body>
</html>`

	var err error
	ws.templates, err = template.New("index").Parse(indexHTML)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse HTML template")
	}
}

func (ws *WebServer) Start(port int) error {
	http.HandleFunc("/", ws.handleIndex)
	http.HandleFunc("/api/stats", ws.handleStats)
	http.HandleFunc("/api/documents", ws.handleDocuments)
	http.HandleFunc("/api/documents/", ws.handleDocumentByID)
	http.HandleFunc("/api/search", ws.handleSearch)

	addr := fmt.Sprintf(":%d", port)
	log.Info().Str("addr", addr).Msg("Starting web server")
	return http.ListenAndServe(addr, nil)
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	if err := ws.templates.Execute(w, nil); err != nil {
		log.Error().Err(err).Msg("Failed to execute template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	docs, err := ws.store.ListDocuments()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get document count")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	stats := map[string]interface{}{
		"total_documents":      len(docs),
		"add_document_count":   ws.stats.AddDocumentCount,
		"search_count":         ws.stats.SearchCount,
		"delete_document_count": ws.stats.DeleteDocumentCount,
		"get_document_count":   ws.stats.GetDocumentCount,
		"get_all_documents":    ws.stats.GetAllDocuments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (ws *WebServer) handleDocuments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ws.stats.GetAllDocuments++
		docs, err := ws.store.ListDocuments()
		if err != nil {
			log.Error().Err(err).Msg("Failed to list documents")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)

	case http.MethodPost:
		var doc Document
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		doc.ID = uuid.New().String()
		doc.CreatedAt = time.Now()
		
		if err := ws.store.AddDocument(doc); err != nil {
			log.Error().Err(err).Msg("Failed to add document")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		ws.stats.AddDocumentCount++
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": doc.ID, "status": "created"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ws *WebServer) handleDocumentByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/documents/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Document ID required", http.StatusBadRequest)
		return
	}
	
	id := parts[0]
	
	// Handle favorite toggle endpoint
	if len(parts) > 1 && parts[1] == "favorite" {
		ws.handleToggleFavorite(w, r, id)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ws.stats.GetDocumentCount++
		docs, err := ws.store.ListDocuments()
		if err != nil {
			log.Error().Err(err).Msg("Failed to list documents")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		for _, doc := range docs {
			if doc.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(doc)
				return
			}
		}
		
		http.Error(w, "Document not found", http.StatusNotFound)

	case "PUT":
		var updateDoc Document
		if err := json.NewDecoder(r.Body).Decode(&updateDoc); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// First delete the old document
		if err := ws.store.DeleteDocument(id); err != nil {
			log.Error().Err(err).Str("id", id).Msg("Failed to delete document for update")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		// Then add the updated document with the same ID
		updateDoc.ID = id
		updateDoc.CreatedAt = time.Now() // Update timestamp
		
		if err := ws.store.AddDocument(updateDoc); err != nil {
			log.Error().Err(err).Msg("Failed to update document")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "updated"})

	case http.MethodDelete:
		if err := ws.store.DeleteDocument(id); err != nil {
			log.Error().Err(err).Str("id", id).Msg("Failed to delete document")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		ws.stats.DeleteDocumentCount++
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ws *WebServer) handleToggleFavorite(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Favorite bool `json:"favorite"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get the current document
	docs, err := ws.store.ListDocuments()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list documents")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var currentDoc *Document
	for _, doc := range docs {
		if doc.ID == id {
			currentDoc = &doc
			break
		}
	}

	if currentDoc == nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Update favorite status
	currentDoc.Favorite = req.Favorite

	// Delete and re-add with updated favorite status
	if err := ws.store.DeleteDocument(id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete document for favorite update")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ws.store.AddDocument(*currentDoc); err != nil {
		log.Error().Err(err).Msg("Failed to update document favorite status")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"favorite": req.Favorite,
		"status":   "updated",
	})
}

func (ws *WebServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	thresholdStr := r.URL.Query().Get("threshold")
	threshold := float32(0.1)
	if thresholdStr != "" {
		if t, err := strconv.ParseFloat(thresholdStr, 32); err == nil && t >= 0 && t <= 1 {
			threshold = float32(t)
		}
	}

	ws.stats.SearchCount++
	docs, err := ws.store.SearchDocuments(query, limit, threshold)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search documents")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}