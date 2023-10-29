package telebot

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidExpDate          = errors.New("🙅 Ваша карта просрочена!")                               //nolint
	ErrInvalidCardNumberLength = errors.New("🙅 Количество цифр в номере карты должно быть равно 16!") //nolint
	ErrInvalidData             = errors.New("🙅 Неверные данные! Попробуйте ещё раз.")                 //nolint
)

// verifyCVV verifies the CVV (Card Verification Value) and returns it as a uint64.
func (b *Bot) verifyCVV(msg string) (uint64, error) {
	// Remove spaces from the message.
	msg = strings.ReplaceAll(msg, " ", "")

	// Check if the length is not 3 digits.
	if len(msg) != 3 {
		return 0, ErrInvalidData
	}

	// Parse the CVV into a uint64.
	cvv, err := strconv.ParseUint(msg, 10, 64)
	if err != nil {
		return 0, ErrInvalidData
	}

	return cvv, nil
}

// verifyCardHolder verifies the cardholder's name and returns it as a string.
func (b *Bot) verifyCardHolder(msg string) (string, error) {
	// Split the name into parts.
	parts := strings.Split(msg, " ")

	// Check if there are not exactly two parts.
	if len(parts) != 2 {
		return "", ErrInvalidData
	}

	return msg, nil
}

// verifyExpDate verifies the card expiration date and returns it as a time.Time.
func (b *Bot) verifyExpDate(msg string) (time.Time, error) {
	// Remove spaces from the message.
	msg = strings.ReplaceAll(msg, " ", "")

	// Split the date into parts.
	parts := strings.Split(msg, "/")

	// Check if there are not exactly two parts.
	if len(parts) != 2 {
		return time.Time{}, ErrInvalidData
	}

	// Parse the month.
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, ErrInvalidData
	}

	// Check for a valid month (1-12).
	if month < 1 || month > 12 {
		return time.Time{}, ErrInvalidData
	}

	// Parse the year.
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, ErrInvalidData
	}

	// Check for an invalid year based on the current year.
	if year < time.Now().Year()%100 {
		return time.Time{}, ErrInvalidExpDate
	}

	// Check if the month and year are not in the past.
	if year == time.Now().Year()%100 && time.Month(month) <= time.Now().Month() {
		return time.Time{}, ErrInvalidExpDate
	}

	// Create a time.Time representing the expiration date.
	exp := time.Date(year+2000, time.Month(month), 0, 0, 0, 0, 0, time.Local) //nolint

	return exp, nil
}

// verifyCardNumber verifies the card number and returns it as a uint64.
func (b *Bot) verifyCardNumber(msg string) (uint64, error) {
	// Remove spaces from the message.
	msg = strings.ReplaceAll(msg, " ", "")

	// Check if the length is not 16 digits.
	if len(msg) != 16 {
		return 0, ErrInvalidCardNumberLength
	}

	// Parse the card number into a uint64.
	num, err := strconv.ParseUint(msg, 10, 64)
	if err != nil {
		return 0, ErrInvalidData
	}

	return num, nil
}
