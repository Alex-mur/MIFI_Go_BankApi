package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/models"
	"mifi-bank/internal/repository"
	"net/http"
)

type CreditHandler struct {
	creditRepo  *repository.CreditRepository
	accountRepo *repository.AccountRepository
	logger      *logrus.Logger
}

func NewCreditHandler(creditRepo *repository.CreditRepository, accountRepo *repository.AccountRepository, logger *logrus.Logger) *CreditHandler {
	return &CreditHandler{creditRepo: creditRepo, accountRepo: accountRepo, logger: logger}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var input struct {
		AccountID    uuid.UUID `json:"account_id"`
		Amount       float64   `json:"amount"`
		InterestRate float64   `json:"interest_rate"`
		TermMonths   int       `json:"term_months"`
	}

	// Проверка владения счетом
	account, err := h.accountRepo.GetByID(input.AccountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "Invalid account", http.StatusForbidden)
		return
	}

	credit := &models.Credit{
		UserID:       userID,
		AccountID:    input.AccountID,
		Amount:       input.Amount,
		InterestRate: input.InterestRate,
		TermMonths:   input.TermMonths,
	}

	if err := h.creditRepo.CreateCredit(credit); err != nil {
		h.logger.Error("Credit creation failed: ", err)
		http.Error(w, "Credit creation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID := r.Context().Value("userID").(uuid.UUID)

	// Парсим creditID из URL
	vars := mux.Vars(r)
	creditID, err := uuid.Parse(vars["creditId"])
	if err != nil {
		h.logger.WithError(err).Error("Invalid credit ID format")
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	// 1. Получаем кредит из репозитория
	credit, err := h.creditRepo.GetByID(creditID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			http.Error(w, "Credit not found", http.StatusNotFound)
			return
		}
		h.logger.WithError(err).Error("Failed to get credit")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 2. Проверяем права доступа
	if credit.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// 3. Получаем график платежей
	schedule, err := h.creditRepo.GetPaymentSchedule(creditID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payment schedule")
		http.Error(w, "Failed to get payment schedule", http.StatusInternalServerError)
		return
	}

	// 4. Формируем ответ
	response := struct {
		CreditID uuid.UUID                `json:"credit_id"`
		Schedule []models.PaymentSchedule `json:"schedule"`
	}{
		CreditID: creditID,
		Schedule: schedule,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
