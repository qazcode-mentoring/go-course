package main

import "fmt"

// Person представляет информацию о человеке
type Person struct {
	Name  string
	Age   int
	Email string
}

// NewPerson создаёт и возвращает новую структуру Person
func NewPerson(name string, age int, email string) Person {
	// TODO: реализуй функцию
	return Person{name, age, email}
}

// PrintPerson выводит информацию о человеке в читаемом формате
func PrintPerson(p Person) {
	// TODO: реализуй функцию
	// Формат:
	// Имя: ...
	// Возраст: ...
	// Email: ...
	fmt.Println("Имя: ", p.Name)
	fmt.Println("Возраст: ", p.Age)
	fmt.Println("Email: ", p.Email)
}

// Birthday увеличивает возраст на 1
func Birthday(p *Person) {
	// TODO: реализуй функцию
	// Используй указатель для изменения оригинала
	p.Age += 1
}

func main() {
	// Создаём человека
	person := NewPerson("Алексей", 25, "alex@example.com")

	fmt.Println("=== Информация о человеке ===")
	PrintPerson(person)

	// День рождения
	fmt.Println("\n=== После дня рождения ===")
	Birthday(&person)
	PrintPerson(person)
}
