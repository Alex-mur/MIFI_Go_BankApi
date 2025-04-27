package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/models"
	"mifi-bank/internal/repository"
	"net/http"
)

type AccountHandler struct {
	accountRepo *repository.AccountRepository
	logger      *logrus.Logger
}

func NewAccountHandler(repo *repository.AccountRepository, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{accountRepo: repo, logger: logger}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var input struct {
		Currency string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("Invalid request: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	account := &models.Account{
		UserID:   userID,
		Currency: input.Currency,
	}

	if err := h.accountRepo.Create(account); err != nil {
		h.logger.Error("Error creating account: ", err)
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
