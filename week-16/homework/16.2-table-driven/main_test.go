package main

import (
	"testing"
)

// almostEqual проверяет равенство float64 с погрешностью
func almostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

// TestCalculator_Add проверяет метод Add с использованием table-driven tests
func TestCalculator_Add(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		// TODO: Добавь минимум 5 тестовых случаев
		// Пример:
		// {"positive numbers", 2, 3, 5},
		// {"negative numbers", -2, -3, -5},
		// {"mixed signs", -2, 5, 3},
		// {"with zero", 5, 0, 5},
		// {"decimals", 1.5, 2.5, 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку
			// result := calc.Add(tt.a, tt.b)
			// if result != tt.expected {
			//     t.Errorf("Add(%.2f, %.2f) = %.2f; want %.2f",
			//         tt.a, tt.b, result, tt.expected)
			// }
			_ = calc // убери эту строку после реализации
		})
	}
}

// TestCalculator_Subtract проверяет метод Subtract
func TestCalculator_Subtract(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		// TODO: Добавь минимум 4 тестовых случая
		// Пример:
		// {"positive result", 5, 3, 2},
		// {"negative result", 3, 5, -2},
		// {"subtract zero", 5, 0, 5},
		// {"both negative", -2, -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку
			_ = calc
		})
	}
}

// TestCalculator_Multiply проверяет метод Multiply
func TestCalculator_Multiply(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		// TODO: Добавь минимум 5 тестовых случаев
		// Не забудь проверить:
		// - умножение на 0
		// - умножение на 1
		// - отрицательные числа
		// - дробные числа
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку
			_ = calc
		})
	}
}

// TestCalculator_Divide проверяет метод Divide с обработкой ошибок
func TestCalculator_Divide(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		expectErr bool
	}{
		// TODO: Добавь минимум 5 тестовых случаев
		// Пример:
		// {"normal division", 10, 2, 5, false},
		// {"division by zero", 10, 0, 0, true},
		// {"divide zero", 0, 5, 0, false},
		// {"negative divisor", 10, -2, -5, false},
		// {"both negative", -10, -2, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку с обработкой ошибок
			// result, err := calc.Divide(tt.a, tt.b)
			//
			// if tt.expectErr {
			//     if err == nil {
			//         t.Error("expected error, got nil")
			//     }
			//     return
			// }
			//
			// if err != nil {
			//     t.Fatalf("unexpected error: %v", err)
			// }
			//
			// if result != tt.expected {
			//     t.Errorf("Divide(%.2f, %.2f) = %.2f; want %.2f",
			//         tt.a, tt.b, result, tt.expected)
			// }
			_ = calc
		})
	}
}

// TestCalculator_Power проверяет метод Power
func TestCalculator_Power(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		base     float64
		exp      float64
		expected float64
	}{
		// TODO: Добавь минимум 4 тестовых случая
		// Не забудь проверить:
		// - степень 0 (любое число в степени 0 = 1)
		// - степень 1
		// - отрицательную степень (2^-1 = 0.5)
		// - дробную степень (4^0.5 = 2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку
			// Используй almostEqual для сравнения float
			// result := calc.Power(tt.base, tt.exp)
			// if !almostEqual(result, tt.expected, 0.0001) {
			//     t.Errorf("Power(%.2f, %.2f) = %.4f; want %.4f",
			//         tt.base, tt.exp, result, tt.expected)
			// }
			_ = calc
		})
	}
}

// TestCalculator_Sqrt проверяет метод Sqrt с обработкой ошибок
func TestCalculator_Sqrt(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name      string
		x         float64
		expected  float64
		expectErr bool
	}{
		// TODO: Добавь минимум 4 тестовых случая
		// Пример:
		// {"perfect square", 16, 4, false},
		// {"zero", 0, 0, false},
		// {"non-perfect square", 2, 1.4142, false}, // используй almostEqual
		// {"negative number", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Реализуй проверку с обработкой ошибок
			// Используй almostEqual для сравнения результата
			_ = calc
		})
	}
}
