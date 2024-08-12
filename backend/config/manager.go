package config

import "github.com/spf13/viper"

type ManagerConfig struct {
	AppsDir             string
	DistDir             string
	ProjectID           string
	Region              string
	ServerIP            string
	ServiceAccountEmail string
}

func loadManagerConfig() ManagerConfig {
	return ManagerConfig{
		AppsDir:             viper.GetString("APPS_DIR"),
		DistDir:             viper.GetString("DIST_DIR"),
		ProjectID:           viper.GetString("PROJECT_ID"),
		Region:              viper.GetString("REGION"),
		ServerIP:            viper.GetString("SERVER_IP"),
		ServiceAccountEmail: viper.GetString("SERVICE_ACCOUNT_EMAIL"),
	}
}
