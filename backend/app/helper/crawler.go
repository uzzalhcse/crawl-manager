package helper

import (
	"bytes"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GenerateBinaryBuild(siteCollection models.SiteCollection, config *config.Config) error {
	GitBranch := siteCollection.GitBranch
	if config.App.Env == "production" {
		GitBranch = "dev"
	}
	// Get the absolute path of the parent directory
	parentDir, err := filepath.Abs(config.Manager.DistDir)
	if err != nil {
		return fmt.Errorf("Error getting parent directory: %v", err)
	}

	files, err := os.ReadDir(config.Manager.AppsDir)
	if err != nil {
		return fmt.Errorf("Error reading directory: %v", err)
	}

	// Pull the latest changes from the remote repository
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git fetch origin && git reset --hard origin/%s", config.Manager.AppsDir, GitBranch))
	output, err := cmd.CombinedOutput()
	fmt.Println("git reset output: ", string(output))

	if err != nil {
		return fmt.Errorf("Error during git reset: %v\nOutput: %s", err, output)
	}

	cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && git checkout %s", config.Manager.AppsDir, GitBranch))
	output, err = cmd.CombinedOutput()
	fmt.Println("git checkout output: ", string(output))

	if err != nil {
		return fmt.Errorf("Error during git checkout: %v\nOutput: %s", err, output)
	}

	// Use --ff-only to avoid divergence issues during git pull
	cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull --ff-only origin %s", config.Manager.AppsDir, GitBranch))
	output, err = cmd.CombinedOutput()
	fmt.Println("git pull output: ", string(output))

	if err != nil {
		return fmt.Errorf("Error during git pull: %v\nOutput: %s", err, output)
	}
	// Small delay to ensure the filesystem reflects the latest changes
	time.Sleep(3 * time.Second)

	siteFound := false
	for _, file := range files {
		if file.IsDir() {
			dirname := file.Name()
			if siteCollection.SiteID == dirname {
				siteFound = true
				fmt.Printf("Generating Binary for: %s\n", dirname)
				outputPath := filepath.Join(parentDir, "dist", dirname)
				sourcePath := fmt.Sprintf("%s/%s", config.Manager.AppsDir, dirname)
				fmt.Println("sourcePath: ", sourcePath)
				fmt.Println("outputPath: ", outputPath)
				cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && go build -o %s", sourcePath, outputPath))
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("Error building site: %v\nOutput: %s", err, output)
				}
			}
		}
	}

	if !siteFound {
		return fmt.Errorf("Invalid site: %s", siteCollection.SiteID)
	}
	return nil
}

func CreateVM(siteCollection models.SiteCollection, config *config.Config) (string, error) {
	//projectID := "lazuli-venturas"
	//region := "asia-northeast1"
	dateTime := time.Now().Format("2006-01-02-15-04-05")
	sanitizedSiteID := strings.ReplaceAll(siteCollection.SiteID, "_", "-")

	instanceName := sanitizedSiteID + "-" + dateTime
	machineType := fmt.Sprintf("projects/%s/zones/%s/machineTypes/e2-custom-%d-%d", config.Manager.ProjectID, siteCollection.VmConfig.Zone, siteCollection.VmConfig.Cores, siteCollection.VmConfig.Memory)
	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	fmt.Println("machineType", machineType)
	// Construct the request body for creating the VM
	vmRequestBody := map[string]interface{}{
		"canIpForward":       false,
		"deletionProtection": false,
		"description":        "",
		"disks": []map[string]interface{}{
			{
				"autoDelete": true,
				"boot":       true,
				"deviceName": instanceName,
				"initializeParams": map[string]interface{}{
					"diskSizeGb":  siteCollection.VmConfig.DiskSize,
					"diskType":    fmt.Sprintf("projects/%s/zones/%s/diskTypes/pd-balanced", config.Manager.ProjectID, siteCollection.VmConfig.Zone),
					"sourceImage": fmt.Sprintf("projects/%s/global/images/boilerplate-for-ninjacrawler-pkg-disk-image", config.Manager.ProjectID),
				},
				"mode": "READ_WRITE",
				"type": "PERSISTENT",
			},
		},
		"displayDevice": map[string]bool{
			"enableDisplay": false,
		},
		"guestAccelerators":       []interface{}{},
		"instanceEncryptionKey":   map[string]string{},
		"keyRevocationActionType": "NONE",
		"labels": map[string]string{
			"goog-ec-src": "vm_add-rest",
		},
		"machineType": machineType,
		"metadata": map[string]interface{}{
			"items": []map[string]string{
				{
					"key":   "enable-osconfig",
					"value": "TRUE",
				},
				{
					"key":   "startup-script",
					"value": fmt.Sprintf("#! /bin/bash\nSiteID=\"%s\"\ninstanceName=\"%s\"\nulimit -n 1000000\ncd /root/apps\ncurl -O %s/binary/$SiteID\nchmod +x $SiteID\ncurl -s %s/api/site-secret/env/$SiteID > .env\nsudo ./$SiteID\ncurl %s/api/stop-crawler/$instanceName", siteCollection.SiteID, instanceName, config.Manager.ServerIP, config.Manager.ServerIP, config.Manager.ServerIP),
				},
			},
		},
		"name": instanceName,
		"networkInterfaces": []map[string]interface{}{
			{
				"accessConfigs": []map[string]string{
					{
						"name":        "External NAT",
						"networkTier": "PREMIUM",
					},
				},
				"stackType":  "IPV4_ONLY",
				"subnetwork": fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", config.Manager.ProjectID, config.Manager.Region),
			},
		},
		"params": map[string]interface{}{
			"resourceManagerTags": map[string]string{},
		},
		"reservationAffinity": map[string]string{
			"consumeReservationType": "ANY_RESERVATION",
		},
		"scheduling": map[string]interface{}{
			"automaticRestart":  true,
			"onHostMaintenance": "MIGRATE",
			"provisioningModel": "STANDARD",
		},
		"serviceAccounts": []map[string]interface{}{
			{
				"email": config.Manager.ServiceAccountEmail,
				"scopes": []string{
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
					"https://www.googleapis.com/auth/service.management.readonly",
					"https://www.googleapis.com/auth/servicecontrol",
					"https://www.googleapis.com/auth/trace.append",
				},
			},
		},
		"shieldedInstanceConfig": map[string]interface{}{
			"enableIntegrityMonitoring": true,
			"enableSecureBoot":          false,
			"enableVtpm":                true,
		},
		"tags": map[string]interface{}{
			"items": []string{"http-server", "https-server"},
		},
		"zone": fmt.Sprintf("projects/%s/zones/%s", config.Manager.ProjectID, siteCollection.VmConfig.Zone),
	}

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(vmRequestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Send the request to create the VM
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances", config.Manager.ProjectID, siteCollection.VmConfig.Zone)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to create instance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("unexpected response status: %s\n", resp.Status)
		fmt.Printf("Response Headers: %v\n", resp.Header)
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response Body: %s\n", string(respBody))
		return "", fmt.Errorf("unexpected response %s", string(respBody))
	}

	return instanceName, nil
}
