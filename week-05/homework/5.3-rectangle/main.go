package main

import "fmt"

// Rectangle представляет прямоугольник
type Rectangle struct {
	Width  float64
	Height float64
}

// Area вычисляет площадь прямоугольника
func Area(r Rectangle) float64 {
	// TODO: реализуй функцию
	return r.Width * r.Height
}

// Perimeter вычисляет периметр прямоугольника
func Perimeter(r Rectangle) float64 {
	// TODO: реализуй функцию
	return (r.Width + r.Height) * 2
}

// Scale масштабирует прямоугольник на заданный коэффициент
func Scale(r *Rectangle, factor float64) {
	// TODO: реализуй функцию
	// Умножь Width и Height на factor
	r.Width = r.Width * factor
	r.Height = r.Height * factor
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Прямоугольник: %+v\n", rect)
	fmt.Printf("Площадь: %.2f\n", Area(rect))
	fmt.Printf("Периметр: %.2f\n", Perimeter(rect))

	fmt.Println("\n=== Масштабирование x2 ===")
	Scale(&rect, 2)
	fmt.Printf("Прямоугольник: %+v\n", rect)
	fmt.Printf("Площадь: %.2f\n", Area(rect))
	fmt.Printf("Периметр: %.2f\n", Perimeter(rect))
}
