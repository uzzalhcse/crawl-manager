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
	"strconv"
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
	frequency := getAvailableFrequency(ctrl.Config, siteCollection.Frequency)
	siteCollection.Frequency = frequency
	if err := ctrl.Service.Create(&siteCollection); err != nil {
		return responses.Error(c, err.Error())
	}
	if err := ctrl.ProxyService.AssignProxiesToSite(siteCollection.SiteID, siteCollection.NumberOfProxies); err != nil {
		return responses.Error(c, "Failed to assign proxies: "+err.Error())
	}
	if siteCollection.Frequency != "" && ctrl.Config.App.Env == "production" {
		err := CreateOrUpdateSchedulerJob(ctrl.Config, siteCollection.Frequency, siteCollection.SiteID, false)
		if err != nil {
			return responses.Error(c, err.Error())
		}
	}

	return responses.Success(c, "Site created successfully")
}
func getAvailableFrequency(config *config.Config, frequency string) string {
	if frequency != "" && config.App.Env == "production" {
		// Retrieve existing jobs and find the next available time slot
		existingJobs, _ := findNextAvailableTimeSlot(config) // Assume this function lists all scheduler jobs
		nextAvailableSchedule := findNextAvailableSlot(existingJobs, frequency)
		// Use the next available schedule with a 5-minute gap
		if nextAvailableSchedule != "" {
			frequency = nextAvailableSchedule
		}
	}
	return frequency
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
	if !siteCollection.UseProxy {
		siteCollection.NumberOfProxies = 0
	}

	siteData, err := ctrl.Service.GetByID(siteID)
	if err != nil {
		return responses.Error(c, err.Error())
	}

	if siteData.Frequency != siteCollection.Frequency {
		frequency := getAvailableFrequency(ctrl.Config, siteCollection.Frequency)
		siteCollection.Frequency = frequency
	}
	// Update the site in the database
	err = ctrl.Service.Update(siteID, &siteCollection)
	if err != nil {
		return responses.Error(c, err.Error())
	}

	// Assign proxies to the site if necessary
	if err := ctrl.ProxyService.AssignProxiesToSite(siteCollection.SiteID, siteCollection.NumberOfProxies); err != nil {
		return responses.Error(c, "Failed to assign proxies: "+err.Error())
	}

	// Update or create the Cloud Scheduler job if frequency is specified and environment is production
	if siteCollection.Frequency != "" && siteData.Frequency != siteCollection.Frequency && ctrl.Config.App.Env == "production" {
		err := CreateOrUpdateSchedulerJob(ctrl.Config, siteCollection.Frequency, siteCollection.SiteID, true)
		if err != nil {
			return responses.Error(c, err.Error())
		}
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
func (ctrl *SiteCollectionController) FindNextAvailableTimeSlot(c *fiber.Ctx) error {
	existingJobs, err := findNextAvailableTimeSlot(ctrl.Config)
	if err != nil {
		return fmt.Errorf("error finding next available time slot: %v", err)
	}

	nextAvailableSchedule := findNextAvailableSlot(existingJobs, "0 18 * * 4")
	return responses.Success(c, fiber.Map{
		"existingJobs":          existingJobs,
		"nextAvailableSchedule": nextAvailableSchedule,
	})
}
func findNextAvailableTimeSlot(config *config.Config) ([]string, error) {
	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// URL to list jobs in the specified project and region
	url := fmt.Sprintf("https://cloudscheduler.googleapis.com/v1/projects/%s/locations/%s/jobs", config.Manager.ProjectID, config.Manager.Region)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request to list jobs: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to list jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(respBody))
	}

	// Parse the response to extract the job schedules
	var responseBody struct {
		Jobs []struct {
			Schedule string `json:"schedule"`
		} `json:"jobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	var schedulers []string
	for _, job := range responseBody.Jobs {
		schedulers = append(schedulers, job.Schedule)
	}
	return schedulers, nil
}

// Create or update the scheduler job with time slot gap
func CreateOrUpdateSchedulerJob(config *config.Config, frequency, siteName string, isUpdate bool) error {
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Define job details, including the calculated time slot
	jobRequestBody := map[string]interface{}{
		"name":     fmt.Sprintf("projects/%s/locations/%s/jobs/%s-job", config.Manager.ProjectID, config.Manager.Region, siteName),
		"schedule": frequency,
		"timeZone": "UTC",
		"httpTarget": map[string]interface{}{
			"uri":        fmt.Sprintf("%s/api/start-crawler/%s", config.Manager.ServerIP, siteName),
			"httpMethod": "GET",
			"headers": map[string]string{
				"Content-Type":  "application/json",
				"Authorization": config.Manager.BearerToken,
			},
		},
	}

	// Convert the job request body to JSON
	requestBody, err := json.Marshal(jobRequestBody)
	if err != nil {
		return fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Determine URL and HTTP method based on whether creating or updating the job
	var url string
	var method string
	if isUpdate {
		url = fmt.Sprintf("https://cloudscheduler.googleapis.com/v1/projects/%s/locations/%s/jobs/%s-job", config.Manager.ProjectID, config.Manager.Region, siteName)
		method = "PATCH"
	} else {
		url = fmt.Sprintf("https://cloudscheduler.googleapis.com/v1/projects/%s/locations/%s/jobs", config.Manager.ProjectID, config.Manager.Region)
		method = "POST"
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(respBody))
	}

	if isUpdate {
		fmt.Println("Cloud Scheduler job updated successfully with time slot gap")
	} else {
		fmt.Println("Cloud Scheduler job created successfully with time slot gap")
	}

	return nil
}

// Find the next available time slot in 5-minute increments within an hour
func findNextAvailableSlot(existingJobs []string, baseSchedule string) string {
	baseComponents := parseCronExpression(baseSchedule)

	// Start from the base hour and try up to 24 hours
	for hourOffset := 0; hourOffset < 24; hourOffset++ {
		// Calculate new hour, wrapping around at 24
		newHour := (parseHour(baseComponents[1]) + hourOffset) % 24

		// Get used slots for this specific hour
		modifiedComponents := make([]string, len(baseComponents))
		copy(modifiedComponents, baseComponents)
		modifiedComponents[1] = strconv.Itoa(newHour)

		usedSlots := getUsedSlotsDynamic(existingJobs, modifiedComponents)

		// Find the next available minute slot for this hour
		for minutes := 0; minutes < 60; minutes += 3 {
			if !usedSlots[minutes] {
				return formatScheduleWithMinuteAndHour(baseSchedule, minutes, strconv.Itoa(newHour))
			}
		}
	}

	// If no slot found after checking all hours, return original base schedule
	return baseSchedule
}

// Find the next available minute slot in 3-minute increments for the current hour
func findNextAvailableMinuteSlot(usedSlots map[int]bool, baseComponents []string) (int, string) {
	hour := parseHour(baseComponents[1])
	for i := 0; i < 60; i += 3 {
		if !usedSlots[i] {
			return i, strconv.Itoa(hour)
		}
	}
	return -1, strconv.Itoa(hour) // No available slots within the hour
}

// Parse the hour as an integer
func parseHour(hourStr string) int {
	hour, _ := strconv.Atoi(hourStr)
	return hour
}

// Parse cron expression into components
func parseCronExpression(schedule string) []string {
	return strings.Fields(schedule)
}

// Get dynamically used slots based on existing jobs and base schedule components
func getUsedSlotsDynamic(existingJobs []string, baseComponents []string) map[int]bool {
	usedSlots := make(map[int]bool)
	for _, job := range existingJobs {
		jobComponents := parseCronExpression(job)
		if matchesBaseSchedule(jobComponents, baseComponents) {
			minute, _ := strconv.Atoi(jobComponents[0])
			usedSlots[minute] = true
		}
	}
	return usedSlots
}

// Check if job components match the base schedule components except for the minute
func matchesBaseSchedule(jobComponents, baseComponents []string) bool {
	for i := 1; i < len(baseComponents); i++ {
		if baseComponents[i] != "*" && baseComponents[i] != jobComponents[i] {
			return false
		}
	}
	return true
}

// Format the cron schedule with specific minute and hour offsets
func formatScheduleWithMinuteAndHour(baseSchedule string, minute int, hour string) string {
	parts := strings.Fields(baseSchedule)
	parts[0] = strconv.Itoa(minute)
	parts[1] = hour
	return strings.Join(parts, " ")
}
