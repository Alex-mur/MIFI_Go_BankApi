package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mifi-bank/internal/config"
	"mifi-bank/internal/models"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграция
	err = db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Card{},
		&models.Transaction{},
		&models.Credit{},
		&models.PaymentSchedule{},
	)

	// Включение расширения pgcrypto
	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")

	return db, err
}
