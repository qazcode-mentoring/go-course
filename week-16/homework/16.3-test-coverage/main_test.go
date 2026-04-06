package main

import (
	"errors"
	"fmt"
	"testing"
)

// TestValidator_ValidateUsername проверяет валидацию username
func TestValidator_ValidateUsername(t *testing.T) {
	v := NewValidator()

	// TODO: Реализуй группу тестов для валидных username
	t.Run("valid usernames", func(t *testing.T) {
		validNames := []string{
			// TODO: Добавь минимум 3 валидных username
			"user",
			"user_123",
			"User_Name_Long",
		}

		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				if err := v.ValidateUsername(name); err != nil {
					t.Errorf("expected %q to be valid, got error: %v", name, err)
				}
			})
		}
	})

	// TODO: Реализуй группу тестов для слишком коротких username
	t.Run("too short", func(t *testing.T) {
		// TODO: Проверь username длиной 1-2 символа
		err := v.ValidateUsername("ab")
		if !errors.Is(err, ErrUsernameTooShort) {
			t.Errorf("expected ErrUsernameTooShort, got %v", err)
		}
	})

	// TODO: Реализуй группу тестов для слишком длинных username
	t.Run("too long", func(t *testing.T) {
		// TODO: Проверь username длиной > 20 символов
		err := v.ValidateUsername("damirdamirdamirdamirdamirdamir")
		if !errors.Is(err, ErrUsernameTooLong) {
			t.Errorf("expected ErrUsernameTooLong, got %v", err)
		}
	})

	// TODO: Реализуй тест для username, начинающегося с цифры
	t.Run("starts with digit", func(t *testing.T) {
		// TODO: Проверь "1user", "123abc" и т.д.
		err := v.ValidateUsername("1user")
		if !errors.Is(err, ErrUsernameStartsDigit) {
			t.Errorf("expected ErrUsernameStartsDigit, got %v", err)
		}
	})

	// TODO: Реализуй тест для username с недопустимыми символами
	t.Run("invalid characters", func(t *testing.T) {
		// TODO: Проверь "user@name", "user name", "user-name"
		invalidNames := []string{
			"user@name",
			"user name",
			"user-name",
		}

		for _, name := range invalidNames {
			t.Run(name, func(t *testing.T) {
				err := v.ValidateUsername(name)

				if !errors.Is(err, ErrUsernameInvalidChar) {
					t.Errorf("expected ErrUsernameInvalidChar, got %v", err)
				}
			})
		}
	})
}

// TestValidator_ValidateEmail проверяет валидацию email
func TestValidator_ValidateEmail(t *testing.T) {
	v := NewValidator()

	t.Run("valid emails", func(t *testing.T) {
		// TODO: Добавь валидные email адреса
		validEmails := []string{
			"user@example.com",
			"test@test.org",
			"a@b.co",
		}

		for _, email := range validEmails {
			t.Run(email, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				if err := v.ValidateEmail(email); err != nil {
					t.Errorf("expected %q to be valid, got error: %v", email, err)
				}
			})
		}
	})

	t.Run("missing @", func(t *testing.T) {
		// TODO: Проверь email без @
		err := v.ValidateEmail("userexample.com")
		if !errors.Is(err, ErrEmailMissingAt) {
			t.Errorf("expected ErrEmailMissingAt, got %v", err)
		}
	})

	t.Run("multiple @", func(t *testing.T) {
		// TODO: Проверь email с несколькими @
		err := v.ValidateEmail("userexample@@.com")
		if !errors.Is(err, ErrEmailMultipleAt) {
			t.Errorf("expected ErrEmailMultipleAt, got %v", err)
		}
	})

	t.Run("empty before @", func(t *testing.T) {
		// TODO: Проверь "@example.com"
		err := v.ValidateEmail("@example.com")
		if !errors.Is(err, ErrEmailInvalidFormat) {
			t.Errorf("expected ErrEmailInvalidFormat, got %v", err)
		}
	})

	t.Run("invalid after @", func(t *testing.T) {
		// TODO: Проверь "user@ab" (слишком короткий домен)
		// и "user@abcd" (нет точки в домене)
		invalidEmails := []string{
			"user@ab",
			"user@abcd",
		}
		for _, email := range invalidEmails {
			t.Run(email, func(t *testing.T) {
				err := v.ValidateEmail(email)
				if !errors.Is(err, ErrEmailInvalidFormat) {
					t.Errorf("expected ErrEmailInvalidFormat for %q, got %v", email, err)
				}
			})
		}
	})
}

// TestValidator_ValidateAge проверяет валидацию возраста
func TestValidator_ValidateAge(t *testing.T) {
	v := NewValidator()

	t.Run("valid ages", func(t *testing.T) {
		validAges := []int{0, 25, 100, 150}

		for _, age := range validAges {
			// TODO: Используй t.Run с именем из age
			t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
				if err := v.ValidateAge(age); err != nil {
					t.Errorf("expected %d to be valid, got error: %v", age, err)
				}
			})
		}
	})

	t.Run("negative age", func(t *testing.T) {
		// TODO: Проверь отрицательный возраст
		err := v.ValidateAge(-1)
		if !errors.Is(err, ErrAgeNegative) {
			t.Errorf("expected ErrAgeNegative, got %v", err)
		}
	})

	t.Run("too old", func(t *testing.T) {
		// TODO: Проверь возраст > 150
		err := v.ValidateAge(151)
		if !errors.Is(err, ErrAgeTooOld) {
			t.Errorf("expected ErrAgeTooOld, got error: %v", err)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// TODO: Проверь граничные случаи: 0 и 150 (должны быть валидны)
		edge := []int{
			0,
			150,
		}
		for _, age := range edge {
			t.Run(fmt.Sprintf("edge: %d", age), func(t *testing.T) {
				if err := v.ValidateAge(age); err != nil {
					t.Errorf("expected %d to be valid, got error: %v", age, err)
				}
			})
		}
	})
}

// TestValidator_ValidatePassword проверяет валидацию пароля
func TestValidator_ValidatePassword(t *testing.T) {
	v := NewValidator()

	t.Run("valid passwords", func(t *testing.T) {
		validPasswords := []string{
			// TODO: Добавь валидные пароли
			"SecurePass1",
			"MyP4ssword",
			"Test1234Abc",
		}

		for _, pwd := range validPasswords {
			t.Run(pwd, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				if err := v.ValidatePassword(pwd); err != nil {
					t.Errorf("expected %v to be valid, got %v", pwd, err)
				}
			})
		}
	})

	t.Run("too short", func(t *testing.T) {
		// TODO: Проверь пароль < 8 символов
		err := v.ValidatePassword("Short1A")
		if !errors.Is(err, ErrPasswordTooShort) {
			t.Errorf("expected ErrPasswordTooShort, got %v", err)
		}
	})

	t.Run("missing digit", func(t *testing.T) {
		// TODO: Проверь пароль без цифр
		err := v.ValidatePassword("damirilimanov")
		if !errors.Is(err, ErrPasswordNoDigit) {
			t.Errorf("expected ErrPasswordNoDigit, got %v", err)
		}
	})

	t.Run("missing uppercase", func(t *testing.T) {
		// TODO: Проверь пароль без заглавных букв
		err := v.ValidatePassword("damir8ilimanov")
		if !errors.Is(err, ErrPasswordNoUpper) {
			t.Errorf("expected ErrPasswordNoUpper, got %v", err)
		}
	})

	t.Run("missing lowercase", func(t *testing.T) {
		// TODO: Проверь пароль без строчных букв
		err := v.ValidatePassword("DAMIR1234")
		if !errors.Is(err, ErrPasswordNoLower) {
			t.Errorf("expected ErrPasswordNoLower, got %v", err)
		}
	})
}

// TestValidator_ValidateUser проверяет комплексную валидацию
func TestValidator_ValidateUser(t *testing.T) {
	v := NewValidator()

	t.Run("valid user", func(t *testing.T) {
		user := User{
			Username: "john_doe",
			Email:    "john@example.com",
			Age:      25,
			Password: "SecurePass123",
		}

		// TODO: Проверь, что ValidateUser возвращает пустой слайс ошибок
		errs := v.ValidateUser(user)
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %v", errs)
		}
	})

	t.Run("single error", func(t *testing.T) {
		user := User{
			Username: "john_doe",
			Email:    "invalid-email", // невалидный email
			Age:      25,
			Password: "SecurePass123",
		}

		// TODO: Проверь, что возвращается ровно 1 ошибка
		errs := v.ValidateUser(user)
		if len(errs) != 1 {
			t.Errorf("expected 1 error, got %d", len(errs))
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		user := User{
			Username: "ab",      // слишком короткий
			Email:    "invalid", // нет @
			Age:      -5,        // отрицательный
			Password: "weak",    // слишком короткий
		}

		// TODO: Проверь, что возвращаются все 4 ошибки
		errs := v.ValidateUser(user)
		expectedErrors := []error{
			ErrUsernameTooShort,
			ErrEmailMissingAt,
			ErrAgeNegative,
			ErrPasswordTooShort,
		}

		if len(errs) != 4 {
			t.Errorf("expected 4 errors, got %d: %v", len(errs), errs)
		}

		for _, exp := range expectedErrors {
			if !containsError(errs, exp) {
				t.Errorf("expected %v not found in results", exp)
			}
		}
	})
}

// Вспомогательная функция для проверки ошибок (можешь использовать)
func containsError(errs []error, target error) bool {
	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
