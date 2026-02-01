package main

import (
	"errors"
	"testing"
)

// TestValidator_ValidateUsername проверяет валидацию username
func TestValidator_ValidateUsername(t *testing.T) {
	v := NewValidator()

	// TODO: Реализуй группу тестов для валидных username
	t.Run("valid usernames", func(t *testing.T) {
		validNames := []string{
			// TODO: Добавь минимум 3 валидных username
			// "user",
			// "user_123",
			// "User_Name_Long",
		}

		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				// if err := v.ValidateUsername(name); err != nil {
				//     t.Errorf("expected %q to be valid, got error: %v", name, err)
				// }
				_ = v
			})
		}
	})

	// TODO: Реализуй группу тестов для слишком коротких username
	t.Run("too short", func(t *testing.T) {
		// TODO: Проверь username длиной 1-2 символа
		// err := v.ValidateUsername("ab")
		// if !errors.Is(err, ErrUsernameTooShort) {
		//     t.Errorf("expected ErrUsernameTooShort, got %v", err)
		// }
		_ = v
	})

	// TODO: Реализуй группу тестов для слишком длинных username
	t.Run("too long", func(t *testing.T) {
		// TODO: Проверь username длиной > 20 символов
		_ = v
	})

	// TODO: Реализуй тест для username, начинающегося с цифры
	t.Run("starts with digit", func(t *testing.T) {
		// TODO: Проверь "1user", "123abc" и т.д.
		_ = v
	})

	// TODO: Реализуй тест для username с недопустимыми символами
	t.Run("invalid characters", func(t *testing.T) {
		// TODO: Проверь "user@name", "user name", "user-name"
		_ = v
	})
}

// TestValidator_ValidateEmail проверяет валидацию email
func TestValidator_ValidateEmail(t *testing.T) {
	v := NewValidator()

	t.Run("valid emails", func(t *testing.T) {
		// TODO: Добавь валидные email адреса
		validEmails := []string{
			// "user@example.com",
			// "test@test.org",
			// "a@b.co",
		}

		for _, email := range validEmails {
			t.Run(email, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				_ = v
			})
		}
	})

	t.Run("missing @", func(t *testing.T) {
		// TODO: Проверь email без @
		// err := v.ValidateEmail("userexample.com")
		// if !errors.Is(err, ErrEmailMissingAt) {
		//     t.Errorf("expected ErrEmailMissingAt, got %v", err)
		// }
		_ = v
	})

	t.Run("multiple @", func(t *testing.T) {
		// TODO: Проверь email с несколькими @
		_ = v
	})

	t.Run("empty before @", func(t *testing.T) {
		// TODO: Проверь "@example.com"
		_ = v
	})

	t.Run("invalid after @", func(t *testing.T) {
		// TODO: Проверь "user@ab" (слишком короткий домен)
		// и "user@abcd" (нет точки в домене)
		_ = v
	})
}

// TestValidator_ValidateAge проверяет валидацию возраста
func TestValidator_ValidateAge(t *testing.T) {
	v := NewValidator()

	t.Run("valid ages", func(t *testing.T) {
		validAges := []int{0, 25, 100, 150}

		for _, age := range validAges {
			// TODO: Используй t.Run с именем из age
			// t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
			//     ...
			// })
			_ = age
		}
		_ = v
	})

	t.Run("negative age", func(t *testing.T) {
		// TODO: Проверь отрицательный возраст
		// err := v.ValidateAge(-1)
		// if !errors.Is(err, ErrAgeNegative) {
		//     t.Errorf("expected ErrAgeNegative, got %v", err)
		// }
		_ = v
	})

	t.Run("too old", func(t *testing.T) {
		// TODO: Проверь возраст > 150
		_ = v
	})

	t.Run("edge cases", func(t *testing.T) {
		// TODO: Проверь граничные случаи: 0 и 150 (должны быть валидны)
		_ = v
	})
}

// TestValidator_ValidatePassword проверяет валидацию пароля
func TestValidator_ValidatePassword(t *testing.T) {
	v := NewValidator()

	t.Run("valid passwords", func(t *testing.T) {
		validPasswords := []string{
			// TODO: Добавь валидные пароли
			// "SecurePass1",
			// "MyP4ssword",
			// "Test1234Abc",
		}

		for _, pwd := range validPasswords {
			t.Run(pwd, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				_ = v
			})
		}
	})

	t.Run("too short", func(t *testing.T) {
		// TODO: Проверь пароль < 8 символов
		// err := v.ValidatePassword("Short1A")
		// if !errors.Is(err, ErrPasswordTooShort) {
		//     t.Errorf("expected ErrPasswordTooShort, got %v", err)
		// }
		_ = v
	})

	t.Run("missing digit", func(t *testing.T) {
		// TODO: Проверь пароль без цифр
		_ = v
	})

	t.Run("missing uppercase", func(t *testing.T) {
		// TODO: Проверь пароль без заглавных букв
		_ = v
	})

	t.Run("missing lowercase", func(t *testing.T) {
		// TODO: Проверь пароль без строчных букв
		_ = v
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
		// errs := v.ValidateUser(user)
		// if len(errs) != 0 {
		//     t.Errorf("expected no errors, got %v", errs)
		// }
		_ = user
		_ = v
	})

	t.Run("single error", func(t *testing.T) {
		user := User{
			Username: "john_doe",
			Email:    "invalid-email", // невалидный email
			Age:      25,
			Password: "SecurePass123",
		}

		// TODO: Проверь, что возвращается ровно 1 ошибка
		// errs := v.ValidateUser(user)
		// if len(errs) != 1 {
		//     t.Errorf("expected 1 error, got %d", len(errs))
		// }
		_ = user
		_ = v
	})

	t.Run("multiple errors", func(t *testing.T) {
		user := User{
			Username: "ab",      // слишком короткий
			Email:    "invalid", // нет @
			Age:      -5,        // отрицательный
			Password: "weak",    // слишком короткий
		}

		// TODO: Проверь, что возвращаются все 4 ошибки
		// errs := v.ValidateUser(user)
		// if len(errs) != 4 {
		//     t.Errorf("expected 4 errors, got %d: %v", len(errs), errs)
		// }
		_ = user
		_ = v
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
