package main

import "fmt"

// AppendUnique добавляет элемент в слайс только если его там ещё нет
func AppendUnique(slice []int, value int) []int {
	// TODO: реализуй функцию
	// Проверь, есть ли value в slice
	// Если нет — добавь через append
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return slice
		}
	}
	return append(slice, value)
}

// RemoveAt удаляет элемент по индексу
func RemoveAt(slice []int, index int) []int {
	// TODO: реализуй функцию
	// Проверь, что index в допустимых границах
	// Используй append для соединения частей до и после index
	if index < 0 || index >= len(slice) {
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}

// RemoveValue удаляет первое вхождение значения из слайса
func RemoveValue(slice []int, value int) []int {
	// TODO: реализуй функцию
	// Найди индекс value и используй RemoveAt
	for i, v := range slice {
		if v == value {
			return RemoveAt(slice, i)
		}
	}
	return slice
}

// SumAll суммирует любое количество чисел (variadic функция)
func SumAll(nums ...int) int {
	// TODO: реализуй функцию
	// nums — это слайс, пройди по нему и сложи
	var sum int = 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("=== Операции со слайсами ===")
	fmt.Println("Исходный слайс:", nums)

	// Добавление уникальных
	nums = AppendUnique(nums, 6)
	fmt.Println("После AppendUnique(6):", nums)
	nums = AppendUnique(nums, 3)
	fmt.Println("После AppendUnique(3):", nums)

	// Удаление по индексу
	nums = RemoveAt(nums, 2)
	fmt.Println("После RemoveAt(2):", nums)

	// Удаление по значению
	nums = RemoveValue(nums, 5)
	fmt.Println("После RemoveValue(5):", nums)

	// Variadic функция
	fmt.Println("\n=== Variadic функция ===")
	fmt.Println("SumAll(1, 2, 3):", SumAll(1, 2, 3))
	fmt.Println("SumAll(nums...):", SumAll(nums...))
}
