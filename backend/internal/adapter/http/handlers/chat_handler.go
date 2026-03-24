package handlers

import (
	"database/sql"
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

	// Step 1: Call Groq to generate SQL - let LLM handle OFFTOPIC decision
	logger.InfoWithDataIP(utils.CategoryGroq, "Chat - Calling Groq to generate SQL", clientIP, map[string]interface{}{
		"query": req.Query,
	})

	sqlQuery, err := usecases.GenerateSQL(h.groqAPIKey, req.Query)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryGroq, "Chat - Failed to generate SQL", clientIP, err, map[string]interface{}{
			"query": req.Query,
		})
		// Return user-friendly error message
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "rate limited") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(domain.ChatResponse{
				Answer: "AI service is temporarily unavailable. Please try again in a moment.",
			})
		}
		if strings.Contains(errorMsg, "unavailable") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(domain.ChatResponse{
				Answer: "AI service is temporarily unavailable. Please try again later.",
			})
		}
		if strings.Contains(errorMsg, "API key") {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ChatResponse{
				Answer: "Server configuration error. Please contact support.",
			})
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(domain.ChatResponse{
			Answer: "Failed to process your query. Please try again later.",
		})
	}

	logger.InfoWithDataIP(utils.CategoryGroq, "Chat - SQL generated", clientIP, map[string]interface{}{
		"sql": sqlQuery,
	})

	// Step 2: Check if response is OFFTOPIC
	if strings.TrimSpace(sqlQuery) == "OFFTOPIC" {
		logger.WarnWithDataIP(utils.CategoryGroq, "Chat - Groq returned OFFTOPIC", clientIP, map[string]interface{}{
			"query": req.Query,
		})
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChatResponse{
			Answer: "This system is designed to answer questions related to the provided dataset only.",
		})
	}

	// Step 3: Execute query
	logger.DebugWithDataIP(utils.CategoryQuery, "Chat - Executing SQL query", clientIP, map[string]interface{}{
		"sql": sqlQuery,
	})

	rows, err := h.db.Query(sqlQuery)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryQuery, "Chat - Failed to execute SQL", clientIP, err, map[string]interface{}{
			"sql":   sqlQuery,
			"error": err.Error(),
		})
		// Return user-friendly error based on SQL error
		errorMsg := err.Error()
		userFriendlyMsg := "Query execution failed"

		// Detailed error detection
		if strings.Contains(errorMsg, "unknown function") && (strings.Contains(sqlQuery, "GETDATE") || strings.Contains(sqlQuery, "DATEDIFF")) {
			userFriendlyMsg = "SQL error: This uses SQL Server functions. This is PostgreSQL. Use NOW() instead of GETDATE(), and INTERVAL for date arithmetic."
		} else if strings.Contains(errorMsg, "must appear in the GROUP BY clause") || strings.Contains(errorMsg, "aggregate function") {
			userFriendlyMsg = "SQL error: When using aggregate functions (COUNT, SUM), all non-aggregated columns must be in GROUP BY. Try rephrasing your question."
		} else if strings.Contains(errorMsg, "column") && strings.Contains(errorMsg, "does not exist") {
			userFriendlyMsg = "The query referenced columns that don't exist in the database. The AI may have misunderstood the data structure."
		} else if strings.Contains(errorMsg, "UNION") {
			userFriendlyMsg = "The query has invalid UNION syntax. Please rephrase your question to be more specific."
		} else if strings.Contains(errorMsg, "syntax error") {
			userFriendlyMsg = "The generated query has a syntax error. Please try a simpler or more specific question."
		} else if strings.Contains(errorMsg, "relation") && strings.Contains(errorMsg, "does not exist") {
			userFriendlyMsg = "The query referenced a table that doesn't exist. Please verify the table name."
		} else if strings.Contains(errorMsg, "invalid") {
			userFriendlyMsg = "The query has an issue with its structure. Please try rephrasing your question."
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ChatResponse{
			Answer: userFriendlyMsg,
			SQL:    sqlQuery,
		})
	}
	defer rows.Close()

	// Parse query results into []map[string]interface{}
	results, err := scanRows(rows)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryQuery, "Chat - Failed to scan query results", clientIP, err, map[string]interface{}{
			"sql":   sqlQuery,
			"error": err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ChatResponse{
			Answer: "Failed to process query results",
			SQL:    sqlQuery,
		})
	}

	logger.InfoWithDataIP(utils.CategoryQuery, "Chat - SQL executed successfully", clientIP, map[string]interface{}{
		"rows_returned": len(results),
		"sql":           sqlQuery,
	})

	if len(results) == 0 {
		logger.InfoWithDataIP(utils.CategoryHandler, "Chat - No data found", clientIP, map[string]interface{}{
			"query": req.Query,
		})
		return c.Status(fiber.StatusOK).JSON(domain.ChatResponse{
			Answer: "No data found matching your query.",
			SQL:    sqlQuery,
			Rows:   results,
		})
	}

	// Step 5: Call Groq to format answer
	logger.InfoWithDataIP(utils.CategoryGroq, "Chat - Calling Groq to format answer", clientIP, map[string]interface{}{
		"rows": len(results),
	})

	answer, err := usecases.FormatAnswer(h.groqAPIKey, req.Query, results)
	if err != nil {
		logger.WarnWithDataIP(utils.CategoryGroq, "Chat - Failed to format answer, returning raw results", clientIP, map[string]interface{}{
			"error": err.Error(),
		})
		// Return results even if formatting fails
		return c.Status(fiber.StatusOK).JSON(domain.ChatResponse{
			Answer: "Retrieved data successfully but failed to format answer.",
			SQL:    sqlQuery,
			Rows:   results,
		})
	}

	logger.InfoWithDataIP(utils.CategoryHandler, "Chat - Success", clientIP, map[string]interface{}{
		"query":         req.Query,
		"rows_returned": len(results),
	})

	return c.JSON(domain.ChatResponse{
		Answer: answer,
		SQL:    sqlQuery,
		Rows:   results,
	})
}

// REMOVED: Handler-level keyword check - LLM now handles OFFTOPIC detection via system prompt
// This allows more flexible question handling and better query understanding by the LLM

// scanRows converts sql.Rows into []map[string]interface{}
func scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}

		results = append(results, entry)
	}

	return results, rows.Err()
}
