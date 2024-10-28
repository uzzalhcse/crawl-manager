package controllers

import (
	"crawl-manager-backend/app/http/responses"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/services"
	"github.com/gofiber/fiber/v2"
)

type ProxyController struct {
	Service *services.ProxyService
	*BaseController
}

func NewProxyController(service *services.ProxyService) *ProxyController {
	that := NewBaseController()
	return &ProxyController{Service: service, BaseController: that}
}

func (ctrl *ProxyController) Index(c *fiber.Ctx) error {
	proxyCollections, err := ctrl.Service.GetAllProxy()
	if len(proxyCollections) == 0 {
		proxyCollections = make([]models.Proxy, 0)
	}
	if err != nil {
		return responses.Error(c, err.Error())
	}
	return responses.Success(c, proxyCollections)
}
func (ctrl *ProxyController) Create(c *fiber.Ctx) error {
	var proxyCollection models.Proxy
	if err := c.BodyParser(&proxyCollection); err != nil {
		return responses.Error(c, err.Error())
	}

	err := ctrl.Service.Create(&proxyCollection)
	if err != nil {
		return responses.Error(c, err.Error())
	}

	return responses.Success(c, proxyCollection)
}
func (ctrl *ProxyController) Update(c *fiber.Ctx) error {
	proxyID := c.Params("id")
	var proxyCollection models.Proxy

	// Parse request body for fields to update
	if err := c.BodyParser(&proxyCollection); err != nil {
		return responses.Error(c, "Failed to parse request body: "+err.Error())
	}

	// Call the service to update the proxy
	if err := ctrl.Service.UpdateProxy(proxyID, &proxyCollection); err != nil {
		return responses.Error(c, "Failed to update proxy: "+err.Error())
	}

	return responses.Success(c, fiber.Map{"status": "proxy updated successfully"})
}
func (ctrl *ProxyController) Delete(c *fiber.Ctx) error {
	server := c.Params("server")
	err := ctrl.Service.Delete(server)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
func (ctrl *ProxyController) AssignProxies(c *fiber.Ctx) error {
	type AssignProxyRequest struct {
		SiteID          string `json:"site_id" form:"site_id"`
		NumberOfProxies int    `json:"number_of_proxies" form:"number_of_proxies"`
	}
	var request AssignProxyRequest
	if err := c.BodyParser(&request); err != nil {
		return responses.Error(c, "Failed to parse request body: "+err.Error())
	}

	// Call repository function to assign proxies
	if err := ctrl.Service.AssignProxiesToSite(request.SiteID, request.NumberOfProxies); err != nil {
		return responses.Error(c, "Failed to assign proxies: "+err.Error())
	}

	return responses.Success(c, fiber.Map{"status": "proxies assigned successfully"})
}
