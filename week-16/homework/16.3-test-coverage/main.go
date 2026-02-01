package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Ошибки валидации
var (
	ErrUsernameTooShort    = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong     = errors.New("username must not exceed 20 characters")
	ErrUsernameStartsDigit = errors.New("username cannot start with a digit")
	ErrUsernameInvalidChar = errors.New("username can only contain letters, digits and underscore")

	ErrEmailMissingAt     = errors.New("email must contain @")
	ErrEmailMultipleAt    = errors.New("email must contain exactly one @")
	ErrEmailInvalidFormat = errors.New("email has invalid format")

	ErrAgeNegative = errors.New("age cannot be negative")
	ErrAgeTooOld   = errors.New("age cannot exceed 150")

	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordNoDigit  = errors.New("password must contain at least one digit")
	ErrPasswordNoUpper  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower  = errors.New("password must contain at least one lowercase letter")
)

// User представляет данные пользователя для валидации
type User struct {
	Username string
	Email    string
	Age      int
	Password string
}

// Validator предоставляет методы валидации
type Validator struct{}

// NewValidator создаёт новый валидатор
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateUsername проверяет имя пользователя
// Правила:
// - Длина от 3 до 20 символов
// - Только буквы, цифры и подчёркивание
// - Не может начинаться с цифры
func (v *Validator) ValidateUsername(username string) error {
	if len(username) < 3 {
		return ErrUsernameTooShort
	}

	if len(username) > 20 {
		return ErrUsernameTooLong
	}

	runes := []rune(username)
	if unicode.IsDigit(runes[0]) {
		return ErrUsernameStartsDigit
	}

	for _, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return ErrUsernameInvalidChar
		}
	}

	return nil
}

// ValidateEmail проверяет email
// Правила:
// - Содержит ровно один символ @
// - До @ минимум 1 символ
// - После @ минимум 3 символа и точка
func (v *Validator) ValidateEmail(email string) error {
	atCount := strings.Count(email, "@")

	if atCount == 0 {
		return ErrEmailMissingAt
	}

	if atCount > 1 {
		return ErrEmailMultipleAt
	}

	parts := strings.Split(email, "@")
	local := parts[0]
	domain := parts[1]

	if len(local) < 1 {
		return ErrEmailInvalidFormat
	}

	if len(domain) < 3 || !strings.Contains(domain, ".") {
		return ErrEmailInvalidFormat
	}

	return nil
}

// ValidateAge проверяет возраст
// Правила:
// - От 0 до 150
func (v *Validator) ValidateAge(age int) error {
	if age < 0 {
		return ErrAgeNegative
	}

	if age > 150 {
		return ErrAgeTooOld
	}

	return nil
}

// ValidatePassword проверяет пароль
// Правила:
// - Минимум 8 символов
// - Содержит хотя бы одну цифру
// - Содержит хотя бы одну заглавную букву
// - Содержит хотя бы одну строчную букву
func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasDigit, hasUpper, hasLower bool

	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		}
	}

	if !hasDigit {
		return ErrPasswordNoDigit
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}

	if !hasLower {
		return ErrPasswordNoLower
	}

	return nil
}

// ValidateUser проверяет все поля пользователя
// Возвращает слайс всех найденных ошибок
func (v *Validator) ValidateUser(user User) []error {
	var errs []error

	if err := v.ValidateUsername(user.Username); err != nil {
		errs = append(errs, fmt.Errorf("username: %w", err))
	}

	if err := v.ValidateEmail(user.Email); err != nil {
		errs = append(errs, fmt.Errorf("email: %w", err))
	}

	if err := v.ValidateAge(user.Age); err != nil {
		errs = append(errs, fmt.Errorf("age: %w", err))
	}

	if err := v.ValidatePassword(user.Password); err != nil {
		errs = append(errs, fmt.Errorf("password: %w", err))
	}

	return errs
}

func main() {
	fmt.Println("=== Validator для тестирования покрытия ===")

	v := NewValidator()

	// Демонстрация валидации
	fmt.Println("\n--- ValidateUsername ---")
	usernames := []string{"user_123", "ab", "1user", "user@name"}
	for _, u := range usernames {
		err := v.ValidateUsername(u)
		if err != nil {
			fmt.Printf("ValidateUsername(%q) = error: %v\n", u, err)
		} else {
			fmt.Printf("ValidateUsername(%q) = OK\n", u)
		}
	}

	fmt.Println("\n--- ValidateEmail ---")
	emails := []string{"user@example.com", "invalid", "user@@test.com", "@test.com"}
	for _, e := range emails {
		err := v.ValidateEmail(e)
		if err != nil {
			fmt.Printf("ValidateEmail(%q) = error: %v\n", e, err)
		} else {
			fmt.Printf("ValidateEmail(%q) = OK\n", e)
		}
	}

	fmt.Println("\n--- ValidateUser ---")
	user := User{
		Username: "john_doe",
		Email:    "john@example.com",
		Age:      25,
		Password: "SecurePass123",
	}
	errs := v.ValidateUser(user)
	if len(errs) == 0 {
		fmt.Printf("User %+v is valid\n", user)
	} else {
		fmt.Printf("User %+v has errors:\n", user)
		for _, err := range errs {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println("\nНапиши тесты и проверь покрытие: go test -cover")
}
