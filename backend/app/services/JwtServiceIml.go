package services

import (
	"crawl-manager-backend/app/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTServiceImpl struct {
	SecretKey string
}

func NewJWTService(secretKey string) *JWTServiceImpl {
	return &JWTServiceImpl{SecretKey: secretKey}
}

func (s *JWTServiceImpl) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"iss":      "crawl-manager",
		"exp":      time.Now().Add(5 * 365 * 24 * time.Hour).Unix(), // 5 years from now
		"iat":      time.Now().Unix(),
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
