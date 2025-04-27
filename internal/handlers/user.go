package handlers

import (
	"github.com/sirupsen/logrus"
	"mifi-bank/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	logger   *logrus.Logger
}

func NewUserHandler(repo *repository.UserRepository, logger *logrus.Logger) *UserHandler {
	return &UserHandler{userRepo: repo, logger: logger}
}
