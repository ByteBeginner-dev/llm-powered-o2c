package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS sets up CORS middleware
func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,Accept,Origin,X-Forwarded-For,X-Real-IP",
		AllowCredentials: false,
		MaxAge:           3600,
	})
}
