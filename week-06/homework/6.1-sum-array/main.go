package main

import "fmt"

// Sum возвращает сумму всех элементов массива
func Sum(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по всем элементам и сложи их
	sum := 0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}
	return sum
}

// Average возвращает среднее арифметическое элементов массива
func Average(arr [5]int) float64 {
	// TODO: реализуй функцию
	// Используй Sum и подели на количество элементов
	sum := Sum(arr)

	return float64(sum) / float64(len(arr))
}

// Max возвращает максимальный элемент массива
func Max(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по массиву и найди максимум
	maxItem := arr[0]
	for i := 0; i < len(arr); i++ {
		if maxItem < arr[i] {
			maxItem = arr[i]
		}
	}
	return maxItem
}

// Min возвращает минимальный элемент массива
func Min(arr [5]int) int {
	// TODO: реализуй функцию
	// Пройди по массиву и найди минимум
	minItem := arr[0]
	for i := 0; i < len(arr); i++ {
		if minItem > arr[i] {
			minItem = arr[i]
		}
	}
	return minItem
}

func main() {
	arr := [5]int{1, 1, 1, 1, 0}

	fmt.Println("=== Операции с массивом ===")
	fmt.Println("Массив:", arr)
	fmt.Println("Сумма:", Sum(arr))
	fmt.Println("Среднее:", Average(arr))
	fmt.Println("Максимум:", Max(arr))
	fmt.Println("Минимум:", Min(arr))
}
