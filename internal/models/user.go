package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User model for user
type User struct {
	ID string `gorm:"primaryKey;type:uuid" json:"id"`

	Password      string `json:"password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `gorm:"unique" json:"email"`
	EmailVerified bool   `json:"email_verified"`
	DepartmentID  uint32 `json:"department_id"`
	RoleID        uint32 `json:"role_id"`
	IsRegistered  bool   `json:"is_registered"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// BeforeCreate - create new uuid
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}
