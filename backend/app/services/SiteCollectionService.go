package services

import (
	"crawl-manager-backend/app/models"
	"crawl-manager-backend/app/repositories"
)

type SiteCollectionService struct {
	Repository *repositories.Repository
}

func NewSiteCollectionService(repo *repositories.Repository) *SiteCollectionService {
	return &SiteCollectionService{Repository: repo}
}
func (s *SiteCollectionService) GetAllSiteCollections() ([]models.SiteCollection, error) {
	return s.Repository.GetAllSiteCollections()
}
func (s *SiteCollectionService) Create(siteCollection *models.SiteCollection) error {
	return s.Repository.CreateSiteCollection(siteCollection)
}

func (s *SiteCollectionService) CreateCrawlingHistory(crawlingHistory *models.CrawlingHistory) error {
	return s.Repository.CreateCrawlingHistory(crawlingHistory)
}

func (s *SiteCollectionService) GetByID(siteID string) (*models.SiteCollection, error) {
	return s.Repository.GetSiteCollectionByID(siteID)
}
func (s *SiteCollectionService) GetCrawlingHistoryByID(siteID string, runningOnly bool) ([]models.CrawlingHistory, error) {
	return s.Repository.GetCrawlingHistoryByID(siteID, runningOnly)
}
func (s *SiteCollectionService) GetCrawlerFromHistory(instanceName string) (*models.CrawlingHistory, error) {
	return s.Repository.GetCrawlerFromHistory(instanceName)
}
func (s *SiteCollectionService) GetCrawlingHistory() ([]models.CrawlingHistory, error) {
	return s.Repository.GetCrawlingHistory()
}

func (s *SiteCollectionService) Update(siteID string, update map[string]interface{}) error {
	return s.Repository.UpdateSiteCollection(siteID, update)
}
func (s *SiteCollectionService) UpdateCrawlingHistory(instanceName string, update map[string]interface{}) error {
	return s.Repository.UpdateCrawlingHistory(instanceName, update)
}

func (s *SiteCollectionService) Delete(siteID string) error {
	return s.Repository.DeleteSiteCollection(siteID)
}
