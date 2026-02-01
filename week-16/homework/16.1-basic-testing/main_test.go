package main

import "testing"

// TestReverse проверяет функцию Reverse
func TestReverse(t *testing.T) {
	// TODO: Проверь переворот обычной строки
	// result := Reverse("hello")
	// if result != "olleh" {
	//     t.Errorf("Reverse(%q) = %q; want %q", "hello", result, "olleh")
	// }

	// TODO: Проверь пустую строку
	// Пустая строка при перевороте остаётся пустой

	// TODO: Проверь строку из одного символа
	// Один символ при перевороте остаётся тем же

	// TODO: Проверь строку с пробелами
	// "a b" -> "b a"
}

// TestIsPalindrome проверяет функцию IsPalindrome
func TestIsPalindrome(t *testing.T) {
	// TODO: Проверь палиндром "radar"
	// if !IsPalindrome("radar") {
	//     t.Error("IsPalindrome(\"radar\") should be true")
	// }

	// TODO: Проверь палиндром "level"

	// TODO: Проверь не-палиндром "hello"
	// if IsPalindrome("hello") {
	//     t.Error("IsPalindrome(\"hello\") should be false")
	// }

	// TODO: Проверь пустую строку (должна быть палиндромом)

	// TODO: Проверь строку из одного символа
}

// TestCountWords проверяет функцию CountWords
func TestCountWords(t *testing.T) {
	// TODO: Проверь обычное предложение
	// result := CountWords("hello world")
	// if result != 2 {
	//     t.Errorf("CountWords(%q) = %d; want %d", "hello world", result, 2)
	// }

	// TODO: Проверь пустую строку (должно быть 0 слов)

	// TODO: Проверь строку с множественными пробелами
	// "  multiple   spaces  " должно дать 2 слова

	// TODO: Проверь одно слово без пробелов
}

// TestTruncate проверяет функцию Truncate
func TestTruncate(t *testing.T) {
	// TODO: Проверь строку короче maxLen (не должна обрезаться)
	// result := Truncate("Hi", 10)
	// if result != "Hi" {
	//     t.Errorf("Truncate(%q, 10) = %q; want %q", "Hi", result, "Hi")
	// }

	// TODO: Проверь строку длиннее maxLen (должна обрезаться с "...")
	// "Hello, World!" с maxLen=10 -> "Hello, ..."

	// TODO: Проверь крайний случай maxLen <= 3
	// При maxLen=3 должно вернуться "..."

	// TODO: Проверь строку равную maxLen (не должна обрезаться)
}
