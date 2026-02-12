package main

import "testing"

// TestReverse проверяет функцию Reverse
func TestReverse(t *testing.T) {
	// TODO: Проверь переворот обычной строки
	result := Reverse("hello")
	if result != "olleh" {
		t.Errorf("Reverse(%q) = %q; want %q", "hello", result, "olleh")
	}

	// TODO: Проверь пустую строку
	// Пустая строка при перевороте остаётся пустой
	result = Reverse("")
	if result != "" {
		t.Errorf("Reverse(%q) = %q; want %q", "", result, "")
	}

	// TODO: Проверь строку из одного символа
	// Один символ при перевороте остаётся тем же
	result = Reverse("a")
	if result != "a" {
		t.Errorf("Reverse(%q) = %q; want %q", "a", result, "a")
	}

	// TODO: Проверь строку с пробелами
	// "a b" -> "b a"
	result = Reverse("a b")
	if result != "b a" {
		t.Errorf("Reverse(%q) = %q; want %q", "b a", result, "b a")
	}
}

// TestIsPalindrome проверяет функцию IsPalindrome
func TestIsPalindrome(t *testing.T) {
	// TODO: Проверь палиндром "radar"
	if !IsPalindrome("radar") {
		t.Error("IsPalindrome(\"radar\") should be true")
	}

	// TODO: Проверь палиндром "level"
	if !IsPalindrome("level") {
		t.Error("IsPalindrome(\"radar\") should be true")
	}
	// TODO: Проверь не-палиндром "hello"
	if IsPalindrome("hello") {
		t.Error("IsPalindrome(\"hello\") should be false")
	}

	// TODO: Проверь пустую строку (должна быть палиндромом)
	if !IsPalindrome("") {
		t.Error("IsPalindrome(\"\") should be true")
	}
	// TODO: Проверь строку из одного символа
	if !IsPalindrome("a") {
		t.Error("IsPalindrome(\"a\") should be true")
	}
}

// TestCountWords проверяет функцию CountWords
func TestCountWords(t *testing.T) {
	// TODO: Проверь обычное предложение
	result := CountWords("hello world")
	if result != 2 {
		t.Errorf("CountWords(%q) = %d; want %d", "hello world", result, 2)
	}

	// TODO: Проверь пустую строку (должно быть 0 слов)
	result = CountWords("")
	if result != 0 {
		t.Errorf("CountWords(%q) = %d; want %d", "", result, 2)
	}
	// TODO: Проверь строку с множественными пробелами
	// "  multiple   spaces  " должно дать 2 слова
	result = CountWords("  multiple   spaces  ")
	if result != 2 {
		t.Errorf("CountWords(%q) = %d; want %d", "  multiple   spaces  ", result, 2)
	}

	// TODO: Проверь одно слово без пробелов
	result = CountWords("islam")
	if result != 1 {
		t.Errorf("CountWords(%q) = %d; want %d", "islam", result, 1)
	}
}

// TestTruncate проверяет функцию Truncate
func TestTruncate(t *testing.T) {
	// TODO: Проверь строку короче maxLen (не должна обрезаться)
	result := Truncate("Hi", 10)
	if result != "Hi" {
		t.Errorf("Truncate(%q, 10) = %q; want %q", "Hi", result, "Hi")
	}

	// TODO: Проверь строку длиннее maxLen (должна обрезаться с "...")
	// "Hello, World!" с maxLen=10 -> "Hello, ..."
	result = Truncate("Hello, World!", 10)
	if result != "Hello, ..." {
		t.Errorf("Truncate(%q, 10) = %q; want %q", "Hello, World!", result, "Hello, ...")
	}

	// TODO: Проверь крайний случай maxLen <= 3
	// При maxLen=3 должно вернуться "..."
	result = Truncate("Hello", 3)
	if result != "..." {
		t.Errorf("Truncate(%q, 3) = %q; want %q", "Hello", result, "...")
	}

	// TODO: Проверь строку равную maxLen (не должна обрезаться)
	result = Truncate("islam", 5)
	if result != "islam" {
		t.Errorf("Truncate(%q, 3) = %q; want %q", "islam", result, "islam")
	}
}
