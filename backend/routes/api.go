package routes

import (
	"crawl-manager-backend/app/http/controllers"
	"crawl-manager-backend/app/http/middleware"
	"crawl-manager-backend/app/repositories"
	"crawl-manager-backend/app/services"
	"crawl-manager-backend/bootstrap"
	"github.com/gofiber/fiber/v2"
)

func SetUpApiRoutes(api fiber.Router) {
	// Initialize repositories and services
	repo := repositories.NewRepository(bootstrap.App().DB)
	siteCollectionService := services.NewSiteCollectionService(repo)
	collectionService := services.NewCollectionService(repo)
	urlCollectionService := services.NewUrlCollectionService(repo)
	secretCollectionService := services.NewSecretCollectionService(repo)
	proxyService := services.NewProxyService(repo)

	// Initialize controllers
	siteCollectionController := controllers.NewSiteCollectionController(siteCollectionService, proxyService)
	collectionController := controllers.NewCollectionController(collectionService)
	urlCollectionController := controllers.NewUrlCollectionController(urlCollectionService)
	secretCollectionController := controllers.NewSecretCollectionController(secretCollectionService)
	proxyController := controllers.NewProxyController(proxyService)

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

	/*
		Private routes
	*/
	api.Group("", middleware.Auth())
	// SiteCollection routes
	site := api.Group("/site")
	site.Get("/", siteCollectionController.Index)
	site.Post("/", siteCollectionController.Create)
	site.Get("/:siteID", siteCollectionController.GetByID)
	site.Put("/:siteID", siteCollectionController.Update)
	site.Delete("/:siteID", siteCollectionController.Delete)

	// Collection routes
	collection := api.Group("/collection")
	collection.Get("/", collectionController.Index)
	collection.Post("/", collectionController.Create)
	collection.Get("/:collectionID", collectionController.GetByID)
	collection.Put("/:collectionID", collectionController.Update)
	collection.Delete("/:collectionID", collectionController.Delete)

	// UrlCollection routes
	url := api.Group("/urlcollections")
	url.Post("/", urlCollectionController.Create)
	url.Get("/:collectionID", urlCollectionController.GetByID)
	url.Put("/:collectionID", urlCollectionController.Update)
	url.Delete("/:collectionID", urlCollectionController.Delete)

	// SiteSecretCollection routes
	secret := api.Group("/site-secret")
	secret.Post("/", secretCollectionController.Create)
	secret.Get("/:siteID", secretCollectionController.GetByID)
	secret.Delete("/:siteID", secretCollectionController.Delete)

}
