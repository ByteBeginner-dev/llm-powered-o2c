package handlers

import (
	"o2c-graph/internal/adapter/db"
	"o2c-graph/internal/core/usecases"
	"o2c-graph/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetNode returns a single node with its neighbors and connecting edges
func (h *Handler) GetNode(c *fiber.Ctx) error {
	clientIP := utils.ExtractClientIP(c.IP(), c.Get("X-Forwarded-For"), c.Get("X-Real-IP"))
	logger := utils.GetLogger()

	nodeType := c.Params("type")
	nodeID := c.Params("id")

	logger.InfoWithDataIP(utils.CategoryHandler, "GetNode - Request received", clientIP, map[string]interface{}{
		"type":     nodeType,
		"id":       nodeID,
		"endpoint": "/api/node/:type/:id",
	})

	if nodeType == "" || nodeID == "" {
		logger.WarnWithDataIP(utils.CategoryValidation, "GetNode - Missing parameters", clientIP, map[string]interface{}{
			"type": nodeType,
			"id":   nodeID,
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing type or id parameter",
		})
	}

	// Create repository and usecase instances
	nodeRepo := db.NewNodeRepository(h.db)
	nodeUsecase := usecases.NewNodeUsecase(nodeRepo)

	// Get node detail
	nodeDetail, err := nodeUsecase.GetNodeDetail(nodeType, nodeID)
	if err != nil {
		logger.ErrorWithDataIP(utils.CategoryHandler, "GetNode - Failed to fetch node detail", clientIP, err, map[string]interface{}{
			"type":  nodeType,
			"id":    nodeID,
			"error": err.Error(),
		})
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	logger.InfoWithDataIP(utils.CategoryHandler, "GetNode - Success", clientIP, map[string]interface{}{
		"type":      nodeType,
		"id":        nodeID,
		"neighbors": len(nodeDetail.Neighbors),
		"edges":     len(nodeDetail.Edges),
	})

	return c.JSON(nodeDetail)
}
