package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"o2c-graph/internal/core/domain"
	"o2c-graph/internal/core/usecases"
	"o2c-graph/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Chat handles natural language queries via Groq API
func (h *Handler) Chat(c *fiber.Ctx) error {
	clientIP := utils.ExtractClientIP(c.IP(), c.Get("X-Forwarded-For"), c.Get("X-Real-IP"))
	logger := utils.GetLogger()

	var req domain.ChatRequest

	if err := c.BodyParser(&req); err != nil {
		logger.ErrorWithDataIP(utils.CategoryHandler, "Chat - Failed to parse request", clientIP, err, map[string]interface{}{
			"endpoint": "/api/chat",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		logger.WarnWithDataIP(utils.CategoryValidation, "Chat - Empty query", clientIP, map[string]interface{}{
			"endpoint": "/api/chat",
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query cannot be empty",
		})
	}

	logger.InfoWithDataIP(utils.CategoryHandler, "Chat - Request received", clientIP, map[string]interface{}{
		"query": req.Query,
	})

	// Step 1: Call Groq to generate SQL
	generatedSQL, err := usecases.GenerateSQL(h.groqAPIKey, req.Query)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryGroq, "Chat - Failed to generate SQL", clientIP, err, map[string]interface{}{
			"query": req.Query,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ChatResponse{
			Answer: "Failed to process your question. Please try again.",
		})
	}

	logger.InfoWithDataIP(utils.CategoryGroq, "Chat - SQL generated", clientIP, map[string]interface{}{
		"query":         req.Query,
		"generated_sql": generatedSQL,
	})

	// Step 2: Check if query is off-topic
	if strings.TrimSpace(generatedSQL) == "OFFTOPIC" {
		logger.WarnWithDataIP(utils.CategoryValidation, "Chat - Off-topic query rejected", clientIP, map[string]interface{}{
			"query": req.Query,
		})
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChatResponse{
			Answer: "This system is designed to answer questions related to the provided dataset only.",
		})
	}

	// Step 3: Sanitize the SQL before execution
	generatedSQL = sanitizeSQL(generatedSQL)

	logger.InfoWithDataIP(utils.CategoryHandler, "Chat - SQL after sanitization", clientIP, map[string]interface{}{
		"sanitized_sql": generatedSQL,
	})

	// Step 4: Validate SQL is not empty after sanitization
	if generatedSQL == "" {
		logger.WarnWithDataIP(utils.CategoryValidation, "Chat - SQL empty after sanitization", clientIP, map[string]interface{}{
			"query": req.Query,
		})
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChatResponse{
			Answer: "Could not generate a valid query for that question. Please try rephrasing.",
		})
	}

	// Step 5: Execute SQL against PostgreSQL
	rows, err := executeQuery(h.db, generatedSQL)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryDatabase, "Chat - Query execution failed", clientIP, err, map[string]interface{}{
			"query":         req.Query,
			"generated_sql": generatedSQL,
			"error":         err.Error(),
		})
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChatResponse{
			Answer: fmt.Sprintf("I generated a query but it failed to execute. Try rephrasing your question. Error: %s", err.Error()),
			SQL:    generatedSQL,
		})
	}

	logger.InfoWithDataIP(utils.CategoryDatabase, "Chat - Query executed successfully", clientIP, map[string]interface{}{
		"query":     req.Query,
		"row_count": len(rows),
	})

	// Step 6: Format the answer using Groq
	answer, err := usecases.FormatAnswer(h.groqAPIKey, req.Query, rows)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryGroq, "Chat - Failed to format answer", clientIP, err, map[string]interface{}{
			"query": req.Query,
		})
		// Return raw data if formatting fails — don't return 500
		return c.JSON(domain.ChatResponse{
			Answer: "Query executed successfully but formatting failed.",
			SQL:    generatedSQL,
			Rows:   rows,
		})
	}

	logger.InfoWithDataIP(utils.CategoryHandler, "Chat - Response ready", clientIP, map[string]interface{}{
		"query":     req.Query,
		"answer":    answer,
		"row_count": len(rows),
	})

	return c.JSON(domain.ChatResponse{
		Answer: answer,
		SQL:    generatedSQL,
		Rows:   rows,
	})
}

// sanitizeSQL cleans LLM output before executing against PostgreSQL
func sanitizeSQL(raw string) string {
	s := strings.TrimSpace(raw)

	// Remove markdown code fences that LLM sometimes adds
	for _, prefix := range []string{"```sql", "```postgresql", "```pgsql", "```"} {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			break
		}
	}

	// Remove trailing code fence
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}

	// Remove trailing semicolons
	s = strings.TrimRight(strings.TrimSpace(s), ";")

	// Remove any leading newlines left after stripping fences
	s = strings.TrimSpace(s)

	return s
}

// executeQuery runs a SQL query and returns results as a slice of maps
func executeQuery(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Build the row map
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for readability
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for clean JSON
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}
