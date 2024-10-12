package authclient

import (
	"github.com/google/uuid"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
)

type Error struct {
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type AuthenticateResp struct {
	UserID uuid.UUID
	Claims auth.Claims
}
