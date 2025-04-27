package scheduler

import (
	"mifi-bank/internal/repository"
	"mifi-bank/internal/service"
	"time"
)

type CreditScheduler struct {
	creditRepo   *repository.CreditRepository
	emailService *service.EmailService
}

func NewScheduler(creditRepo *repository.CreditRepository, emailService *service.EmailService) *CreditScheduler {
	return &CreditScheduler{creditRepo: creditRepo, emailService: emailService}
}

func (s *CreditScheduler) Start() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.creditRepo.ProcessDailyPayments(s.emailService)
			if err != nil {
				return
			}
		}
	}
}
