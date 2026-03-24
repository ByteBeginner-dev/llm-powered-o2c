package routes

import (
	"o2c-graph/internal/adapter/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// Register wires all API routes to their handlers
func Register(app *fiber.App, h *handlers.Handler) {
	api := app.Group("/api")

	// Graph endpoints
	api.Get("/graph", h.GetGraph)         // returns all nodes + edges
	api.Get("/node/:type/:id", h.GetNode) // returns single node + neighbors

	// Chat endpoint
	api.Post("/chat", h.Chat) // natural language → SQL → answer
}
