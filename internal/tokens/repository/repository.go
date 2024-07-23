package repository

import (
	"account-service/internal/models"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	DB     *gorm.DB
	Tracer trace.Tracer
}

// NewRepository create new Repository
func NewRepository(db *gorm.DB, tr trace.Tracer) *Repository {
	return &Repository{
		DB:     db,
		Tracer: tr,
	}
}

// CreateToken create new token with specific variety
func (r *Repository) CreateToken(ctx context.Context, variety string, identity string) (*models.Token, error) {
	tr := r.Tracer
	_, span := tr.Start(ctx, "db-create-token")
	defer span.End()

	token := models.Token{
		Variety:  variety,
		Identity: identity,
	}
	result := r.DB.Create(&token)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.Create error: %w", result.Error)
	}

	return &token, nil
}

// GetTokenByID get token from db by id
func (r *Repository) GetTokenByID(ctx context.Context, id string, variety string) (*models.Token, error) {
	tr := r.Tracer
	_, span := tr.Start(ctx, "db-get-token-by-id")
	defer span.End()

	var resultToken models.Token
	result := r.DB.Where(&models.Token{
		ID:      id,
		Variety: variety,
	}).First(&resultToken)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.First error: %w", result.Error)
	}

	return &resultToken, nil
}

// SetUse set token is used
func (r *Repository) SetUse(token *models.Token) {
	token.LastUse = time.Now()
	r.DB.Save(token)
}

// Revoke set token is revoked
func (r *Repository) Revoke(token *models.Token) {
	token.IsRevoked = true
	r.DB.Save(token)
}
