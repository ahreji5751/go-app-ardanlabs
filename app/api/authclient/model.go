package authclient

import (
	"github.com/ardanlabs/service/app/api/auth"
	"github.com/google/uuid"
)

type Authorize struct {
	UserID uuid.UUID
	Claims auth.Claims
	Rule   string
}

type AuthenticateResp struct {
	UserID uuid.UUID
	Claims auth.Claims
}
