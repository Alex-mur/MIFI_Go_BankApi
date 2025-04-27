package repository

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mifi-bank/internal/models"
	"mifi-bank/pkg/crypto"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *models.Card, pubKey *openpgp.Entity, hmacSecret []byte) error {
	// Шифрование данных
	encryptedNumber, err := crypto.EncryptPGP([]byte(card.Number), pubKey)
	if err != nil {
		return err
	}
	encryptedExpiry, err := crypto.EncryptPGP([]byte(card.Expiry), pubKey)
	if err != nil {
		return err
	}

	// Генерация HMAC
	card.NumberHMAC = crypto.GenerateHMAC(encryptedNumber, hmacSecret)
	card.ExpiryHMAC = crypto.GenerateHMAC(encryptedExpiry, hmacSecret)

	// Хеширование CVV
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(card.CVV), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	card.CVVHash = string(cvvHash)

	// Сохранение в БД
	return r.db.Create(card).Error
}

func (r *CardRepository) GetByID(id uuid.UUID, pgpPrivKey *openpgp.Entity, hmacKey []byte) (*models.Card, error) {
	var card *models.Card
	err := r.db.First(&card, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	// Проверка HMAC перед дешифровкой
	if !crypto.VerifyHMAC(card.EncryptedNumber, card.NumberHMAC, hmacKey) {
		return nil, errors.New("invalid HMAC for card number")
	}

	if !crypto.VerifyHMAC(card.EncryptedExpiry, card.ExpiryHMAC, hmacKey) {
		return nil, errors.New("invalid HMAC for expiry date")
	}

	// Дешифровка данных
	decryptedNum, err := crypto.DecryptPGP(card.EncryptedNumber, pgpPrivKey)
	decryptedExp, err := crypto.DecryptPGP(card.EncryptedExpiry, pgpPrivKey)

	card.Number = string(decryptedNum)
	card.Expiry = string(decryptedExp)

	return card, nil
}

func (r *CardRepository) GetByUserID(userID uuid.UUID) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Preload("Account").
		Where("user_id = ?", userID).
		Find(&cards).
		Error
	return cards, err
}

func (r *CardRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Card{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *CardRepository) BlockCard(id uuid.UUID) error {
	return r.UpdateStatus(id, "blocked")
}

func (r *CardRepository) GetByLastFourDigits(lastFour string) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("last_four_digits = ?", lastFour).
		Find(&cards).
		Error
	return cards, err
}

func (r *CardRepository) GetByAccountID(accountID uuid.UUID) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("account_id = ?", accountID).
		Find(&cards).
		Error
	return cards, err
}

func (r *CardRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Мягкое удаление
		return tx.Model(&models.Card{}).
			Where("id = ?", id).
			Update("deleted_at", gorm.Expr("NOW()")).
			Error
	})
}

// Дополнительные методы безопасности

func (r *CardRepository) GetForUpdate(id uuid.UUID) (*models.Card, error) {
	var card models.Card
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&card, "id = ?", id).
		Error
	return &card, err
}

func (r *CardRepository) ExistsByNumberHash(hash string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Card{}).
		Where("number_hash = ?", hash).
		Count(&count).
		Error
	return count > 0, err
}

// Ошибки репозитория
var (
	ErrDuplicateCard = errors.New("card with this number already exists")
	ErrCardNotFound  = errors.New("card not found")
)
