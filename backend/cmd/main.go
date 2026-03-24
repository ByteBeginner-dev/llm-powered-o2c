package main

import (
	"log"

	"o2c-graph/internal/adapter/db"
	"o2c-graph/internal/adapter/http/handlers"
	"o2c-graph/internal/adapter/http/routes"
	"o2c-graph/internal/config"
	ingest "o2c-graph/internal/infra/ingest"
	migrate "o2c-graph/internal/infra/migrate"
	"o2c-graph/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// 2. Load config
	cfg := config.Load()

	// 3. Initialize logger
	if err := utils.InitLogger("./logs"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.GetLogger().Close()
	logger := utils.GetLogger()

	logger.Info(utils.CategoryServer, "Starting O2C Graph API server")

	// 4. Connect to PostgreSQL
	database, err := db.Init(cfg.DatabaseURL)
	if err != nil {
		logger.Error(utils.CategoryDatabase, "Failed to connect to database", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	logger.Info(utils.CategoryDatabase, "✓ Connected to PostgreSQL")

	// 5. Run migrations (creates all 19 tables)
	if err := migrate.Run(database, "./schema.sql"); err != nil {
		logger.Error(utils.CategoryMigration, "Migration failed", err)
		log.Fatalf("Migration failed: %v", err)
	}
	logger.Info(utils.CategoryMigration, "✓ Migrations applied - All 19 tables created")

	// 6. Run ingestion (fills tables from JSONL files, skips if already loaded)
	if err := ingest.Run(database, cfg.DataDir); err != nil {
		logger.Error(utils.CategoryIngestion, "Data ingestion failed", err)
		log.Fatalf("Ingestion failed: %v", err)
	}
	logger.Info(utils.CategoryIngestion, "✓ Data ingestion complete - All JSONL files processed")

	// 7. Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "O2C-Graph-API",
	})

	// 8. Request logging middleware
	app.Use(func(c *fiber.Ctx) error {
		clientIP := utils.ExtractClientIP(c.IP(), c.Get("X-Forwarded-For"), c.Get("X-Real-IP"))
		method := c.Method()
		path := c.Path()
		query := c.Query("")

		logger.InfoWithDataIP(utils.CategoryHandler,
			"Incoming request",
			clientIP,
			map[string]interface{}{
				"method": method,
				"path":   path,
				"query":  query,
			},
		)

		// Process request
		err := c.Next()

		// Log response
		statusCode := c.Response().StatusCode()
		logger.InfoWithDataIP(utils.CategoryHandler,
			"Response sent",
			clientIP,
			map[string]interface{}{
				"method":      method,
				"path":        path,
				"status_code": statusCode,
			},
		)

		return err
	})

	// 9. CORS — allow React frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,Accept,Origin,X-Forwarded-For,X-Real-IP",
		AllowCredentials: false,
		MaxAge:           3600,
	}))

	// 10. Register all routes
	h := handlers.New(database, cfg.GroqApiKey)
	routes.Register(app, h)
	logger.Info(utils.CategoryRoute, "All routes registered successfully")

	// 11. Serve React frontend build (static files)
	app.Static("/", "./frontend/dist")

	// 12. Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		clientIP := utils.ExtractClientIP(c.IP(), c.Get("X-Forwarded-For"), c.Get("X-Real-IP"))
		logger.InfoWithIP(utils.CategoryServer, "Health check", clientIP)
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// 13. Start server
	port := cfg.Port
	if port == "" {
		port = "8189"
	}

	logger.InfoWithData(utils.CategoryServer, "Server starting", map[string]interface{}{
		"port": port,
		"app":  "O2C-Graph-API",
	})

	if err := app.Listen(":" + port); err != nil {
		logger.Error(utils.CategoryServer, "Server error", err)
		log.Fatalf("Server error: %v", err)
	}
}
