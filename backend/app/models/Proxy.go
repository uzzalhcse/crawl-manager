package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Proxy struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProxyID               string             `json:"proxy_id" bson:"proxy_id"` // Adjusted to match "id" from the response
	Server                string             `json:"server" bson:"server"`
	Username              string             `json:"username" bson:"username"`
	Password              string             `json:"password" bson:"password"`
	ProxyAddress          string             `json:"proxy_address" bson:"proxy_address"`
	Port                  int                `json:"port" bson:"port"`
	Valid                 bool               `json:"valid" bson:"valid"`
	LastVerification      time.Time          `json:"last_verification" bson:"last_verification"`
	CountryCode           string             `json:"country_code" bson:"country_code"`
	CityName              string             `json:"city_name" bson:"city_name"`
	ASNName               string             `json:"asn_name" bson:"asn_name"`
	ASNNumber             int                `json:"asn_number" bson:"asn_number"`
	HighCountryConfidence bool               `json:"high_country_confidence" bson:"high_country_confidence"`
	CreatedAt             time.Time          `json:"created_at" bson:"created_at"`
	ErrorLog              string             `json:"error_log" bson:"error_log"`
	SiteProxies           []SiteProxy        `json:"site_proxies" bson:"site_proxies"`
}

func (c *Proxy) GetTableName() string {
	return "proxies"
}
