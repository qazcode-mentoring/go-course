package main

import (
	"errors"
	"fmt"
	"math"
)

// ErrDivisionByZero возвращается при делении на ноль
var ErrDivisionByZero = errors.New("division by zero")

// ErrNegativeSqrt возвращается при попытке извлечь корень из отрицательного числа
var ErrNegativeSqrt = errors.New("cannot calculate square root of negative number")

// Calculator предоставляет математические операции
type Calculator struct{}

// NewCalculator создаёт новый калькулятор
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add складывает два числа
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract вычитает b из a
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply умножает два числа
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide делит a на b
// Возвращает ошибку ErrDivisionByZero при делении на ноль
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Power возводит base в степень exp
func (c *Calculator) Power(base, exp float64) float64 {
	return math.Pow(base, exp)
}

// Sqrt возвращает квадратный корень числа
// Возвращает ошибку ErrNegativeSqrt для отрицательных чисел
func (c *Calculator) Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(x), nil
}

func main() {
	fmt.Println("=== Calculator для Table-Driven Tests ===")

	calc := NewCalculator()

	// Демонстрация методов
	fmt.Printf("\nAdd(2, 3) = %.2f\n", calc.Add(2, 3))
	fmt.Printf("Subtract(5, 3) = %.2f\n", calc.Subtract(5, 3))
	fmt.Printf("Multiply(4, 5) = %.2f\n", calc.Multiply(4, 5))

	if result, err := calc.Divide(10, 2); err == nil {
		fmt.Printf("Divide(10, 2) = %.2f\n", result)
	}

	if _, err := calc.Divide(10, 0); err != nil {
		fmt.Printf("Divide(10, 0) = error: %v\n", err)
	}

	fmt.Printf("Power(2, 10) = %.2f\n", calc.Power(2, 10))

	if result, err := calc.Sqrt(16); err == nil {
		fmt.Printf("Sqrt(16) = %.2f\n", result)
	}

	if result, err := calc.Sqrt(-1); err != nil {
		fmt.Printf("Sqrt(-1) = error: %v\n", err)
	} else {
		fmt.Printf("Sqrt(-1) = %.2f\n", result)
	}

	fmt.Println("\nНапиши table-driven тесты в main_test.go и запусти: go test -v")
}
