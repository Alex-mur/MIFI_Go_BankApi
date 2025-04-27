package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Credit struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount       float64   `gorm:"type:decimal(15,2);not null"`
	InterestRate float64   `gorm:"type:decimal(5,2);not null"`
	TermMonths   int       `gorm:"not null"`
	Status       string    `gorm:"type:varchar(20);default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Payments     []PaymentSchedule
}

func (c *Credit) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	return
}
