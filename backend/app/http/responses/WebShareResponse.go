package responses

import (
	"crawl-manager-backend/app/models"
	"fmt"
	"time"
)

type ProxyListResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []ProxyResponse `json:"results"`
}
type ProxyResponse struct {
	ID                    string    `json:"id"`
	Username              string    `json:"username"`
	Password              string    `json:"password"`
	ProxyAddress          string    `json:"proxy_address"`
	Port                  int       `json:"port"`
	Valid                 bool      `json:"valid"`
	LastVerification      time.Time `json:"last_verification"`
	CountryCode           string    `json:"country_code"`
	CityName              string    `json:"city_name"`
	ASNName               string    `json:"asn_name"`
	ASNNumber             int       `json:"asn_number"`
	HighCountryConfidence bool      `json:"high_country_confidence"`
	CreatedAt             time.Time `json:"created_at"`
}

func (resp *ProxyListResponse) ConvertToProxy() []models.Proxy {
	var proxy []models.Proxy
	for i := 0; i < len(resp.Results); i++ {
		proxy = append(proxy, resp.Results[i].Convert())
	}
	return proxy
}
func (response *ProxyResponse) Convert() models.Proxy {
	return models.Proxy{
		ProxyID:               response.ID,
		Server:                fmt.Sprintf("%s:%d", response.ProxyAddress, response.Port),
		Username:              response.Username,
		Password:              response.Password,
		ProxyAddress:          response.ProxyAddress,
		Port:                  response.Port,
		Valid:                 response.Valid,
		LastVerification:      response.LastVerification,
		CountryCode:           response.CountryCode,
		CityName:              response.CityName,
		ASNName:               response.ASNName,
		ASNNumber:             response.ASNNumber,
		HighCountryConfidence: response.HighCountryConfidence,
		CreatedAt:             response.CreatedAt,
		ErrorLog:              "",
		SiteProxies:           nil,
	}
}
