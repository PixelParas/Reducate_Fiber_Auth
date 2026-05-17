package main

import (
	"fiber_auth/config"
	"fiber_auth/database"
	"fiber_auth/handlers"
	"fiber_auth/middleware"
	"os/signal"

	"log"
	"os"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to laod configuration: %w", err)
	}

	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Databse initilization failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	handlers.InitHandlers(cfg)

	middleware.InitMiddleware(cfg)

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api")
	api.Post("/signup", handlers.Signup)
	api.Post("/login", handlers.Login)

	//Protexted routes
	protected := api.Group("/", middleware.RequireAuth)
	protected.Get("/profile", handlers.GetProfile)

	//admin Only routes
	admin := api.Group("/", middleware.RequireAdmin)
	admin.Get("/users", handlers.GetAllUsers)

	app.Get("api/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server is running on port: %v", cfg.Port)

	<-quit

	log.Println("Server Shutting down....")

	if err := app.Shutdown(); err != nil {
		log.Fatalf(("Server forced to shutdown: %v"), err)
	}

	log.Println("server exited gracefully")
}
