package service

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"mifi-bank/internal/config"
	"mifi-bank/internal/models"
)

type EmailService struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		dialer:   gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass),
		from:     cfg.EmailFrom,
		fromName: cfg.EmailFromName,
	}
}

func (s *EmailService) SendPaymentNotification(userEmail string, payment *models.PaymentSchedule, credit *models.Credit) error {
	// Генерация HTML-письма
	tpl, err := template.ParseFiles("templates/payment_notification.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := struct {
		Payment *models.PaymentSchedule
		Credit  *models.Credit
	}{payment, credit}

	if err := tpl.Execute(&body, data); err != nil {
		return err
	}

	// Формирование письма
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.from, s.fromName)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Payment Notification")
	m.SetBody("text/html", body.String())

	// Асинхронная отправка
	go func() {
		if err := s.dialer.DialAndSend(m); err != nil {
			logrus.Errorf("Failed to send email: %v", err)
		}
	}()

	return nil
}
