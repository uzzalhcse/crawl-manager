package models

type SiteCollection struct {
	SiteID    string   `json:"site_id" bson:"site_id"`
	Url       string   `json:"url" bson:"url"`
	GitBranch string   `json:"git_branch" bson:"git_branch"`
	Status    string   `json:"status" bson:"status"`
	Frequency string   `json:"frequency" bson:"frequency"`
	VmConfig  VmConfig `json:"vm_config" bson:"vm_config"`
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
