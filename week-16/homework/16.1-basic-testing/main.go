package main

import (
	"fmt"
	"strings"
)

// Reverse переворачивает строку
// Пример: "hello" -> "olleh"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome проверяет, является ли строка палиндромом
// Палиндром читается одинаково в обе стороны
// Пустая строка и один символ считаются палиндромами
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// CountWords считает количество слов в строке
// Слова разделяются пробелами
// Множественные пробелы не создают "пустых" слов
func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// Truncate обрезает строку до maxLen символов
// Если строка была обрезана, добавляет "..." в конец
// Если maxLen <= 3, возвращает только "..."
// Если строка короче или равна maxLen, возвращает как есть
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if maxLen <= 3 {
		return "..."[:maxLen]
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen-3]) + "..."
}

func main() {
	fmt.Println("=== Функции для тестирования ===")

	// Демонстрация Reverse
	fmt.Println("\n--- Reverse ---")
	fmt.Printf("Reverse(%q) = %q\n", "hello", Reverse("hello"))
	fmt.Printf("Reverse(%q) = %q\n", "Go", Reverse("Go"))
	fmt.Printf("Reverse(%q) = %q\n", "", Reverse(""))

	// Демонстрация IsPalindrome
	fmt.Println("\n--- IsPalindrome ---")
	fmt.Printf("IsPalindrome(%q) = %v\n", "radar", IsPalindrome("radar"))
	fmt.Printf("IsPalindrome(%q) = %v\n", "hello", IsPalindrome("hello"))
	fmt.Printf("IsPalindrome(%q) = %v\n", "A", IsPalindrome("A"))

	// Демонстрация CountWords
	fmt.Println("\n--- CountWords ---")
	fmt.Printf("CountWords(%q) = %d\n", "hello world", CountWords("hello world"))
	fmt.Printf("CountWords(%q) = %d\n", "  multiple   spaces  ", CountWords("  multiple   spaces  "))
	fmt.Printf("CountWords(%q) = %d\n", "", CountWords(""))

	// Демонстрация Truncate
	fmt.Println("\n--- Truncate ---")
	fmt.Printf("Truncate(%q, 10) = %q\n", "Hello, World!", Truncate("Hello, World!", 10))
	fmt.Printf("Truncate(%q, 20) = %q\n", "Short", Truncate("Short", 20))
	fmt.Printf("Truncate(%q, 3) = %q\n", "Hello", Truncate("Hello", 3))

	fmt.Println("\nТеперь напиши тесты в main_test.go и запусти: go test -v")
}
