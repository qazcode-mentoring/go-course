package main

import "fmt"

// swap принимает два целых числа и возвращает их в обратном порядке
func swap(a, b int) (int, int) {
	// TODO: реализуй функцию
	return b, a
}

func main() {
	x, y := 10, 20
	fmt.Println("До swap:", x, y)

	x, y = swap(x, y)
	fmt.Println("После swap:", x, y)
}
