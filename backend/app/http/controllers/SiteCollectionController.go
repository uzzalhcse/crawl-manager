package controllers

import (
	"bytes"
	"crawl-manager-backend/app/http/responses"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/services"
	"crawl-manager-backend/config"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

type SiteCollectionController struct {
	Service      *services.SiteCollectionService
	ProxyService *services.ProxyService
	*BaseController
}

func NewSiteCollectionController(service *services.SiteCollectionService, proxyService *services.ProxyService) *SiteCollectionController {
	that := NewBaseController()
	return &SiteCollectionController{BaseController: that, Service: service, ProxyService: proxyService}
}

func (ctrl *SiteCollectionController) Index(c *fiber.Ctx) error {
	siteCollections, err := ctrl.Service.GetAllSiteCollections()
	if err != nil {
		return responses.Error(c, err.Error())
	}
	return responses.Success(c, siteCollections)
}
func (ctrl *SiteCollectionController) Create(c *fiber.Ctx) error {
	var siteCollection models.SiteCollection
	if err := c.BodyParser(&siteCollection); err != nil {
		return responses.Error(c, err.Error())
	}

	if err := ctrl.Service.Create(&siteCollection); err != nil {
		return responses.Error(c, err.Error())
	}
	if err := ctrl.ProxyService.AssignProxiesToSite(siteCollection.SiteID, siteCollection.NumberOfProxies); err != nil {
		return responses.Error(c, "Failed to assign proxies: "+err.Error())
	}
	if siteCollection.Frequency != "" && ctrl.Config.App.Env == "production" {
		err := CreateSchedulerJob(ctrl.Config, siteCollection.Frequency, siteCollection.SiteID)
		if err != nil {
			return responses.Error(c, err.Error())
		}
	}

	return responses.Success(c, "Site created successfully")
}
func CreateSchedulerJob(config *config.Config, frequency, siteName string) error {
	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error retrieving access token: %v\n", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Define the request body for the Cloud Scheduler job
	jobRequestBody := map[string]interface{}{
		"name":     fmt.Sprintf("projects/%s/locations/%s/jobs/%s-job", config.Manager.ProjectID, config.Manager.Region, siteName),
		"schedule": frequency,
		"timeZone": "UTC",
		"httpTarget": map[string]interface{}{
			"uri":        fmt.Sprintf("%s/api/start-crawler/%s", config.Manager.ServerIP, siteName),
			"httpMethod": "GET",
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(jobRequestBody)
	if err != nil {
		return fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Send the request to create the Cloud Scheduler job
	url := fmt.Sprintf("https://cloudscheduler.googleapis.com/v1/projects/%s/locations/%s/jobs", config.Manager.ProjectID, config.Manager.Region)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to create scheduler job: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("unexpected response status: %s\n", resp.Status)
		fmt.Printf("Response Headers: %v\n", resp.Header)
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response Body: %s\n", string(respBody))
		return fmt.Errorf("unexpected response %s", string(respBody))
	}

	fmt.Println("Cloud Scheduler job created successfully")
	return nil
}
func UpdateSchedulerJob(config *config.Config, frequency, siteName string) error {
	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Define the request body for updating the Cloud Scheduler job
	jobRequestBody := map[string]interface{}{
		"schedule": frequency,
		"timeZone": "UTC",
		"httpTarget": map[string]interface{}{
			"uri":        fmt.Sprintf("%s/api/start-crawler/%s", config.Manager.ServerIP, siteName),
			"httpMethod": "GET",
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(jobRequestBody)
	if err != nil {
		return fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Construct the URL to update the job
	url := fmt.Sprintf("https://cloudscheduler.googleapis.com/v1/projects/%s/locations/%s/jobs/%s", config.Manager.ProjectID, config.Manager.Region, siteName)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to update scheduler job: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(respBody))
	}

	fmt.Println("Cloud Scheduler job updated successfully")
	return nil
}
func (ctrl *SiteCollectionController) GetByID(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := ctrl.Service.GetByID(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(siteCollection)
}

func (ctrl *SiteCollectionController) Update(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	var siteCollection models.SiteCollection
	if err := c.BodyParser(&siteCollection); err != nil {
		return responses.Error(c, err.Error())
	}

	err := ctrl.Service.Update(siteID, &siteCollection)

	if err := ctrl.ProxyService.AssignProxiesToSite(siteCollection.SiteID, siteCollection.NumberOfProxies); err != nil {
		return responses.Error(c, "Failed to assign proxies: "+err.Error())
	}
	if err != nil {
		return responses.Error(c, err.Error())
	}

	return responses.Success(c, "Site updated successfully")
}

func (ctrl *SiteCollectionController) Delete(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	err := ctrl.Service.Delete(siteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
