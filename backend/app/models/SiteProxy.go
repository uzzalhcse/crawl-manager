package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SiteProxy struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SiteID  string             `json:"site_id" bson:"site_id"`
	ProxyID primitive.ObjectID `json:"proxy_id" bson:"proxy_id"` // Reference to Proxy
}

func (c *SiteProxy) GetTableName() string {
	return "site_proxies"
}
