package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Balance   float64   `gorm:"type:decimal(15,2);default:0.0"`
	Currency  string    `gorm:"size:3;not null"`
	CreatedAt time.Time
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	return
}
