package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mifi-bank/internal/models"
	"time"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) Transfer(from, to uuid.UUID, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var fromAccount models.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fromAccount, "id = ?", from).Error; err != nil {
			return err
		}

		if fromAccount.Balance < amount {
			return errors.New("insufficient funds")
		}

		if err := tx.Model(&models.Account{}).
			Where("id = ?", from).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Account{}).
			Where("id = ?", to).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return tx.Create(&models.Transaction{
			FromAccountID: from,
			ToAccountID:   to,
			Amount:        amount,
			Status:        "completed",
		}).Error
	})
}

func (r *TransactionRepository) GetByUserID(
	userID uuid.UUID,
	fromDate time.Time,
	toDate time.Time,
	page int,
	pageSize int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	// Получаем все счета пользователя
	var accountIDs []uuid.UUID
	if err := r.db.Model(&models.Account{}).
		Where("user_id = ?", userID).
		Pluck("id", &accountIDs).Error; err != nil {
		return nil, 0, err
	}

	// Базовый запрос
	query := r.db.Model(&models.Transaction{}).
		Where("(from_account_id IN (?) OR to_account_id IN (?))", accountIDs, accountIDs).
		Order("created_at DESC")

	// Фильтрация по дате
	if !fromDate.IsZero() {
		query = query.Where("created_at >= ?", fromDate)
	}
	if !toDate.IsZero() {
		query = query.Where("created_at <= ?", toDate)
	}

	// Получаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// Выполняем запрос
	err := query.
		Preload("FromAccount").
		Preload("ToAccount").
		Find(&transactions).Error

	return transactions, total, err
}
