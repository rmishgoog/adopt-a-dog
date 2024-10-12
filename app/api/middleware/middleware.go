// This package provides the 'app' layer middleware which is protocol agnostic.

package middleware

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
)

type Handler func(context.Context) error

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
)

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}
