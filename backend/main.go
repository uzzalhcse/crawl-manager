package main

import (
	"crawl-manager-backend/app/exceptions"
	"crawl-manager-backend/bootstrap"
	"crawl-manager-backend/routes"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	// Set up day-wise logging within month folders
	logFile := setupDailyLogFile()
	defer logFile.Close()

	// Add middleware to log HTTP requests
	app.App.Use(createRequestLogger(logFile))

	// Register routes
	routes.RegisterRoutes(app.App)

	// Launch the application in a goroutine
	go startApplication(app)

	// Graceful shutdown
	app.GracefulShutdown(func() {
		shutdownApplication(app)
	})
}

// setupDailyLogFile creates a log file with day-wise naming in month folders
func setupDailyLogFile() *os.File {
	// Create logs directory if it doesn't exist
	logBaseDir := "logs"
	if err := os.MkdirAll(logBaseDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Generate month folder name
	currentMonth := time.Now().Format("2006-01")
	monthDir := filepath.Join(logBaseDir, currentMonth)

	// Create month folder
	if err := os.MkdirAll(monthDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create month directory: %v", err)
	}

	// Generate daily log filename
	currentDay := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(monthDir, fmt.Sprintf("requests-%s.log", currentDay))

	// Open log file with append mode
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set output to both file and standard output
	log.SetOutput(logFile)
	return logFile
}

// createRequestLogger creates a middleware for detailed request logging
func createRequestLogger(logFile *os.File) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Collect detailed request information
		method := c.Method()
		path := c.OriginalURL()
		ip := c.IP()
		userAgent := c.Get("User-Agent")

		// Log incoming request with more details
		log.Printf("Incoming Request: Method=%s, Path=%s, IP=%s, UserAgent=%s",
			method, path, ip, userAgent)

		// Log headers (optional, can be verbose)
		headers := c.GetReqHeaders()
		headerLog := "Request Headers:\n"
		for key, value := range headers {
			headerLog += fmt.Sprintf("  %s: %s\n", key, value)
		}
		log.Println(headerLog)

		// Continue processing the request
		err := c.Next()

		// Log completion details
		duration := time.Since(start)
		log.Printf("Completed: Status=%d, Duration=%v, Method=%s, Path=%s",
			c.Response().StatusCode(), duration, method, path)

		return err
	}
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
