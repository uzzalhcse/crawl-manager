package config

import "github.com/spf13/viper"

type ManagerConfig struct {
	AppsDir                   string
	DistDir                   string
	ProjectID                 string
	Region                    string
	ServerIP                  string
	ServiceAccountEmail       string
	ServiceAccountCredentials string
	GcpBucketName             string
	WebShareApiKey            string
	WebShareApiUrl            string
}

func loadManagerConfig() ManagerConfig {
	return ManagerConfig{
		AppsDir:                   viper.GetString("APPS_DIR"),
		DistDir:                   viper.GetString("DIST_DIR"),
		ProjectID:                 viper.GetString("PROJECT_ID"),
		Region:                    viper.GetString("REGION"),
		ServerIP:                  viper.GetString("SERVER_IP"),
		ServiceAccountEmail:       viper.GetString("SERVICE_ACCOUNT_EMAIL"),
		ServiceAccountCredentials: viper.GetString("GCP_SERVICE_ACCOUNT"),
		GcpBucketName:             viper.GetString("GCP_BUCKET_NAME"),
		WebShareApiKey:            viper.GetString("WEB_SHARE_API_KEY"),
		WebShareApiUrl:            viper.GetString("WEB_SHARE_API_URL"),
	}
}
