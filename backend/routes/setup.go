package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// RegisterRoutes registers all routes
func RegisterRoutes(router fiber.Router) {
	// Enable CORS with default options to allow all origins
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                      // Allow all origins
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS", // Allowed methods
	}))
	router.Static("/public", "./storage")
	router.Static("/binary", "./dist")
	web := router.Group("")
	SetUpWebRoutes(web)

	api := router.Group("/api")
	auth := router.Group("/api/auth")
	SetUpPublicRoutes(api)
	SetUpAuthRoutes(auth)
	SetUpApiRoutes(api)
	SetUpManagerRoutes(api)

}
