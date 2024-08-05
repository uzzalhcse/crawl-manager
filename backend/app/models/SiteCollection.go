package models

type SiteCollection struct {
	SiteID               string   `json:"site_id" bson:"site_id"`
	Name                 string   `json:"name" bson:"name"`
	Url                  string   `json:"url" bson:"url"`
	Status               string   `json:"status" bson:"status"`
	NoOfCrawlingPerMonth int      `json:"no_of_crawling_per_month" bson:"no_of_crawling_per_month"`
	VmConfig             VmConfig `json:"vm_config" bson:"vm_config"`
}

func (c *SiteCollection) GetTableName() string {
	return "sites"
}

type VmConfig struct {
	Cores    int    `json:"cores" bson:"cores"`
	Memory   int    `json:"memory" bson:"memory"`
	DiskSize int    `json:"disk" bson:"disk"`
	Zone     string `json:"zone" bson:"zone"`
}
