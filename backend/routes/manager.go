package routes

import (
	"crawl-manager-backend/app/http/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetUpManagerRoutes(api fiber.Router) {
	managerController := controllers.NewManagerController()
	// Define routes
	api.Get("/start-crawler/:SiteID", managerController.StartCrawler)
	api.Get("/build-crawler/:SiteID", managerController.BuildCrawler)
	api.Get("/stop-crawler/:instanceName", managerController.StopCrawler)
	api.Get("/crawling-history", managerController.CrawlingHistory)
	api.Post("/add-crawler-logs/:instanceName", managerController.CrawlingHistoryLog)
	api.Post("/add-crawler-summary/:instanceName", managerController.CrawlingSummary)
	api.Get("/crawler-summary/:instanceName", managerController.GetCrawlingSummary)
}
