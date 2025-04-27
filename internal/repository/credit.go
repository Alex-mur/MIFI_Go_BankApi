package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"mifi-bank/internal/models"
	"mifi-bank/internal/service"
	"time"
)

type CreditRepository struct {
	db *gorm.DB
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

func (r *CreditRepository) GetByID(id uuid.UUID) (*models.Credit, error) {
	var credit models.Credit
	err := r.db.Preload("Payments").
		First(&credit, "id = ?", id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return &credit, err
}

func NewCreditRepository(db *gorm.DB) *CreditRepository {
	return &CreditRepository{db: db}
}

func (r *CreditRepository) CreateCredit(credit *models.Credit) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Расчет графика платежей
		schedule, err := r.calculatePaymentSchedule(credit)
		if err != nil {
			return err
		}

		if err := tx.Create(credit).Error; err != nil {
			return err
		}

		for _, payment := range schedule {
			payment.CreditID = credit.ID
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *CreditRepository) calculatePaymentSchedule(credit *models.Credit) ([]models.PaymentSchedule, error) {
	// Аннуитетный платеж
	monthlyRate := credit.InterestRate / 100 / 12
	annuity := (credit.Amount * monthlyRate) / (1 - math.Pow(1+monthlyRate, float64(-credit.TermMonths)))

	schedule := make([]models.PaymentSchedule, credit.TermMonths)
	date := time.Now().AddDate(0, 1, 0) // Первый платеж через месяц

	for i := 0; i < credit.TermMonths; i++ {
		schedule[i] = models.PaymentSchedule{
			DueDate: date.AddDate(0, i, 0),
			Amount:  math.Round(annuity*100) / 100,
			Status:  "pending",
		}
	}

	return schedule, nil
}

func (r *CreditRepository) GetPaymentSchedule(creditID uuid.UUID) ([]models.PaymentSchedule, error) {
	var payments []models.PaymentSchedule

	err := r.db.
		Where("credit_id = ?", creditID).
		Order("due_date ASC").
		Find(&payments).
		Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *CreditRepository) ProcessDailyPayments(emailService *service.EmailService) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var payments []models.PaymentSchedule
		now := time.Now().UTC().Truncate(24 * time.Hour)

		// Находим просроченные платежи
		if err := tx.Where("status = 'pending' AND due_date < ?", now).Find(&payments).Error; err != nil {
			return err
		}

		for _, payment := range payments {
			var credit models.Credit
			if err := tx.First(&credit, payment.CreditID).Error; err != nil {
				continue
			}

			// Проверка баланса
			var account models.Account
			if err := tx.First(&account, credit.AccountID).Error; err != nil {
				continue
			}

			amountToPay := payment.Amount + payment.Penalty

			if account.Balance >= amountToPay {
				// Списание средств
				if err := tx.Model(&account).Update("balance", gorm.Expr("balance - ?", amountToPay)).Error; err != nil {
					continue
				}

				// Обновление статуса платежа
				payment.Status = "paid"
				payment.PaidAt = &now

				// Получение данных пользователя
				var user models.User
				if err := tx.Model(&models.User{}).
					Joins("JOIN credits ON credits.user_id = users.id").
					Where("credits.id = ?", credit.ID).
					First(&user).Error; err == nil {

					// Отправка уведомления
					err := emailService.SendPaymentNotification(user.Email, &payment, &credit)
					if err != nil {
						return err
					}
				}
			} else {
				// Начисление штрафа
				payment.Penalty += payment.Amount * 0.10
				payment.Status = "overdue"
			}

			if err := tx.Save(&payment).Error; err != nil {
				continue
			}
		}
		return nil
	})
}
