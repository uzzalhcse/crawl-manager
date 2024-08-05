package services

import (
	"crawl-manager-backend/app/repositories"
)

type ManagerService struct {
	Repository *repositories.Repository
}

func NewManagerService(repo *repositories.Repository) *ManagerService {
	return &ManagerService{Repository: repo}
}
