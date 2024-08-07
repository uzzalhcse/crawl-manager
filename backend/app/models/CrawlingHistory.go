package models

type CrawlingHistory struct {
	SiteID       string          `json:"site_id" bson:"site_id"`
	Status       string          `json:"status" bson:"status"`               // running,stopped
	InstanceName string          `json:"instance_name" bson:"instance_name"` // running,stopped
	StartDate    string          `json:"start_date" bson:"start_date"`       // 2024-05-25
	EndDate      string          `json:"end_date" bson:"end_date"`           // 2024-08-20
	Site         *SiteCollection `json:"site" bson:"site"`
	Logs         string          `json:"logs" bson:"logs"`
}

func (c *CrawlingHistory) GetTableName() string {
	return "crawling_history"
}
