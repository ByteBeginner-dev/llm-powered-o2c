package handlers

import (
	"o2c-graph/internal/adapter/db"
	"o2c-graph/internal/core/usecases"
	"o2c-graph/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetGraph returns all nodes and edges in the O2C graph
// Handler layer: receives HTTP request, delegates to usecase, returns JSON
func (h *Handler) GetGraph(c *fiber.Ctx) error {
	clientIP := utils.ExtractClientIP(c.IP(), c.Get("X-Forwarded-For"), c.Get("X-Real-IP"))
	logger := utils.GetLogger()

	logger.InfoWithIP(utils.CategoryHandler, "GetGraph - Request received", clientIP)

	// Create repository instance (data access layer)
	graphRepo := db.NewGraphRepository(h.db)

	// Create usecase instance (business logic layer)
	graphUsecase := usecases.NewGraphUsecase(graphRepo)

	// Call business logic
	graph, err := graphUsecase.GetGraph()
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryHandler, "GetGraph - Failed to fetch graph", clientIP, err, map[string]interface{}{
			"endpoint": "/api/graph",
			"error":    err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch graph",
		})
	}

	logger.InfoWithDataIP(utils.CategoryHandler, "GetGraph - Success", clientIP, map[string]interface{}{
		"nodes_count": len(graph.Nodes),
		"edges_count": len(graph.Edges),
	})

	// Return JSON response
	return c.JSON(graph)
}
