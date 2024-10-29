package controllers

import (
	"crawl-manager-backend/app/http/responses"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type SecretCollectionController struct {
	Service *services.SecretCollectionService
	*BaseController
}

func NewSecretCollectionController(service *services.SecretCollectionService) *SecretCollectionController {
	that := NewBaseController()
	return &SecretCollectionController{Service: service, BaseController: that}
}

func (ctrl *SecretCollectionController) Index(c *fiber.Ctx) error {
	siteCollections, err := ctrl.Service.GetAllGlobalSecret()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(siteCollections)
}
func (ctrl *SecretCollectionController) Create(c *fiber.Ctx) error {
	var secretCollection models.SiteSecret
	if err := c.BodyParser(&secretCollection); err != nil {
		return responses.Error(c, err.Error())
	}

	err := ctrl.Service.Create(&secretCollection)
	if err != nil {
		return responses.Error(c, err.Error())
	}

	return responses.Success(c, secretCollection)
}

func (ctrl *SecretCollectionController) GetByID(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := ctrl.Service.GetByID(siteID)
	if err != nil {
		return responses.Success(c, fiber.Map{
			"site_id": siteID,
			"secrets": "",
		})
	}

	return responses.Success(c, siteCollection)
}
func (ctrl *SecretCollectionController) GetEnvBySite(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, _ := ctrl.Service.GetByID(siteID)

	// Define default data
	defaultData := map[string]interface{}{
		"PROJECT_ID":               ctrl.Config.Manager.ProjectID,
		"SERVER_IP":                ctrl.Config.Manager.ServerIP,
		"DB_USERNAME":              "lazuli",
		"DB_PASSWORD":              "x1RWo6cqFtHiaAHce5HB",
		"DB_HOST":                  "localhost",
		"DB_PORT":                  27017,
		"USER_AGENT":               "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"APP_ENV":                  ctrl.Config.App.Env,
		"API_USERNAME":             "lazuli",
		"API_PASSWORD":             "ninja",
		"GCP_BUCKET_NAME":          ctrl.Config.Manager.GcpBucketName,
		"GCP_CREDENTIALS_PATH":     "/root/apps/gcp-file-upload-key.json",
		"GCP_LOG_CREDENTIALS_PATH": "/root/apps/gcp-log-key.json",
		"GCP_SERVICE_ACCOUNT":      ctrl.Config.Manager.ServiceAccountCredentials,
		"BIGQUERY_DATASET":         "data_source",
		"UPLOAD_LOGS":              "false",
		"PROXY_SERVERS":            "",
		"LOG_TO_GCP":               "true",
		"ENABLE_FILE_LOGGING":      "false",
	}
	if siteCollection != nil && siteCollection.Secrets != nil {
		// Merge siteCollection.Secrets into defaultData
		for key, value := range siteCollection.Secrets {
			defaultData[key] = value
		}
	}

	// Create envData string
	var envData string
	for key, value := range defaultData {
		envData += fmt.Sprintf("%s=%v\n", key, value)
	}

	return c.SendString(envData)
}

func (ctrl *SecretCollectionController) Delete(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	err := ctrl.Service.Delete(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
