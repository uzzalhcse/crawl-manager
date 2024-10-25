package services

import (
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/repositories"
)

type ProxyService struct {
	Repository *repositories.Repository
}

func NewProxyService(repo *repositories.Repository) *ProxyService {
	return &ProxyService{Repository: repo}
}
func (s *ProxyService) GetAllProxy() ([]models.Proxy, error) {
	return s.Repository.GetAllProxy()
}
func (s *ProxyService) Create(proxy *models.Proxy) error {
	return s.Repository.CreateProxy(proxy)
}

func (s *ProxyService) Delete(server string) error {
	return s.Repository.DeleteProxy(server)
}
