package models

import "time"

type CrawlingPerformance struct {
	SiteID            string    `json:"site_id" bson:"site_id"`
	InstanceName      string    `json:"instance_name" bson:"instance_name"`
	TotalRequests     int32     `json:"total_requests" bson:"total_requests"`
	FailedRequests    int32     `json:"failed_requests" bson:"failed_requests"`
	SuccessRate       int32     `json:"success_rate" bson:"success_rate"`
	ElapsedTime       int32     `json:"elapsed_time" bson:"elapsed_time"`
	RequestsPerMinute int32     `json:"requests_per_minute" bson:"requests_per_minute"`
	RequestsPerHour   int32     `json:"requests_per_hour" bson:"requests_per_hour"`
	RequestsPerDay    int32     `json:"requests_per_day" bson:"requests_per_day"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
}

func (c *CrawlingPerformance) GetTableName() string {
	return "crawling_performance"
}
