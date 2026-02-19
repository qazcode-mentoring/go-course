package main

import (
	"fmt"
	"strings"
)

// WordCount возвращает map с количеством каждого слова в тексте
func WordCount(text string) map[string]int {
	// TODO: реализуй функцию
	// Используй strings.Fields для разбиения на слова
	// Пройди по словам и увеличивай счётчик в map
	str := strings.Fields(text)
	counts := map[string]int{}

	for i := 0; i < len(str); i++ {
		word := str[i]
		counts[word] += 1
	}

	return counts
}

// CharCount возвращает map с количеством каждого символа
func CharCount(text string) map[rune]int {
	// TODO: реализуй функцию
	// Пройди по строке с range — получишь rune
	charCount := map[rune]int{}

	for _, r := range text {
		charCount[r]++
	}
	return charCount
}

// MostFrequent возвращает самое часто встречающееся слово
func MostFrequent(text string) string {
	// TODO: реализуй функцию
	// Используй WordCount, затем найди слово с максимальным значением
	wordCount := WordCount(text)

	var frequentWord string
	maxCount := 1
	for key, value := range wordCount {
		if maxCount < value {
			frequentWord = key
			maxCount = value
		}
	}
	return frequentWord
}

func main() {
	text := "go go go python java go python"

	fmt.Println("=== Подсчёт слов ===")
	fmt.Println("Текст:", text)
	fmt.Println("Количество слов:", WordCount(text))
	fmt.Println("Самое частое слово:", MostFrequent(text))

	fmt.Println("\n=== Подсчёт символов ===")
	word := "hello"
	fmt.Println("Слово:", word)
	fmt.Println("Количество символов:", CharCount(word))
}
