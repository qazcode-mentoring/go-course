package main

import (
	"fmt"
	"sort"
)

// Person представляет информацию о человеке
type Person struct {
	Name string
	Age  int
}

// ByAge реализует sort.Interface для сортировки по возрасту
type ByAge []Person

// Len возвращает количество элементов
func (a ByAge) Len() int {
	// TODO: реализуй метод
	return len(a)
}

// Swap меняет элементы местами
func (a ByAge) Swap(i, j int) {
	// TODO: реализуй метод
	a[i], a[j] = a[j], a[i]
}

// Less возвращает true если элемент i должен быть перед j (по возрасту)
func (a ByAge) Less(i, j int) bool {
	// TODO: реализуй метод
	return a[i].Age < a[j].Age
}

// ByName реализует sort.Interface для сортировки по имени
type ByName []Person

// Len возвращает количество элементов
func (a ByName) Len() int {
	// TODO: реализуй метод
	return len(a)
}

// Swap меняет элементы местами
func (a ByName) Swap(i, j int) {
	// TODO: реализуй метод
	a[i], a[j] = a[j], a[i]
}

// Less возвращает true если элемент i должен быть перед j (по имени)
func (a ByName) Less(i, j int) bool {
	// TODO: реализуй метод
	// Для строк можно использовать < для лексикографического сравнения
	return a[i].Name < a[j].Name
}

func main() {
	people := []Person{
		{"Мария", 30},
		{"Алексей", 25},
		{"Яна", 22},
		{"Иван", 20},
		{"Борис", 35},
	}

	fmt.Println("=== Реализация sort.Interface ===")
	fmt.Println("Исходный список:")
	for _, p := range people {
		fmt.Printf("  %s, %d лет\n", p.Name, p.Age)
	}

	// Сортировка по возрасту
	fmt.Println("\nСортировка по возрасту (ByAge):")
	sort.Sort(ByAge(people))
	for _, p := range people {
		fmt.Printf("  %s, %d лет\n", p.Name, p.Age)
	}

	// Сортировка по имени
	fmt.Println("\nСортировка по имени (ByName):")
	sort.Sort(ByName(people))
	for _, p := range people {
		fmt.Printf("  %s, %d лет\n", p.Name, p.Age)
	}

	// Обратная сортировка
	fmt.Println("\nОбратная сортировка по возрасту:")
	sort.Sort(sort.Reverse(ByAge(people)))
	for _, p := range people {
		fmt.Printf("  %s, %d лет\n", p.Name, p.Age)
	}
}
