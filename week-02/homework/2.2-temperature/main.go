package main

import "fmt"

// Celsius представляет температуру в градусах Цельсия
type Celsius float64

// Fahrenheit представляет температуру в градусах Фаренгейта
type Fahrenheit float64

// CelsiusToFahrenheit конвертирует Цельсий в Фаренгейт
func CelsiusToFahrenheit(c Celsius) Fahrenheit {
	// TODO: реализуй функцию
	// Формула: F = C × 9/5 + 32
	f := Fahrenheit((c * 9.0 / 5.0) + 32.0)
	return f
}

// FahrenheitToCelsius конвертирует Фаренгейт в Цельсий
func FahrenheitToCelsius(f Fahrenheit) Celsius {
	// TODO: реализуй функцию
	// Формула: C = (F - 32) × 5/9
	c := Celsius((f - 32.0) * (5.0 / 9.0))
	return c
}

func main() {
	// Тест 1: 100°C = 212°F
	c := Celsius(100.0)
	f := CelsiusToFahrenheit(c)
	fmt.Printf("%.1f°C = %.1f°F\n", c, f)

	// Тест 2: 32°F = 0°C
	f2 := Fahrenheit(32.0)
	c2 := FahrenheitToCelsius(f2)
	fmt.Printf("%.1f°F = %.1f°C\n", f2, c2)

	// TODO: Добавь проверку остальных тестовых значений:
	// 0°C = 32°F
	c1 := Celsius(0.0)
	f1 := CelsiusToFahrenheit(c1)
	fmt.Printf("%.1f°C = %.1f°F\n", c1, f1)

	// -40°C = -40°F
	c3 := Celsius(-40.0)
	f3 := CelsiusToFahrenheit(c3)
	fmt.Printf("%.1f°C = %1.f°F\n", c3, f3)

	// 37°C = 98.6°F
	c4 := Celsius(37.0)
	f4 := CelsiusToFahrenheit(c4)
	fmt.Printf("%.1f°C = %.1f°F\n", c4, f4)
}
