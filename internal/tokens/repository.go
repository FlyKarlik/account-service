package tokens

import (
	"account-service/internal/models"
	"context"
)

// Repository interface for repository
type Repository interface {
	CreateToken(ctx context.Context, variety string, identity string) (*models.Token, error)
	GetTokenByID(ctx context.Context, id string, variety string) (*models.Token, error)
	SetUse(token *models.Token)
	Revoke(token *models.Token)
}
