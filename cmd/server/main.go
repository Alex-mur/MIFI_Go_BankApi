package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"mifi-bank/internal/config"
	"mifi-bank/internal/handlers"
	"mifi-bank/internal/middleware"
	"mifi-bank/internal/repository"
	"mifi-bank/internal/scheduler"
	"mifi-bank/internal/service"
	"net/http"
)

func main() {
	cfg := config.Load()

	// Инициализация БД
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	creditRepo := repository.NewCreditRepository(db)
	cardRepo := repository.NewCardRepository(db)

	// Инициализация сервисов
	emailService := service.NewEmailService(cfg)
	creditScheduler := scheduler.NewScheduler(creditRepo, emailService)
	go creditScheduler.Start()

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(
		userRepo,
		logrus.StandardLogger(),
		cfg.JWTSecret,
	)
	accountHandler := handlers.NewAccountHandler(accountRepo, logrus.StandardLogger())
	cardHandler := handlers.NewCardHandler(cardRepo, accountRepo, logrus.StandardLogger(), cfg)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo, logrus.StandardLogger())
	analyticsHandler := handlers.NewAnalyticsHandler(transactionRepo, logrus.StandardLogger())
	creditHandler := handlers.NewCreditHandler(creditRepo, accountRepo, logrus.StandardLogger())

	// Маршрутизация
	// Публичные маршруты
	r := mux.NewRouter()
	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Защищенные маршруты
	secured := r.PathPrefix("/api/user").Subrouter()
	secured.Use(middleware.JWTAuth(cfg.JWTSecret))

	secured.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	secured.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	secured.HandleFunc("/transfer", transactionHandler.Transfer).Methods("POST")
	secured.HandleFunc("/analytics", analyticsHandler.GetTransactions).Methods("GET")
	secured.HandleFunc("/credits/create", creditHandler.CreateCredit).Methods("POST")
	secured.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetPaymentSchedule).Methods("GET")

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
