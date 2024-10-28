package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Proxy struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Server      string             `json:"server" bson:"server"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password" bson:"password"`
	Status      string             `json:"status" bson:"status"`
	SiteProxies []SiteProxy        `json:"site_proxies" bson:"site_proxies" gorm:"foreignKey:ProxyID"`
}

func (c *Proxy) GetTableName() string {
	return "proxies"
}
