package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/repository"
	"net/http"
)

type TransactionHandler struct {
	transRepo   *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	logger      *logrus.Logger
}

func NewTransactionHandler(repo *repository.TransactionRepository, logger *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{transRepo: repo, logger: logger}
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var input struct {
		FromAccountID uuid.UUID `json:"from_account_id"`
		ToAccountID   uuid.UUID `json:"to_account_id"`
		Amount        float64   `json:"amount"`
	}

	// Проверка владения счетом-отправителем
	fromAccount, err := h.accountRepo.GetByID(input.FromAccountID)
	if err != nil || fromAccount.UserID != userID {
		http.Error(w, "Invalid sender account", http.StatusForbidden)
		return
	}

	// Выполнение перевода
	if err := h.transRepo.Transfer(
		input.FromAccountID,
		input.ToAccountID,
		input.Amount,
	); err != nil {
		h.logger.Error("Transfer failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
