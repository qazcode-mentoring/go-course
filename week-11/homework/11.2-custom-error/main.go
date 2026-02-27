package main

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

// Error реализует интерфейс error
func (e ValidationError) Error() string {
	// TODO: реализуй метод
	// Формат: "поле 'Field': Message"
	return fmt.Sprintf("поле '%v': %v", e.Field, e.Message)
}

// NotFoundError представляет ошибку "не найдено"
type NotFoundError struct {
	Resource string
	ID       int
}

// Error реализует интерфейс error
func (e NotFoundError) Error() string {
	// TODO: реализуй метод
	// Формат: "Resource с ID X не найден"
	return fmt.Sprintf("%v c ID %d не найден", e.Resource, e.ID)
}

// ValidateAge проверяет корректность возраста
func ValidateAge(age int) error {
	// TODO: реализуй функцию
	// Если age < 0 — вернуть ValidationError
	// Если age > 150 — вернуть ValidationError
	// Иначе nil
	if age < 0 {
		return ValidationError{
			Field:   "age",
			Message: "не может быть отрицательным",
		}
	}
	if age > 150 {
		return ValidationError{
			Field:   "age",
			Message: "слишком большой возраст",
		}
	}

	return nil
}

// ValidateEmail проверяет корректность email
func ValidateEmail(email string) error {
	// TODO: реализуй функцию
	// Если email пустой — вернуть ValidationError
	// Если нет @ — вернуть ValidationError
	// Иначе nil
	if email == "" {
		return ValidationError{
			Field:   "email",
			Message: "напишите свою почту",
		}
	}
	if !strings.Contains(email, "@") {
		return ValidationError{
			Field:   "email",
			Message: "почта введена некоректно",
		}
	}

	return nil
}

// FindUser ищет пользователя по ID
func FindUser(id int) (string, error) {
	// TODO: реализуй функцию
	// Известные пользователи: 1 -> "Алексей", 2 -> "Мария"
	// Если не найден — вернуть NotFoundError
	if id == 1 {
		return "Алексей", nil
	}
	if id == 2 {
		return "Мария", nil
	}
	return "", NotFoundError{Resource: "User", ID: id}
}

func main() {
	fmt.Println("=== Пользовательские ошибки ===")

	// ValidationError
	fmt.Println("\n--- ValidationError ---")
	ages := []int{-5, 25, 200}
	for _, age := range ages {
		if err := ValidateAge(age); err != nil {
			fmt.Printf("Ошибка для age=%d: %v\n", age, err)

			var valErr ValidationError
			if errors.As(err, &valErr) {
				fmt.Printf("  Поле: %s\n", valErr.Field)
			}
		} else {
			fmt.Printf("age=%d валиден\n", age)
		}
	}

	fmt.Println("\nПроверка email:")
	emails := []string{"", "invalid", "valid@example.com"}
	for _, email := range emails {
		if err := ValidateEmail(email); err != nil {
			fmt.Printf("Ошибка для '%s': %v\n", email, err)
		} else {
			fmt.Printf("'%s' валиден\n", email)
		}
	}

	// NotFoundError
	fmt.Println("\n--- NotFoundError ---")
	ids := []int{1, 2, 999}
	for _, id := range ids {
		name, err := FindUser(id)
		if err != nil {
			fmt.Printf("Ошибка для id=%d: %v\n", id, err)

			var notFound NotFoundError
			if errors.As(err, &notFound) {
				fmt.Printf("  Ресурс: %s, ID: %d\n", notFound.Resource, notFound.ID)
			}
		} else {
			fmt.Printf("Найден пользователь id=%d: %s\n", id, name)
		}
	}
}
