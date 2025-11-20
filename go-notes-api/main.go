package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Note struct - ek note ka data
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// In-memory notes (with file persistence)
var (
	notes  = []Note{}
	nextID = 1
	mu     sync.Mutex
)

const dataFile = "notes.json"

func main() {
	// Load notes from file at startup (permanent storage)
	if err := loadNotesFromFile(); err != nil {
		log.Println("Could not load notes from file:", err)
	}

	r := gin.Default()

	// CORS + custom logging middleware
	r.Use(corsMiddleware())

	// Routes
	r.GET("/notes", getNotes)          // saare notes
	r.POST("/notes", addNote)          // naya note add
	r.GET("/notes/:id", getNoteByID)   // id se ek note
	r.PUT("/notes/:id", updateNote)    // note update
	r.DELETE("/notes/:id", deleteNote) // note delete

	// Simple docs endpoint (Swagger-style info)
	r.GET("/docs", docsHandler)

	// Simple frontend serve
	r.StaticFile("/", "./frontend.html")

	log.Println("Server listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// ================== Handlers ==================

// GET /notes
func getNotes(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, notes)
}

// POST /notes
func addNote(c *gin.Context) {
	var newNote Note

	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newNote.ID = nextID
	nextID++
	notes = append(notes, newNote)

	if err := saveNotesToFile(); err != nil {
		log.Println("Failed to save notes:", err)
	}

	c.JSON(http.StatusCreated, newNote)
}

// GET /notes/:id
func getNoteByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, note := range notes {
		if note.ID == id {
			c.JSON(http.StatusOK, note)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// PUT /notes/:id  (Update note)
func updateNote(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	var updated Note
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, note := range notes {
		if note.ID == id {
			notes[i].Title = updated.Title
			notes[i].Content = updated.Content

			if err := saveNotesToFile(); err != nil {
				log.Println("Failed to save notes:", err)
			}

			c.JSON(http.StatusOK, notes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// DELETE /notes/:id
func deleteNote(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)

			if err := saveNotesToFile(); err != nil {
				log.Println("Failed to save notes:", err)
			}

			c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// ================== Helper functions ==================

func parseIDParam(c *gin.Context) (int, error) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, err
	}
	return id, nil
}

// ================== Persistence (file as DB) ==================

func loadNotesFromFile() error {
	file, err := os.Open(dataFile)
	if os.IsNotExist(err) {
		// No existing file - first run
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var loaded []Note
	if err := json.NewDecoder(file).Decode(&loaded); err != nil {
		return err
	}

	notes = loaded
	// Set nextID correctly
	maxID := 0
	for _, n := range notes {
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	nextID = maxID + 1

	return nil
}

func saveNotesToFile() error {
	file, err := os.Create(dataFile) // truncate + create
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(notes)
}

// ================== Docs Handler ==================

func docsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":       "Go Notes API",
		"description": "Simple REST API for managing notes (Go + Gin)",
		"endpoints": []gin.H{
			{"method": "GET", "path": "/notes", "description": "Get all notes"},
			{"method": "POST", "path": "/notes", "description": "Create a new note"},
			{"method": "GET", "path": "/notes/:id", "description": "Get note by ID"},
			{"method": "PUT", "path": "/notes/:id", "description": "Update note by ID"},
			{"method": "DELETE", "path": "/notes/:id", "description": "Delete note by ID"},
		},
	})
}

// ================== Middleware ==================

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s\n", c.Request.Method, c.Request.URL.Path)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
