package controllers

import (
	"crawl-manager-backend/app/helper"
	"crawl-manager-backend/app/http/responses"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/repositories"
	"crawl-manager-backend/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os/exec"
	"time"
)

// TestController defines a controller for handling test-related requests
type ManagerController struct {
	*BaseController
	managerService *services.ManagerService
	siteService    *services.SiteCollectionService
}

func NewManagerController() *ManagerController {
	that := NewBaseController()

	// Manager controller
	return &ManagerController{
		BaseController: that,
		managerService: services.NewManagerService(repositories.NewRepository(that.DB)),
		siteService:    services.NewSiteCollectionService(repositories.NewRepository(that.DB)),
	}
}

func (that *ManagerController) Manager(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "Hello World",
		"status":  "Success",
	})
}

func (that *ManagerController) StopCrawler(c *fiber.Ctx) error {

	instanceName := c.Params("instanceName")

	crawlingHistory, err := that.siteService.GetCrawlerFromHistory(instanceName)
	if err != nil {
		return responses.Error(c, "Crawler not running")
	}

	command := []string{"gcloud", "compute", "instances", "stop", instanceName, "--zone", crawlingHistory.Site.VmConfig.Zone}
	res := ExecuteCommand(command[0], command[1:])

	err = that.siteService.UpdateCrawlingHistory(instanceName, map[string]interface{}{"status": "stopped", "logs": res, "end_date": time.Now().Format("2006-01-02 15:04:05")})
	if err != nil {
		return responses.Error(c, err.Error())
	}

	fmt.Println("cmd output: ", res)
	return nil
}

func (that *ManagerController) CrawlingHistoryLog(c *fiber.Ctx) error {

	instanceName := c.Params("instanceName")

	var crawlingPerformance models.CrawlingPerformance
	if err := c.BodyParser(&crawlingPerformance); err != nil {
		return responses.Error(c, err.Error())
	}
	crawlingPerformance.InstanceName = instanceName
	crawlingPerformance.CreatedAt = time.Now()
	if err := that.siteService.AddCrawlingPerformance(&crawlingPerformance); err != nil {
		return responses.Error(c, err.Error())
	}

	return nil
}
func (that *ManagerController) CrawlingSummary(c *fiber.Ctx) error {

	instanceName := c.Params("instanceName")

	var crawlingSummary models.CrawlingSummary
	if err := c.BodyParser(&crawlingSummary); err != nil {
		return responses.Error(c, err.Error())
	}
	crawlingSummary.InstanceName = instanceName
	crawlingSummary.CreatedAt = time.Now()
	if err := that.siteService.AddCrawlingSummary(&crawlingSummary); err != nil {
		return responses.Error(c, err.Error())
	}

	return nil
}
func ExecuteCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return string(output)
}
func (that *ManagerController) StartCrawler(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := that.siteService.GetByID(siteID)
	if err != nil {
		return responses.Error(c, err.Error())
	}
	if siteCollection.Status != "active" {
		return responses.Error(c, fmt.Sprintf("%s Crawler is not active", siteCollection.SiteID))
	}

	err = helper.GenerateBinaryBuild(*siteCollection, that.Config)
	if err != nil {
		return responses.Error(c, err.Error())
	}
	fmt.Println("Creating VM for: ", siteID)
	instanceName, instanceID, err := helper.CreateVM(*siteCollection, that.Config)
	if err != nil {
		return responses.Error(c, err.Error())
	}
	fmt.Println("instanceName: ", instanceName)
	logUrl := fmt.Sprintf(
		`https://console.cloud.google.com/logs/query;query=resource.type%%3D%%22gce_instance%%22%%0Aresource.labels.instance_id%%3D%%22%s%%22%%0AlogName%%3D%%22projects%%2F%s%%2Flogs%%2Fgoogle_metadata_script_runner%%22;?project=%s`,
		instanceID, that.Config.Manager.ProjectID, that.Config.Manager.ProjectID,
	)
	err = that.siteService.CreateCrawlingHistory(&models.CrawlingHistory{
		SiteID:       siteID,
		Status:       "running",
		InstanceName: instanceName,
		InstanceID:   instanceID,
		Site:         siteCollection,
		LogUrl:       logUrl,
		StartDate:    time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return responses.Error(c, err.Error())
	}
	return responses.Success(c, "Crawler Started Successfully")
}
func (that *ManagerController) BuildCrawler(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	siteCollection, err := that.siteService.GetByID(siteID)
	if err != nil {
		return responses.Error(c, err.Error())
	}
	if siteCollection.Status != "active" {
		return responses.Error(c, fmt.Sprintf("%s Crawler is not active", siteCollection.SiteID))
	}

	err = helper.GenerateBinaryBuild(*siteCollection, that.Config)
	if err != nil {
		return responses.Error(c, err.Error())
	}
	return responses.Success(c, "Crawler Build Successfully")
}
func (that *ManagerController) CrawlingHistory(c *fiber.Ctx) error {

	crawlingHistory, err := that.siteService.GetCrawlingHistory()
	if err != nil {
		return responses.Error(c, err.Error())
	}

	return responses.Success(c, crawlingHistory)
}

func (that *ManagerController) GetCrawlingSummary(c *fiber.Ctx) error {
	instanceName := c.Params("instanceName")
	crawlingSummary, err := that.siteService.GetCrawlingSummary(instanceName)
	if err != nil {
		return responses.Error(c, err.Error())
	}

	return responses.Success(c, crawlingSummary)
}
