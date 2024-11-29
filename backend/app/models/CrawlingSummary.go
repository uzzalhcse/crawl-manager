package models

import "time"

type CrawlingSummary struct {
	SiteID         string          `json:"site_id" bson:"site_id"`
	InstanceName   string          `json:"instance_name" bson:"instance_name"`
	CollectionName string          `json:"collection_name" bson:"collection_name"`
	DataCount      int32           `json:"data_count" bson:"data_count"`
	ErrorCount     int32           `json:"error_count" bson:"error_count"`
	Errors         []CrawlingError `json:"errors" bson:"errors"`
	CreatedAt      time.Time       `json:"created_at" bson:"created_at"`
}

func (c *CrawlingSummary) GetTableName() string {
	return "crawling_summary"
}

type CrawlingError struct {
	Url          string `json:"url" bson:"url"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
}
