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
			// "user",
			// "user_123",
			// "User_Name_Long",

			"user",
			"user_123",
			"User_Name_Long",
		}

		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				// if err := v.ValidateUsername(name); err != nil {
				//     t.Errorf("expected %q to be valid, got error: %v", name, err)
				// }
				if err := v.ValidateUsername(name); err != nil {
					t.Errorf("expected %q to be valid, got error: %v", name, err)
				}
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

		err := v.ValidateUsername("ab")
		if !errors.Is(err, ErrUsernameTooShort) {
			t.Errorf("expected ErrUsernameTooShort, got %v", err)
		}

		err = v.ValidateUsername("a")
		if !errors.Is(err, ErrUsernameTooShort) {
			t.Errorf("expected ErrUsernameTooShort, got %v", err)
		}

	})

	// TODO: Реализуй группу тестов для слишком длинных username
	t.Run("too long", func(t *testing.T) {
		// TODO: Проверь username длиной > 20 символов
		err := v.ValidateUsername("reallylongusernameeeblablablablagolang")
		if !errors.Is(err, ErrUsernameTooLong) {
			t.Errorf("expected ErrUsernameTooLong, got %v", err)
		}
	})

	// TODO: Реализуй тест для username, начинающегося с цифры
	t.Run("starts with digit", func(t *testing.T) {
		// TODO: Проверь "1user", "123abc" и т.д
		err := v.ValidateUsername("1user")
		if !errors.Is(err, ErrUsernameStartsDigit) {
			t.Errorf("expected ErrUsernameStartsDigit, got %v", err)
		}
	})

	t.Run("starts with digit", func(t *testing.T) {
		err := v.ValidateUsername("123abc")
		if !errors.Is(err, ErrUsernameStartsDigit) {
			t.Errorf("expected ErrUsernameStartsDigit, got %v", err)
		}
	})

	// TODO: Реализуй тест для username с недопустимыми символами
	t.Run("invalid characters", func(t *testing.T) {
		// TODO: Проверь "user@name", "user name", "user-name"
		err := v.ValidateUsername("user@name")
		if !errors.Is(err, ErrUsernameInvalidChar) {
			t.Errorf("expected ErrUsernameInvalidChar, got %v", err)
		}
	})

	t.Run("invalid characters", func(t *testing.T) {
		err := v.ValidateUsername("user name")
		if !errors.Is(err, ErrUsernameInvalidChar) {
			t.Errorf("expected ErrUsernameInvalidChar, got %v", err)
		}
	})

	t.Run("invalid characters", func(t *testing.T) {
		err := v.ValidateUsername("user-name")
		if !errors.Is(err, ErrUsernameInvalidChar) {
			t.Errorf("expected ErrUsernameInvalidChar, got %v", err)
		}
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

			"user@example.com",
			"test@test.org",
			"a@b.co",
			"islam.uzakpai@narxoz.kz",
			"iuzakbayy@gmail.com",
		}

		for _, email := range validEmails {
			t.Run(email, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				if err := v.ValidateEmail(email); err != nil {
					t.Errorf("expected %q to be valid, got %v", email, err)
				}
			})
		}
	})

	t.Run("missing @", func(t *testing.T) {
		// TODO: Проверь email без @
		// err := v.ValidateEmail("userexample.com")
		// if !errors.Is(err, ErrEmailMissingAt) {
		//     t.Errorf("expected ErrEmailMissingAt, got %v", err)
		// }
		err := v.ValidateEmail("userexample.com")
		if !errors.Is(err, ErrEmailMissingAt) {
			t.Errorf("expected ErrEmailMissingAt, got %v", err)
		}
	})

	t.Run("multiple @", func(t *testing.T) {
		// TODO: Проверь email с несколькими @
		err := v.ValidateEmail("user@@example.com")
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
		err := v.ValidateEmail("user@abcd")
		if !errors.Is(err, ErrEmailInvalidFormat) {
			t.Errorf("expected ErrEmailInvalidFormat, got %v", err)
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
			// t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
			//     ...
			// })
			t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
				if err := v.ValidateAge(age); err != nil {
					t.Errorf("expected %q to be valid, got %v", age, err)
				}
			})
		}
	})

	t.Run("negative age", func(t *testing.T) {
		// TODO: Проверь отрицательный возраст
		// err := v.ValidateAge(-1)
		// if !errors.Is(err, ErrAgeNegative) {
		//     t.Errorf("expected ErrAgeNegative, got %v", err)
		// }
		err := v.ValidateAge(-1)
		if !errors.Is(err, ErrAgeNegative) {
			t.Errorf("expected ErrAgeNegative, got %v", err)
		}
	})

	t.Run("too old", func(t *testing.T) {
		// TODO: Проверь возраст > 150
		err := v.ValidateAge(200)
		if !errors.Is(err, ErrAgeTooOld) {
			t.Errorf("expected ErrAgeTooOld, got %v", err)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// TODO: Проверь граничные случаи: 0 и 150 (должны быть валидны)
		if err := v.ValidateAge(0); err != nil {
			t.Errorf("expected %q to be valid, got %v", 0, err)
		}

		if err := v.ValidateAge(150); err != nil {
			t.Errorf("expected %q to be valid, got %v", 150, err)
		}
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
			"SecurePass1",
			"MyP4ssword",
			"Test1234Abc",
			"Test123456",
			"MyPass1234",
		}

		for _, pwd := range validPasswords {
			t.Run(pwd, func(t *testing.T) {
				// TODO: Проверь, что ошибки нет
				if err := v.ValidatePassword(pwd); err != nil {
					t.Errorf("expected %q to be valid, got %v", pwd, err)
				}
			})
		}
	})

	t.Run("too short", func(t *testing.T) {
		// TODO: Проверь пароль < 8 символов
		// err := v.ValidatePassword("Short1A")
		// if !errors.Is(err, ErrPasswordTooShort) {
		//     t.Errorf("expected ErrPasswordTooShort, got %v", err)
		// }

		err := v.ValidatePassword("Short1A")
		if !errors.Is(err, ErrPasswordTooShort) {
			t.Errorf("expected ErrPasswordTooShort, got %v", err)
		}
	})

	t.Run("missing digit", func(t *testing.T) {
		// TODO: Проверь пароль без цифр
		err := v.ValidatePassword("PasswordWithoutDigits")
		if !errors.Is(err, ErrPasswordNoDigit) {
			t.Errorf("expected ErrPasswordNoDigit, got %v", err)
		}
	})

	t.Run("missing uppercase", func(t *testing.T) {
		// TODO: Проверь пароль без заглавных букв
		err := v.ValidatePassword("passwordwithoutupper1")
		if !errors.Is(err, ErrPasswordNoUpper) {
			t.Errorf("expected ErrPasswordNoUpper, got %v", err)
		}
	})

	t.Run("missing lowercase", func(t *testing.T) {
		// TODO: Проверь пароль без строчных букв
		err := v.ValidatePassword("PASSWORDWITHOUTLOWER1")
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
		// errs := v.ValidateUser(user)
		// if len(errs) != 0 {
		//     t.Errorf("expected no errors, got %v", errs)
		// }
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
		// errs := v.ValidateUser(user)
		// if len(errs) != 1 {
		//     t.Errorf("expected 1 error, got %d", len(errs))
		// }
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
		// errs := v.ValidateUser(user)
		// if len(errs) != 4 {
		//     t.Errorf("expected 4 errors, got %d: %v", len(errs), errs)
		// }
		errs := v.ValidateUser(user)
		if len(errs) != 4 {
			t.Errorf("expected 4 errors, got %d: %v", len(errs), errs)
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
