package main

import (
	"fmt"
	"strings"
)

// StringProcessor — тип для функции обработки строки
type StringProcessor func(string) string

// ToUpper возвращает процессор, переводящий строку в верхний регистр
func ToUpper() StringProcessor {
	// TODO: реализуй функцию
	// Верни функцию, использующую strings.ToUpper
	return func(s string) string {
		return strings.ToUpper(s)
	}
}

// ToLower возвращает процессор, переводящий строку в нижний регистр
func ToLower() StringProcessor {
	// TODO: реализуй функцию
	return func(s string) string {
		return strings.ToLower(s)
	}
}

// Trim возвращает процессор, удаляющий пробелы по краям
func Trim() StringProcessor {
	// TODO: реализуй функцию
	// Используй strings.TrimSpace
	return func(s string) string {
		return strings.TrimSpace(s)
	}
}

// AddPrefix возвращает процессор, добавляющий префикс к строке
func AddPrefix(prefix string) StringProcessor {
	// TODO: реализуй функцию
	// Это замыкание — функция запоминает prefix
	return func(s string) string {
		return prefix + s
	}
}

// Compose объединяет несколько процессоров в один
// Процессоры применяются слева направо
func Compose(processors ...StringProcessor) StringProcessor {
	// TODO: реализуй функцию
	// Верни функцию, которая применяет все процессоры по очереди
	return func(s string) string {
		for _, p := range processors {
			s = p(s)
		}
		return s
	}
}

// Pipeline применяет процессоры к строке последовательно
func Pipeline(input string, processors ...StringProcessor) string {
	// TODO: реализуй функцию
	// Пройди по процессорам и применяй каждый к результату
	for _, p := range processors {
		input = p(input)
	}
	return input
}

func main() {
	fmt.Println("=== Композиция функций ===")

	// Отдельные процессоры
	fmt.Println("ToUpper(\"hello\"):", ToUpper()("hello"))
	fmt.Println("ToLower(\"HELLO\"):", ToLower()("HELLO"))
	fmt.Println("Trim(\"  hello  \"):", Trim()("  hello  "))
	fmt.Println("AddPrefix(\">>> \")(\"hello\"):", AddPrefix(">>> ")("hello"))

	// Compose
	fmt.Println("\n=== Compose ===")
	processor := Compose(Trim(), ToLower(), AddPrefix(">>> "))
	result := processor("  HELLO WORLD  ")
	fmt.Println("Compose(Trim, ToLower, AddPrefix):", result)

	// Pipeline
	fmt.Println("\n=== Pipeline ===")
	result = Pipeline("  HELLO WORLD  ",
		Trim(),
		ToLower(),
		AddPrefix("Result: "),
	)
	fmt.Println("Pipeline result:", result)
}
