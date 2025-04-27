package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PaymentSchedule struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreditID uuid.UUID `gorm:"type:uuid;not null;index"`
	DueDate  time.Time `gorm:"not null"`
	Amount   float64   `gorm:"type:decimal(15,2);not null"`
	Penalty  float64   `gorm:"type:decimal(15,2);default:0.0"`
	Status   string    `gorm:"type:varchar(20);default:'pending'"`
	PaidAt   *time.Time
}

func (p *PaymentSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
