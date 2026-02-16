package main

import "fmt"

// template
// Sum возвращает сумму всех элементов массива
func Sum(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по всем элементам и сложи их
	return 0
}

// Average возвращает среднее арифметическое элементов массива
func Average(arr [5]int) float64 {
	// TODO: реализуй функцию
	// Используй Sum и подели на количество элементов
	return 0.0
}

// Max возвращает максимальный элемент массива
func Max(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по массиву и найди максимум
	return 0
}

// Min возвращает минимальный элемент массива
func Min(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по массиву и найди минимум
	return 0
}

func main() {
	arr := [5]int{10, 20, 5, 15, 30}

	fmt.Println("=== Операции с массивом ===")
	fmt.Println("Массив:", arr)
	fmt.Println("Сумма:", Sum(arr))
	fmt.Println("Среднее:", Average(arr))
	fmt.Println("Максимум:", Max(arr))
	fmt.Println("Минимум:", Min(arr))
}
