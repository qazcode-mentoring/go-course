package main

import (
	"fmt"
	"math"
)

// Circle представляет круг
type Circle struct {
	Radius float64
}

// Area возвращает площадь круга (π × r²)
func (c Circle) Area() float64 {
	// TODO: реализуй метод
	s := (c.Radius * c.Radius) * math.Pi
	return s
}

// Perimeter возвращает периметр (длину окружности) круга (2 × π × r)
func (c Circle) Perimeter() float64 {
	// TODO: реализуй метод
	return 2 * math.Pi * c.Radius
}

// Scale масштабирует круг (умножает радиус на factor)
func (c *Circle) Scale(factor float64) {
	// TODO: реализуй метод
	// Используй указатель для изменения оригинала
	c.Radius *= factor
}

// Rectangle представляет прямоугольник
type Rectangle struct {
	Width, Height float64
}

// Area возвращает площадь прямоугольника (width × height)
func (r Rectangle) Area() float64 {
	// TODO: реализуй метод
	return r.Width * r.Height
}

// Perimeter возвращает периметр прямоугольника (2 × (width + height))
func (r Rectangle) Perimeter() float64 {
	// TODO: реализуй метод
	return 2 * (r.Width + r.Height)
}

// Scale масштабирует прямоугольник (умножает стороны на factor)
func (r *Rectangle) Scale(factor float64) {
	// TODO: реализуй метод
	r.Height *= factor
	r.Width *= factor
}

func main() {
	fmt.Println("=== Circle ===")
	c := Circle{Radius: 5}
	fmt.Printf("Радиус: %.2f\n", c.Radius)
	fmt.Printf("Площадь: %.2f\n", c.Area())
	fmt.Printf("Периметр: %.2f\n", c.Perimeter())

	c.Scale(2)
	fmt.Printf("После Scale(2), радиус: %.2f\n", c.Radius)

	fmt.Println("\n=== Rectangle ===")
	r := Rectangle{Width: 4, Height: 3}
	fmt.Printf("Размеры: %.2f × %.2f\n", r.Width, r.Height)
	fmt.Printf("Площадь: %.2f\n", r.Area())
	fmt.Printf("Периметр: %.2f\n", r.Perimeter())

	r.Scale(1.5)
	fmt.Printf("После Scale(1.5): %.2f × %.2f\n", r.Width, r.Height)
}
