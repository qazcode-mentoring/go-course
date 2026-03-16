package main

import "fmt"

// generateNumbers генерирует числа от start до end (включительно)
// Возвращает канал, из которого можно читать числа
func generateNumbers(start, end int) <-chan int {
	// TODO: реализуй функцию
	// 1. Создай канал ch := make(chan int)
	// 2. Запусти горутину, которая:
	//    - Использует defer close(ch)
	//    - В цикле отправляет числа от start до end в канал
	// 3. Верни канал
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := start; i <= end; i++ {
			ch <- i
		}
	}()
	return ch
}

// generateFibonacci генерирует первые n чисел Фибоначчи
// Возвращает канал, из которого можно читать числа
// F(0)=0, F(1)=1, F(n)=F(n-1)+F(n-2)
func generateFibonacci(n int) <-chan int {
	// TODO: реализуй функцию
	// 1. Создай канал
	// 2. Запусти горутину, которая:
	//    - Закрывает канал при завершении (defer close)
	//    - Генерирует n чисел Фибоначчи
	//    - Отправляет каждое число в канал
	// 3. Верни канал
	ch := make(chan int)
	x, y := 0, 1
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			ch <- x
			x, y = y, x+y
		}
	}()
	return ch
}

// isPrime проверяет, является ли число простым
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// generatePrimes генерирует простые числа до maxNum (включительно)
// Возвращает канал, из которого можно читать простые числа
func generatePrimes(maxNum int) <-chan int {
	// TODO: реализуй функцию
	// 1. Создай канал
	// 2. Запусти горутину, которая:
	//    - Закрывает канал при завершении
	//    - Проверяет каждое число от 2 до maxNum
	//    - Если число простое, отправляет в канал
	// 3. Верни канал
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 2; i <= maxNum; i++ {
			if isPrime(i) {
				ch <- i
			}
		}
	}()
	return ch
}

func main() {
	fmt.Println("=== Channel Generator ===")

	// Тест generateNumbers
	fmt.Println("\n--- Числа от 1 до 10 ---")
	for n := range generateNumbers(1, 10) {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// Тест generateFibonacci
	fmt.Println("\n--- Первые 15 чисел Фибоначчи ---")
	for n := range generateFibonacci(15) {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// Тест generatePrimes
	fmt.Println("\n--- Простые числа до 50 ---")
	for n := range generatePrimes(50) {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// Дополнительный тест: сумма чисел от генератора
	fmt.Println("\n--- Сумма чисел от 1 до 100 ---")
	sum := 0
	for n := range generateNumbers(1, 100) {
		sum += n
	}
	fmt.Printf("Сумма: %d\n", sum) // Ожидается: 5050
}
