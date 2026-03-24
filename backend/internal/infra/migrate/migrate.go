package migrate

import (
	"database/sql"
	"fmt"
	"os"
)

// Run reads schema.sql and executes it against the database.
// It is idempotent — all CREATE TABLE statements use IF NOT EXISTS.
func Run(db *sql.DB, schemaPath string) error {
	// Read the schema.sql file
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file at %s: %w", schemaPath, err)
	}

	// Execute the entire schema as one statement block
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
