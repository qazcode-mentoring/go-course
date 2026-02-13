package main

import "fmt"

// Address представляет адрес
type Address struct {
	City     string
	Street   string
	Building string
}

// Person представляет человека с адресом
type Person struct {
	Name    string
	Age     int
	Email   string
	Address Address // вложенная структура
}

// PrintFullInfo выводит полную информацию о человеке включая адрес
func PrintFullInfo(p Person) {
	// TODO: реализуй функцию
	// Вывести: имя, возраст, email, город, улицу, дом
	fmt.Printf("Имя: %s\n", p.Name)
	fmt.Printf("Возраст: %d\n", p.Age)
	fmt.Printf("Почта: %s\n", p.Email)
	fmt.Printf("Город: %s\n", p.Address.City)
	fmt.Printf("Улица: %s\n", p.Address.Street)
	fmt.Printf("Дом: %s\n", p.Address.Building)
}

func main() {
	// Создание Person с вложенным Address
	person := Person{
		Name:  "Мария",
		Age:   30,
		Email: "maria@example.com",
		Address: Address{
			City:     "Москва",
			Street:   "Тверская",
			Building: "1",
		},
	}

	fmt.Println("=== Полная информация ===")
	PrintFullInfo(person)

	// Доступ к вложенным полям
	fmt.Println("\n=== Доступ к полям ===")
	fmt.Println("Имя:", person.Name)
	fmt.Println("Город:", person.Address.City)
	fmt.Println("Улица:", person.Address.Street)
	fmt.Println("Дом:", person.Address.Building)
}
