package helper

import (
	"bytes"
	"crawl-manager-backend/app/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GenerateBinaryBuild(SiteID string) error {

	appsDir := "/root/ninja-combined-crawler/apps"
	distDir := "/root//crawl-manager/backend"

	//appsDir := "/home/uzzal/Workplace/Lazuli/ninja-combined-crawler/apps"
	//distDir := "/home/uzzal/Workplace/github/crawl-manager-backend"

	// Get the absolute path of the parent directory
	parentDir, err := filepath.Abs(distDir)
	if err != nil {
		return fmt.Errorf("Error getting parent directory: %v", err)
	}

	files, err := os.ReadDir(appsDir)
	if err != nil {
		return fmt.Errorf("Error reading directory:", err)
	}

	siteFound := false
	for _, file := range files {
		if file.IsDir() {
			dirname := file.Name()
			if SiteID == dirname {
				siteFound = true
				fmt.Printf("Generating Binary for: %s\n", dirname)
				outputPath := filepath.Join(parentDir, "dist", dirname)
				sourcePath := fmt.Sprintf("%s/%s", appsDir, dirname)
				fmt.Println("sourcePath: ", sourcePath)
				fmt.Println("outputPath: ", outputPath)
				cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull && go build -o %s", sourcePath, outputPath))
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("Error building site: %v\nOutput: %s", err, output)
				}
			}
		}
	}
	if !siteFound {
		return fmt.Errorf("invalid site: %s", SiteID)
	}
	return nil
}

func CreateVM(siteCollection models.SiteCollection) (string, error) {
	projectID := "lazuli-venturas"
	date := time.Now().Format("2006-01-02")
	instanceName := siteCollection.SiteID + "-" + date
	machineType := fmt.Sprintf("projects/%s/zones/%s/machineTypes/e2-custom-%d-%d", projectID, siteCollection.VmConfig.Zone, siteCollection.VmConfig.Cores, siteCollection.VmConfig.Memory)
	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

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
					"diskType":    fmt.Sprintf("projects/lazuli-venturas/zones/%s/diskTypes/pd-balanced", siteCollection.VmConfig.Zone),
					"sourceImage": "projects/lazuli-venturas/global/images/reference-crawler-disk-image",
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
					"value": "#! /bin/bash\\nSiteID=\\\"sumitool\\\"\\ninstanceName=\\\"sumitool\\\"\\nulimit -n 1000000\\ncd /root\\ncurl -O http://35.243.109.168:8080/binary/$SiteID\\nchmod +x $SiteID\\ncurl -s http://35.243.109.168:8080/api/site-secret/$SiteID  > .env\\nsudo ./$SiteID\\ncurl http://35.243.109.168:8080/api/stop-crawler/$instanceName",
				},
				{
					"key":   "ssh-keys",
					"value": "haris_dipto:ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBJKv0j3qfsCVR3vWVNs94qfTiKX1/yYKorP/zAcm+Xh+gy/4v4P5cA9ZLZvGjXqVjUwEu1yJSxe9P/fD8LxE0ec= google-ssh {\"userName\":\"haris.dipto@lazuli.ninja\",\"expireOn\":\"2024-03-04T06:10:06+0000\"}\nharis_dipto:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCE+ZM7cG+i5CAYAMxLPGUG6snCFqIc/OC7f/ttm4+DxvdRDAlmPlYrhMTBMVfga2CP/Idq6Gc69nONrRVGSr8ZiGWqMHiyxZaQ6HuTViatY+8TtO6PIKcr59PbiMaPdehSSxB5C1lbXEtzSK0mpGek8yLg8yHrPD0uN5xfcJai4fI6bPydK6GBM/LXjo+pnc0/B7FBcpCUzpaXZgxB2X6I9eQdB7f80zoX+e00yhr6CP2ZdQQnpXA5E67iIx3FlBXqf+ORWwg1mzC6lt8YeXYHvGXH04AtNZBjrTuMFc/zScWZjlFgBBwLFdHNUmaSKC/yIPKL8Xhi7RtplcYYl1KV google-ssh {\"userName\":\"haris.dipto@lazuli.ninja\",\"expireOn\":\"2024-03-04T06:10:24+0000\"}",
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
				"subnetwork": "projects/lazuli-venturas/regions/asia-northeast1/subnetworks/default",
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
				"email": "845643578999-compute@developer.gserviceaccount.com",
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
		"zone": fmt.Sprintf("projects/%s/zones/%s", projectID, siteCollection.VmConfig.Zone),
	}

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(vmRequestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Send the request to create the VM
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances", projectID, siteCollection.VmConfig.Zone)
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
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return instanceName, nil
}
