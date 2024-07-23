package models

type Role struct {
	ID uint32 `gorm:"primaryKey;type=uint32" json:"id"`

	Name string `gorm:"unique" json:"name"`
}