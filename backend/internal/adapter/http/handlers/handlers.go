package handlers

import (
	"database/sql"
)

// Handler contains dependencies for all HTTP handlers
type Handler struct {
	db         *sql.DB
	groqAPIKey string
}

// New creates a new Handler with database and API key dependencies
func New(db *sql.DB, groqAPIKey string) *Handler {
	return &Handler{
		db:         db,
		groqAPIKey: groqAPIKey,
	}
}
