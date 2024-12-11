package routes

import (
	"crawl-manager-backend/app/http/controllers"
	"crawl-manager-backend/app/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetUpManagerRoutes(api fiber.Router) {
	managerController := controllers.NewManagerController()

	// public routes
	api.Get("/stop-crawler/:instanceName", managerController.StopCrawler)
	api.Post("/add-crawler-logs/:instanceName", managerController.CrawlingHistoryLog)
	api.Post("/add-crawler-summary/:instanceName", managerController.CrawlingSummary)

	manager := api.Group("", middleware.Auth())
	// Define routes
	manager.Get("/start-crawler/:SiteID", managerController.StartCrawler)
	manager.Get("/build-crawler/:SiteID", managerController.BuildCrawler)
	manager.Get("/crawling-history", managerController.CrawlingHistory)
	manager.Get("/crawler-summary/:instanceName", managerController.GetCrawlingSummary)

}
