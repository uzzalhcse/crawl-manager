package routes

import (
	"crawl-manager-backend/app/http/controllers"
	"crawl-manager-backend/app/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetUpManagerRoutes(api fiber.Router) {
	managerController := controllers.NewManagerController()

	manager := api.Group("", middleware.Auth())
	// Define routes
	manager.Get("/start-crawler/:SiteID", managerController.StartCrawler)
	manager.Get("/build-crawler/:SiteID", managerController.BuildCrawler)
	manager.Get("/crawling-history", managerController.CrawlingHistory)
	manager.Get("/crawler-summary/:instanceName", managerController.GetCrawlingSummary)

}
