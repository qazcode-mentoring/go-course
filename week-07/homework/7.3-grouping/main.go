package main

import "fmt"

// Student представляет информацию о студенте
type Student struct {
	Name  string
	Grade int
}

// GroupByGrade группирует студентов по классу
func GroupByGrade(students []Student) map[int][]Student {
	// TODO: реализуй функцию
	// Создай map[int][]Student
	stds := make(map[int][]Student)
	// Пройди по студентам и добавляй в соответствующий слайс
	for _, v := range students {
		stds[v.Grade] = append(stds[v.Grade], v)
	}
	return stds
}

// GroupByFirstLetter группирует слова по первой букве
func GroupByFirstLetter(words []string) map[rune][]string {
	// TODO: реализуй функцию
	// Для получения первой буквы используй []rune(word)[0]
	firstLetter := make(map[rune][]string)

	for _, word := range words {
		first := []rune(word)[0]
		firstLetter[first] = append(firstLetter[first], word)
	}
	return firstLetter
}

// GetStudentNames возвращает имена студентов указанного класса
func GetStudentNames(groups map[int][]Student, grade int) []string {
	// TODO: реализуй функцию
	// Получи слайс студентов по grade
	// Извлеки только имена
	names := make([]string, 0)
	for _, v := range groups[grade] {
		names = append(names, v.Name)
	}

	return names
}

func main() {
	students := []Student{
		{"Алексей", 10},
		{"Мария", 10},
		{"Иван", 11},
		{"Анна", 11},
		{"Пётр", 10},
	}

	fmt.Println("=== Группировка студентов ===")
	fmt.Println("Все студенты:", students)

	groups := GroupByGrade(students)
	fmt.Println("\nПо классам:")
	for grade, list := range groups {
		fmt.Printf("  %d класс: %v\n", grade, list)
	}

	fmt.Println("\nИмена 10 класса:", GetStudentNames(groups, 10))

	fmt.Println("\n=== Группировка слов ===")
	words := []string{"apple", "apricot", "banana", "avocado", "blueberry"}
	byLetter := GroupByFirstLetter(words)
	fmt.Println("Слова:", words)
	fmt.Println("По первой букве:", byLetter)
}
