package card

import (
	"fmt"
	"math/rand"
	"time"
)

type CardData struct {
	Number string // Format: 4242424242424242
	Expiry string // Format: MM/YY
	CVV    string // 3 digits
}

func GenerateCardData() CardData {
	return CardData{
		Number: generateCardNumber(),
		Expiry: generateExpiryDate(),
		CVV:    generateCVV(),
	}
}

// Генерация номера карты с валидной контрольной суммой (Luhn algorithm)
func generateCardNumber() string {
	// Префикс для Visa
	prefix := "4"

	// Генерация 15 цифр (prefix + 14 случайных + контрольная сумма)
	rand.Seed(time.Now().UnixNano())
	partial := prefix
	for i := 0; i < 14; i++ {
		partial += fmt.Sprintf("%d", rand.Intn(10))
	}

	// Вычисление контрольной цифры
	checkDigit := luhnCheckDigit(partial)
	return partial + fmt.Sprintf("%d", checkDigit)
}

func luhnCheckDigit(number string) int {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		sum += digit
		alternate = !alternate
	}

	return (10 - (sum % 10)) % 10
}

// Генерация срока действия (MM/YY)
func generateExpiryDate() string {
	now := time.Now()
	year := now.Year() + 3 // Срок действия +3 года
	month := rand.Intn(12) + 1

	return fmt.Sprintf("%02d/%02d", month, year%100)
}

// Генерация CVV (3 цифры)
func generateCVV() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
