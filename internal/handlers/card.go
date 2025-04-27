package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/config"
	"mifi-bank/internal/models"
	"mifi-bank/internal/repository"
	"mifi-bank/pkg/card"
	"mifi-bank/pkg/crypto"
	"net/http"
)

type CardHandler struct {
	cardRepo    *repository.CardRepository
	accountRepo *repository.AccountRepository
	logger      *logrus.Logger
	cfg         *config.Config
}

func NewCardHandler(cardRepo *repository.CardRepository,
	accountRepo *repository.AccountRepository,
	logger *logrus.Logger,
	cfg *config.Config) *CardHandler {
	return &CardHandler{cardRepo: cardRepo, accountRepo: accountRepo, logger: logger, cfg: cfg}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var input struct {
		AccountID uuid.UUID `json:"account_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.WithError(err).Error("Failed to decode card request")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 1. Проверка владения счетом
	account, err := h.accountRepo.GetByID(input.AccountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "Account not found or access denied", http.StatusForbidden)
		return
	}

	// 2. Генерация данных карты
	cardData := card.GenerateCardData()

	// 3. Загрузка PGP ключей
	pubKey, err := crypto.LoadPublicKey(h.cfg.PGPPublicKey)
	if err != nil {
		h.logger.WithError(err).Error("Failed to load PGP public key")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 4. Создание модели карты
	newCard := &models.Card{
		UserID:    userID,
		AccountID: input.AccountID,
		Number:    cardData.Number,
		Expiry:    cardData.Expiry,
		CVV:       cardData.CVV,
	}

	// 5. Вызов метода Create репозитория
	if err := h.cardRepo.Create(newCard, pubKey, h.cfg.HMACSecret); err != nil {
		h.logger.WithError(err).Error("Failed to create card")
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	// 6. Формирование ответа (без чувствительных данных)
	response := map[string]interface{}{
		"id":         newCard.ID,
		"account_id": newCard.AccountID,
		"last_four":  newCard.LastFourDigits,
		"expiry":     newCard.Expiry,
		"created_at": newCard.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err1 := json.NewEncoder(w).Encode(response)
	if err1 != nil {
		return
	}
}

func (h *CardHandler) generateCardData() {

}
