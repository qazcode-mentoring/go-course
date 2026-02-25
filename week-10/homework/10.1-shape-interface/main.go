package main

import (
	"fmt"
	"math"
)

// Shape — интерфейс для геометрических фигур
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle представляет круг
type Circle struct {
	Radius float64
}

// Area возвращает площадь круга
func (c Circle) Area() float64 {
	// TODO: реализуй метод
	s := math.Pi * (c.Radius * c.Radius) // подсказка

	return s
}

// Perimeter возвращает периметр круга
func (c Circle) Perimeter() float64 {
	// TODO: реализуй метод
	return 2 * math.Pi * c.Radius
}

// Rectangle представляет прямоугольник
type Rectangle struct {
	Width, Height float64
}

// Area возвращает площадь прямоугольника
func (r Rectangle) Area() float64 {
	// TODO: реализуй метод
	return r.Width * r.Height
}

// Perimeter возвращает периметр прямоугольника
func (r Rectangle) Perimeter() float64 {
	// TODO: реализуй метод
	return 2 * (r.Width + r.Height)
}

// PrintShapeInfo выводит информацию о любой фигуре
func PrintShapeInfo(s Shape) {
	// TODO: реализуй функцию
	// Выведи тип фигуры, площадь и периметр
	// Подсказка: для типа используй fmt.Printf("%T", s)
	fmt.Printf("тип фигуры: %T, площадь: %.2f, периметр: %.2f", s, s.Area(), s.Perimeter())
}

// TotalArea возвращает сумму площадей всех фигур
func TotalArea(shapes []Shape) float64 {
	// TODO: реализуй функцию
	// Пройди по всем фигурам и сложи их площади
	var total float64
	for i := 0; i < len(shapes); i++ {
		total += shapes[i].Area()
	}
	return total
}

func main() {
	fmt.Println("=== Интерфейс Shape ===")

	// Создаём разные фигуры
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 3},
		Circle{Radius: 2},
		Rectangle{Width: 10, Height: 5},
	}

	fmt.Printf("shapes[1]: %v\n", shapes[1])

	// Выводим информацию о каждой
	fmt.Println("Информация о фигурах:")
	for _, s := range shapes {
		PrintShapeInfo(s)
	}

	// Общая площадь
	fmt.Printf("\nОбщая площадь всех фигур: %.2f\n", TotalArea(shapes))
}
