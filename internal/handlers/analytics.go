package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/repository"
	"net/http"
	"strconv"
	"time"
)

type AnalyticsHandler struct {
	transRepo *repository.TransactionRepository
	logger    *logrus.Logger
}

func NewAnalyticsHandler(repo *repository.TransactionRepository, logger *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{transRepo: repo, logger: logger}
}

func (h *AnalyticsHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	// Параметры запроса
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	fromDate, _ := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	toDate, _ := time.Parse(time.RFC3339, r.URL.Query().Get("to"))

	transactions, total, err := h.transRepo.GetByUserID(
		userID,
		fromDate,
		toDate,
		page,
		pageSize,
	)

	if err != nil {
		h.logger.Error("Error fetching transactions: ", err)
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":     transactions,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}

	json.NewEncoder(w).Encode(response)
}
