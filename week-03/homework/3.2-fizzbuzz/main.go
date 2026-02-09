package main

import "fmt"

// FizzBuzz выводит числа от 1 до n по правилам FizzBuzz
func FizzBuzz(n int) {
	// TODO: реализуй функцию
	// Правила:
	// - Делится на 3 и 5 → "FizzBuzz"
	// - Делится на 3 → "Fizz"
	// - Делится на 5 → "Buzz"
	// - Иначе → само число
	for i := 1; i <= n; i++ {
		// TODO: добавь логику
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}

func main() {
	FizzBuzz(15)
}
