package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	FromAccountID uuid.UUID `gorm:"type:uuid"`
	ToAccountID   uuid.UUID `gorm:"type:uuid;not null"`
	Amount        float64   `gorm:"type:decimal(15,2);not null"`
	Status        string    `gorm:"type:varchar(20);not null"` // pending, completed, failed
	CreatedAt     time.Time
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	return
}
