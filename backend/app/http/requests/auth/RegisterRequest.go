// Http/Requests/register_request.go

package authrequests

import (
	"crawl-manager-backend/app/http/requests"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	*requests.Request
}
