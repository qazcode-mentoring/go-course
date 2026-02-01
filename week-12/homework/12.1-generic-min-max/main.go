package main

import (
	"cmp"
	"fmt"
)

// Min возвращает меньшее из двух значений
func Min[T cmp.Ordered](a, b T) T {
	// TODO: реализуй функцию
	// Если a < b, верни a, иначе верни b
	var zero T
	return zero
}

// Max возвращает большее из двух значений
func Max[T cmp.Ordered](a, b T) T {
	// TODO: реализуй функцию
	// Если a > b, верни a, иначе верни b
	var zero T
	return zero
}

// MinSlice возвращает минимальный элемент слайса
// Возвращает (min, true) если слайс не пустой
// Возвращает (zero, false) если слайс пустой
func MinSlice[T cmp.Ordered](slice []T) (T, bool) {
	// TODO: реализуй функцию
	// 1. Проверь, что слайс не пустой
	// 2. Инициализируй min первым элементом
	// 3. Пройди по всем элементам и найди минимум
	var zero T
	return zero, false
}

// MaxSlice возвращает максимальный элемент слайса
// Возвращает (max, true) если слайс не пустой
// Возвращает (zero, false) если слайс пустой
func MaxSlice[T cmp.Ordered](slice []T) (T, bool) {
	// TODO: реализуй функцию
	// 1. Проверь, что слайс не пустой
	// 2. Инициализируй max первым элементом
	// 3. Пройди по всем элементам и найди максимум
	var zero T
	return zero, false
}

// Clamp ограничивает значение v в диапазоне [min, max]
func Clamp[T cmp.Ordered](v, minVal, maxVal T) T {
	// TODO: реализуй функцию
	// Если v < minVal, верни minVal
	// Если v > maxVal, верни maxVal
	// Иначе верни v
	// Подсказка: можно использовать Min и Max
	return v
}

func main() {
	fmt.Println("=== Generic Min/Max ===")

	// Тест Min/Max с разными типами
	fmt.Println("\n--- Min/Max ---")
	fmt.Printf("Min(3, 5) = %d (ожидается: 3)\n", Min(3, 5))
	fmt.Printf("Max(3, 5) = %d (ожидается: 5)\n", Max(3, 5))
	fmt.Printf("Min(3.14, 2.71) = %.2f (ожидается: 2.71)\n", Min(3.14, 2.71))
	fmt.Printf("Max(3.14, 2.71) = %.2f (ожидается: 3.14)\n", Max(3.14, 2.71))
	fmt.Printf("Min(\"apple\", \"banana\") = %s (ожидается: apple)\n", Min("apple", "banana"))
	fmt.Printf("Max(\"apple\", \"banana\") = %s (ожидается: banana)\n", Max("apple", "banana"))

	// Тест MinSlice/MaxSlice
	fmt.Println("\n--- MinSlice/MaxSlice ---")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("Слайс: %v\n", nums)

	if min, ok := MinSlice(nums); ok {
		fmt.Printf("MinSlice = %d (ожидается: 1)\n", min)
	} else {
		fmt.Println("MinSlice: слайс пустой")
	}

	if max, ok := MaxSlice(nums); ok {
		fmt.Printf("MaxSlice = %d (ожидается: 9)\n", max)
	} else {
		fmt.Println("MaxSlice: слайс пустой")
	}

	// Тест с пустым слайсом
	fmt.Println("\n--- Пустой слайс ---")
	empty := []int{}
	if _, ok := MinSlice(empty); !ok {
		fmt.Println("MinSlice: корректно вернул false для пустого слайса")
	}
	if _, ok := MaxSlice(empty); !ok {
		fmt.Println("MaxSlice: корректно вернул false для пустого слайса")
	}

	// Тест со слайсом строк
	fmt.Println("\n--- Слайс строк ---")
	words := []string{"go", "python", "rust", "java"}
	fmt.Printf("Слайс: %v\n", words)

	if min, ok := MinSlice(words); ok {
		fmt.Printf("MinSlice = %s (ожидается: go)\n", min)
	}
	if max, ok := MaxSlice(words); ok {
		fmt.Printf("MaxSlice = %s (ожидается: rust)\n", max)
	}

	// Тест Clamp
	fmt.Println("\n--- Clamp ---")
	fmt.Printf("Clamp(5, 0, 10) = %d (ожидается: 5)\n", Clamp(5, 0, 10))
	fmt.Printf("Clamp(-5, 0, 10) = %d (ожидается: 0)\n", Clamp(-5, 0, 10))
	fmt.Printf("Clamp(15, 0, 10) = %d (ожидается: 10)\n", Clamp(15, 0, 10))
	fmt.Printf("Clamp(0.5, 0.0, 1.0) = %.1f (ожидается: 0.5)\n", Clamp(0.5, 0.0, 1.0))
	fmt.Printf("Clamp(-0.5, 0.0, 1.0) = %.1f (ожидается: 0.0)\n", Clamp(-0.5, 0.0, 1.0))
	fmt.Printf("Clamp(1.5, 0.0, 1.0) = %.1f (ожидается: 1.0)\n", Clamp(1.5, 0.0, 1.0))
}
