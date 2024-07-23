package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Token model for token
type Token struct {
	ID string `gorm:"primaryKey;type:uuid" json:"id"`

	Identity  string    `json:"identity"`
	Variety   string    `json:"variety"`
	IsRevoked bool      `json:"is_revoked"`
	LastUse   time.Time `json:"last_use"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate add uuid to id
func (token *Token) BeforeCreate(tx *gorm.DB) (err error) {
	token.ID = uuid.NewString()
	return
}
