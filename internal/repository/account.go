package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mifi-bank/internal/models"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "id = ?", id).Error
	return &account, err
}

func (r *AccountRepository) UpdateBalance(id uuid.UUID, amount float64) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}
