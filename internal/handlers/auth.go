package handlers

import (
	"encoding/json"
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/models"
	"mifi-bank/internal/repository"
	"net/http"
	"time"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	logger    *logrus.Logger
	jwtSecret []byte
}

func NewAuthHandler(repo *repository.UserRepository, logger *logrus.Logger, secret []byte) *AuthHandler {
	return &AuthHandler{
		userRepo:  repo,
		logger:    logger,
		jwtSecret: secret,
	}
}

// RegisterRequest структура запроса регистрации
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse структура ответа регистрации
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	AuthToken string `json:"auth_token,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode registration request")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	user := &models.User{Email: req.Email}

	// Создание пользователя
	if err := h.userRepo.CreateUser(user, req.Password); err != nil {
		switch err {
		case repository.ErrEmailExists:
			http.Error(w, "Email already registered", http.StatusConflict)
		case repository.ErrInvalidEmail:
			http.Error(w, "Invalid email format", http.StatusBadRequest)
		case repository.ErrWeakPassword:
			http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		default:
			h.logger.WithError(err).Error("Failed to create user")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Генерация JWT токена
	token, err := h.generateJWT(user.ID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate JWT token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		AuthToken: token,
	})
}

func (h *AuthHandler) generateJWT(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(h.jwtSecret)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode login request")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.Authenticate(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		h.logger.WithError(err).Error("Authentication error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate JWT token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		AuthToken: token,
	})
}
