package interfaces

import (
	"account-service/internal/models"
	"context"
	"github.com/golang-jwt/jwt"
)

// TokenService token service interface
type TokenService interface {
	NewJWT(ctx context.Context, variety string, identity string, username string, extra jwt.MapClaims) (*models.JWT, error)
	ParseJWT(ctx context.Context, token string) (*models.JWT, error)
	ValidateRefreshJWT(ctx context.Context, token string) (*models.JWT, string, error)
	CreateAuthJWT(ctx context.Context, identity string, email string) (*models.JWT, *models.JWT, error)
	CreateAccessJWT(ctx context.Context, identity string, email string) (*models.JWT, *models.JWT, error)
	Validate(j *models.JWT, variety string) error
	Revoke(j *models.JWT)
}
