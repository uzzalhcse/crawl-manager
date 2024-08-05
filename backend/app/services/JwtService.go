package services

import "crawl-manager-backend/app/models"

type JWTService interface {
	GenerateToken(user *models.User) (string, error)
}
