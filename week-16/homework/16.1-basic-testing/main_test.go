package main

import "testing"

// TestReverse проверяет функцию Reverse
func TestReverse(t *testing.T) {
	// TODO: Проверь переворот обычной строки
	// TODO: Проверь пустую строку
	// TODO: Проверь строку из одного символа
	// TODO: Проверь строку с пробелами

	res1 := Reverse("hello")
	if res1 != "olleh" {
		t.Errorf("Reverse(%q) = %q; want %q", "hello", res1, "olleh")
	}

	res2 := Reverse("")
	if res2 != "" {
		t.Errorf("Reverse(%q) = %q; want %q", "", res2, "")
	}

	res3 := Reverse("a")
	if res3 != "a" {
		t.Errorf("Reverse(%q) = %q; want %q", "a", res3, "a")
	}

	res4 := Reverse("a b")
	if res4 != "b a" {
		t.Errorf("Reverse(%q) = %q; want %q", "a b", res4, "b a")
	}
}

// TestIsPalindrome проверяет функцию IsPalindrome
func TestIsPalindrome(t *testing.T) {
	// TODO: Проверь палиндром "radar"
	// TODO: Проверь палиндром "level"
	// TODO: Проверь не-палиндром "hello"
	// TODO: Проверь пустую строку (должна быть палиндромом)
	// TODO: Проверь строку из одного символа

	if !IsPalindrome("radar") {
		t.Errorf("IsPalindrome(%q) = false; want true", "radar")
	}

	if !IsPalindrome("level") {
		t.Errorf("IsPalindrome(%q) = false; want true", "level")
	}

	if IsPalindrome("hello") {
		t.Errorf("IsPalindrome(%q) = true; want false", "hello")
	}

	if !IsPalindrome("") {
		t.Errorf("IsPalindrome(%q) = false; want true", "")
	}

	if !IsPalindrome("a") {
		t.Errorf("IsPalindrome(%q) = false; want true", "a")
	}
}

// TestCountWords проверяет функцию CountWords
func TestCountWords(t *testing.T) {
	// TODO: Проверь обычное предложение
	// TODO: Проверь пустую строку (должно быть 0 слов)
	// TODO: Проверь строку с множественными пробелами
	// TODO: Проверь одно слово без пробелов

	res1 := CountWords("hello world")
	if res1 != 2 {
		t.Errorf("CountWords(%q) = %d; want %d", "hello world", res1, 2)
	}

	res2 := CountWords("")
	if res2 != 0 {
		t.Errorf("CountWords(%q) = %d; want %d", "", res2, 0)
	}

	res3 := CountWords("  multiple   spaces")
	if res3 != 2 {
		t.Errorf("CountWords(%q) = %d; want %d", "  multiple   spaces", res3, 2)
	}

	res4 := CountWords("hi")
	if res4 != 1 {
		t.Errorf("CountWords(%q) = %d; want %d", "hi", res4, 1)
	}
}

// TestTruncate проверяет функцию Truncate
func TestTruncate(t *testing.T) {
	// TODO: Проверь строку короче maxLen (не должна обрезаться)
	// TODO: Проверь строку длиннее maxLen (должна обрезаться с "...")
	// TODO: Проверь крайний случай maxLen <= 3
	// TODO: Проверь строку равную maxLen (не должна обрезаться)

	res1 := Truncate("hi", 3)
	if res1 != "hi" {
		t.Errorf("Truncate(%q, 3) = %q; want %q", "hi", res1, "hi")
	}

	res2 := Truncate("hello world", 8)
	if res2 != "hello..." {
		t.Errorf("Truncate(%q, 8) = %q; want %q", "hello world", res2, "hello...")
	}

	res3 := Truncate("long string", 3)
	if res3 != "..." {
		t.Errorf("Truncate(%q, 3) = %q; want %q", "long string", res3, "...")
	}

	res4 := Truncate("golang", 6)
	if res4 != "golang" {
		t.Errorf("Truncate(%q, 6) = %q; want %q", "golang", res4, "golang")
	}
}
