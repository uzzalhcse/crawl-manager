package authrequests

import (
	"crawl-manager-backend/app/http/requests"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
	*requests.Request
}
