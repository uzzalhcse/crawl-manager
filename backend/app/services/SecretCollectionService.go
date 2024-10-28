package services

import (
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/repositories"
)

type SecretCollectionService struct {
	Repository *repositories.Repository
}

func NewSecretCollectionService(repo *repositories.Repository) *SecretCollectionService {
	return &SecretCollectionService{Repository: repo}
}
func (s *SecretCollectionService) GetAllSiteSecret() ([]models.SiteSecret, error) {
	return s.Repository.GetAllSiteSecretCollections()
}
func (s *SecretCollectionService) GetAllGlobalSecret() ([]models.GlobalSecret, error) {
	return s.Repository.GetAllGlobalSecretCollections()
}
func (s *SecretCollectionService) Create(siteSecret *models.SiteSecret) error {
	return s.Repository.CreateSecretCollection(siteSecret)
}

func (s *SecretCollectionService) GetByID(siteID string) (*models.SiteSecret, error) {
	return s.Repository.GetSiteSecretCollectionByID(siteID)
}

func (s *SecretCollectionService) Delete(siteID string) error {
	return s.Repository.DeleteSiteCollection(siteID)
}
