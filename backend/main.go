package main

import (
	"crawl-manager-backend/app/exceptions"
	"crawl-manager-backend/bootstrap"
	"crawl-manager-backend/routes"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	startServer()
}

func startServer() {
	app := bootstrap.App()
	defer app.CloseDBConnection()
	app.ConnectDB()

	// Add middleware to log HTTP requests
	app.App.Use(requestLogger)

	// Register routes
	routes.RegisterRoutes(app.App)

	// Launch the application in a goroutine
	go startApplication(app)

	// Graceful shutdown
	app.GracefulShutdown(func() {
		shutdownApplication(app)
	})
}

func startApplication(app *bootstrap.Application) {
	port := ":" + app.Config.App.Port
	if err := app.App.Listen(port); err != nil {
		exceptions.PanicIfNeeded(err.Error())
	}
}

func shutdownApplication(app *bootstrap.Application) {
	if err := app.App.Shutdown(); err != nil {
		fmt.Println("Error during shutdown:", err.Error())
	}

	app.CloseDBConnection()
}

// Middleware to log HTTP requests
func requestLogger(c *fiber.Ctx) error {
	start := time.Now()
	log.Printf("Incoming %s %s from %s", c.Method(), c.OriginalURL(), c.IP())

	// Continue processing the request
	err := c.Next()

	// Log completion details
	duration := time.Since(start)
	log.Printf("Completed %d in %v", c.Response().StatusCode(), duration)

	return err
}
