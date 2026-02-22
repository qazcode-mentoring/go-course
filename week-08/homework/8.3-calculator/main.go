package main

import (
	"errors"
	"fmt"
)

// Operation — тип для математической операции
type Operation func(a, b int) int

// operations — map с зарегистрированными операциями
var operations map[string]Operation

// init инициализирует базовые операции
func init() {
	// TODO: реализуй инициализацию
	// Создай map и добавь операции: "+", "-", "*", "/"
	// Для деления верни 0 при делении на ноль (или обработай в Calculate)
	operations = make(map[string]Operation)
	operations["+"] = func(a, b int) int { return a + b }
	operations["-"] = func(a, b int) int { return a - b }
	operations["*"] = func(a, b int) int { return a * b }
	operations["/"] = func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	}
}

// Calculate выполняет операцию op над числами a и b
func Calculate(op string, a, b int) (int, error) {
	// TODO: реализуй функцию
	// Найди операцию в map
	// Если не найдена — верни ошибку
	// Для деления проверь b != 0
	v, ok := operations[op]
	if ok {
		return v(a, b), nil
	}
	return 0, errors.New("not implemented")
}

// RegisterOperation добавляет новую операцию в калькулятор
func RegisterOperation(op string, fn Operation) {
	// TODO: реализуй функцию
	operations[op] = fn
}

func main() {
	fmt.Println("=== Калькулятор ===")

	// Базовые операции
	ops := []string{"+", "-", "*", "/"}
	a, b := 10, 5

	for _, op := range ops {
		result, err := Calculate(op, a, b)
		if err != nil {
			fmt.Printf("%d %s %d = ошибка: %v\n", a, op, b, err)
		} else {
			fmt.Printf("%d %s %d = %d\n", a, op, b, result)
		}
	}

	// Деление на ноль
	fmt.Println("\nДеление на ноль:")
	result, err := Calculate("/", 10, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}

	// Добавление новой операции
	fmt.Println("\nДобавляем операцию возведения в степень (^):")
	RegisterOperation("^", func(a, b int) int {
		result := 1
		for i := 0; i < b; i++ {
			result *= a
		}
		return result
	})

	result, _ = Calculate("^", 2, 8)
	fmt.Printf("2 ^ 8 = %d\n", result)

	// Неизвестная операция
	fmt.Println("\nНеизвестная операция:")
	_, err = Calculate("%", 10, 3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
}
