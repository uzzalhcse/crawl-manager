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

func (ctrl *ProxyController) Delete(c *fiber.Ctx) error {
	server := c.Params("server")
	err := ctrl.Service.Delete(server)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
