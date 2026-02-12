package main

import "fmt"

// Countdown выводит обратный отсчёт от n до 1, затем "Поехали!"
func Countdown(n int) {
	// TODO: реализуй функцию используя defer
	defer fmt.Println("Поехали!")
	for i := 1; i <= n; i++ {
		defer fmt.Println(i)
	}
	// первый defer "Поехали" кладется в отложенный стек и defer fmt.Println(i) за ним, затем все выводится с конца.
	// Стек работает по принципу LIFO
}

func main() {
	fmt.Println("=== Обратный отсчёт ===")
	Countdown(5)
}
