package services

import (
	"crawl-manager-backend/app/http/responses"
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/repositories"
	"crawl-manager-backend/bootstrap"
	"encoding/json"
	"fmt"
	"log"
)

type ProxyService struct {
	Repository        *repositories.Repository
	HttpClientService HttpClientService
}

func NewProxyService(repo *repositories.Repository) *ProxyService {
	return &ProxyService{Repository: repo, HttpClientService: NewHttpClientService()}
}
func (s *ProxyService) SyncProxy() ([]models.Proxy, error) {

	var proxies []models.Proxy
	apiUrl := fmt.Sprintf("%s/%s", bootstrap.App().Config.Manager.WebShareApiUrl, "proxy/list/?mode=direct&page=1&page_size=50")

	headers := make(map[string]string)
	headers["Authorization"] = bootstrap.App().Config.Manager.WebShareApiKey
	for {
		var response responses.ProxyListResponse
		apiResponse, err := s.HttpClientService.DoRequest("GET", apiUrl, nil, headers)
		if err != nil {
			return proxies, err
		}
		// Unmarshal body into ProxyListResponse
		if err = json.Unmarshal(apiResponse, &response); err != nil {
			return proxies, err
		}

		proxies = append(proxies, response.ConvertToProxy()...)
		if err := s.Repository.UpdateProxies(proxies); err != nil {
			log.Printf("Error updating proxies: %v", err)
		}
		if response.Next != "" {
			apiUrl = response.Next
		} else {
			break
		}
	}
	return proxies, nil

}
func (s *ProxyService) GetAllProxy() ([]models.Proxy, error) {
	return s.Repository.GetAllProxy()
}
func (s *ProxyService) GetProxyBySiteID(siteID string) ([]models.Proxy, error) {
	return s.Repository.GetSiteProxiesBySiteID(siteID)
}
func (s *ProxyService) FindProxy(id string) (*models.Proxy, error) {
	return s.Repository.FindProxy(id)
}

//	func (s *ProxyService) Create(proxy *models.Proxy) error {
//		return s.Repository.CreateProxy(proxy)
//	}
func (s *ProxyService) UpdateProxy(id string, proxy *models.Proxy) error {
	proxy.SiteProxies = nil
	return s.Repository.UpdateProxy(id, proxy)
}

func (s *ProxyService) Delete(server string) error {
	return s.Repository.DeleteProxy(server)
}
func (s *ProxyService) AssignProxiesToSite(siteID string, proxyCount int) error {
	return s.Repository.AssignProxiesToSite(siteID, proxyCount)
}
