package main

import "fmt"

// IsPrime проверяет, является ли число простым
func IsPrime(n int) bool {
	// TODO: реализуй функцию
	if n <= 1 {
		return false
	}
	var d int = 2
	for n%d != 0 {
		d += 1
	}
	return d == n
}

func main() {
	// Тесты
	testCases := []int{-5, 0, 1, 2, 3, 4, 5, 7, 11, 13, 15, 17, 18, 19, 20, 23, 100}

	for _, n := range testCases {
		result := IsPrime(n)
		fmt.Printf("IsPrime(%d) = %v\n", n, result)
	}
}
