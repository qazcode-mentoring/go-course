package main

import "fmt"

// Greet возвращает строку приветствия для указанного имени
func Greet(name string) string {
	// TODO: реализуй функцию
	// Формат: "Привет, {имя}! Добро пожаловать в мир Go!"
	// Бонус: если имя пустое, вернуть "Привет, незнакомец!"
	if name == "" {
		return "Привет, незнакомец!"
	}

	return fmt.Sprintf("Привет, %s! Добро пожаловать в мир Go!", name)
}

func main() {
	message := Greet("Алексей")
	fmt.Println(message)

	fmt.Println(Greet("Мария"))

	// Бонус: проверка пустого имени
	fmt.Println(Greet(""))
}
