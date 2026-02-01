package main

import (
	"fmt"
)

// TypeName возвращает название типа переданного значения
func TypeName(v any) string {
	// TODO: реализуй функцию
	// Используй type switch:
	// switch v.(type) {
	// case int: return "int"
	// case string: return "string"
	// ...
	// }
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case float64:
		return "float64"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

// ToString преобразует любое значение в строку
func ToString(v any) string {
	// TODO: реализуй функцию
	// Можно использовать fmt.Sprintf("%v", v)
	// Или обработать типы отдельно через type switch

	return fmt.Sprintf("%v", v)
}

// Sum суммирует все числовые значения из слайса
// Игнорирует нечисловые значения
func Sum(values []any) float64 {
	// TODO: реализуй функцию
	// Используй type switch для int и float64
	var sum float64
	for _, v := range values {
		switch value := v.(type) {
		case int:
			sum += float64(value)
		case float64:
			sum += value
		default:
			continue
		}
	}
	return sum
}

// FilterByType возвращает слайс значений указанного типа
// Используй type assertion для проверки типа
func FilterByType[T any](values []any) []T {
	// TODO: реализуй функцию
	// Для каждого значения попробуй v.(T)
	// Если ok == true, добавь в результат
	var result []T
	for _, v := range values {
		if typed, ok := v.(T); ok {
			result = append(result, typed)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Пустой интерфейс (any) ===")

	// TypeName
	fmt.Println("\n--- TypeName ---")
	fmt.Println("TypeName(42):", TypeName(42))
	fmt.Println("TypeName(\"hello\"):", TypeName("hello"))
	fmt.Println("TypeName(3.14):", TypeName(3.14))
	fmt.Println("TypeName(true):", TypeName(true))

	// ToString
	fmt.Println("\n--- ToString ---")
	fmt.Println("ToString(42):", ToString(42))
	fmt.Println("ToString(true):", ToString(true))
	fmt.Println("ToString(3.14):", ToString(3.14))
	fmt.Println("ToString([]int{1,2,3}):", ToString([]int{1, 2, 3}))

	// Sum
	fmt.Println("\n--- Sum ---")
	values := []any{1, 2.5, "skip", 3, 4.5, true}
	fmt.Println("Values:", values)
	fmt.Println("Sum:", Sum(values))

	// FilterByType
	fmt.Println("\n--- FilterByType ---")
	mixed := []any{1, "hello", 2, "world", 3, true}
	fmt.Println("Mixed:", mixed)
	fmt.Println("Strings:", FilterByType[string](mixed))
	fmt.Println("Ints:", FilterByType[int](mixed))
}
