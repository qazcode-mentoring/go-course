package main

import (
	"errors"
	"fmt"
)

// Factorial вычисляет факториал числа n
// Возвращает ошибку для отрицательных чисел
func Factorial(n int) (int, error) {
	// TODO: реализуй функцию
	// 1. Проверь на отрицательное число → вернуть ошибку
	// 2. 0! = 1
	// 3. n! = 1 * 2 * ... * n

	if n < 0 {
		return 0, errors.New("факториал отрицательного числа не определён")
	} else if n == 0 {
		return 1, nil
	}
	// TODO: вычисли факториал
	var result = 1
	for i := 1; i < n; i++ {
		result *= i + 1
	}

	return result, nil
}

func main() {
	// Тест 1: 5! = 120
	result, err := Factorial(5)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("5! =", result)
	}

	// Тест 2: 0! = 1
	result, err = Factorial(0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("0! =", result)
	}

	// Тест 3: отрицательное число
	result, err = Factorial(-3)
	if err != nil {
		fmt.Println("Factorial(-3):", err)
	} else {
		fmt.Println("-3! =", result)
	}

	// TODO: добавь тесты для 1!, 10!
}
