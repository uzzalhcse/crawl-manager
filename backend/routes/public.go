package routes

import (
	"crawl-manager-backend/app/http/controllers"
	"crawl-manager-backend/app/repositories"
	"crawl-manager-backend/app/services"
	"crawl-manager-backend/bootstrap"
	"github.com/gofiber/fiber/v2"
)

func SetUpPublicRoutes(api fiber.Router) {

	repo := repositories.NewRepository(bootstrap.App().DB)
	secretCollectionService := services.NewSecretCollectionService(repo)
	proxyService := services.NewProxyService(repo)
	siteCollectionService := services.NewSiteCollectionService(repo)

	managerController := controllers.NewManagerController()
	proxyController := controllers.NewProxyController(proxyService)
	secretCollectionController := controllers.NewSecretCollectionController(secretCollectionService)
	siteCollectionController := controllers.NewSiteCollectionController(siteCollectionService, proxyService)

	// Test controller
	testService := services.NewTestService(repo)
	testController := controllers.NewTestController(testService)

	/*
		Public api
	*/

	api.Get("/site-secret/env/:siteID", secretCollectionController.GetEnvBySite)

	// Proxy Management
	proxy := api.Group("/proxy")
	proxy.Get("/sync", proxyController.Sync)
	//proxy.Post("/", proxyController.Create)
	proxy.Get("/", proxyController.Index)
	proxy.Get("/:siteID", proxyController.Show)
	proxy.Post("/stop", proxyController.StopProxy)
	proxy.Put("/:id", proxyController.Update)
	proxy.Delete("/:server", proxyController.Delete)
	//proxy.Post("/allocate-proxies", proxyController.AssignProxies)

	// Define test routes
	api.Get("/", testController.Test)
	api.Get("/test", testController.GetAllHandler)
	api.Get("/start-crawler/:SiteID/:zone", testController.StartCrawler)
	api.Get("/test/available-slots", siteCollectionController.FindNextAvailableTimeSlot) // test api

	// manager routes
	api.Get("/stop-crawler/:instanceName", managerController.StopCrawler)
	api.Post("/add-crawler-logs/:instanceName", managerController.CrawlingHistoryLog)
	api.Post("/add-crawler-summary/:instanceName", managerController.CrawlingSummary)

}
