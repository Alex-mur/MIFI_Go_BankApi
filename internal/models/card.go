package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Card struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Number          string    `gorm:"type:varchar(255);not null"`
	NumberHash      string    `gorm:"type:varchar(64);uniqueIndex"`
	LastFourDigits  string    `gorm:"type:varchar(4);not null;index"`
	Expiry          string    `gorm:"type:varchar(5);not null"`
	CVV             string    `gorm:"type:varchar(255);not null"`
	Status          string    `gorm:"type:varchar(20);default:'active'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	EncryptedNumber []byte         `gorm:"type:bytea;not null"`
	EncryptedExpiry []byte         `gorm:"type:bytea;not null"`
	NumberHMAC      []byte         `gorm:"type:bytea;not null"`
	ExpiryHMAC      []byte         `gorm:"type:bytea;not null"`
	CVVHash         string         `gorm:"not null"`

	User    User    `gorm:"foreignKey:UserID"`
	Account Account `gorm:"foreignKey:AccountID"`
}

func (c *Card) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}
