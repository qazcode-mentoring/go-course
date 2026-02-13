package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
	// Раскомментируй после реализации fetchMultiple
	// "sync"
)

// FetchResult содержит результат загрузки одного URL
type FetchResult struct {
	URL   string
	Data  string
	Error error
}

// simulateFetch эмулирует загрузку данных с URL.
// Возвращает канал, из которого можно прочитать результат.
// Загрузка занимает случайное время от 100ms до 2s.
func simulateFetch(url string) <-chan string {
	ch := make(chan string)
	go func() {
		// Случайная задержка от 100ms до 2s
		delay := time.Duration(100+rand.Intn(1900)) * time.Millisecond
		time.Sleep(delay)
		ch <- fmt.Sprintf("Data from %s (took %v)", url, delay)
	}()
	return ch
}

// fetchWithTimeout эмулирует загрузку данных с указанного URL.
// Возвращает данные или ошибку по таймауту.
func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
	// TODO: реализуй функцию
	// 1. Вызови simulateFetch(url) для получения канала с результатом
	// 2. Используй select для ожидания результата или таймаута:
	//    - case для чтения из канала результата
	//    - case для time.After(timeout)
	// 3. Верни данные или ошибку "timeout exceeded"

	ch := simulateFetch(url)

	select {
	case data := <-ch:
		return data, nil
	case <-time.After(timeout):
		return "", errors.New("timeout exceeded")
	}

	return "", errors.New("not implemented")
}

// fetchMultiple загружает данные из нескольких URL параллельно.
// Возвращает результаты по мере их получения через канал.
func fetchMultiple(urls []string, timeout time.Duration) <-chan FetchResult {
	// TODO: реализуй функцию
	// 1. Создай канал для результатов
	// 2. Создай WaitGroup для отслеживания горутин
	// 3. Для каждого URL:
	//    - Добавь 1 к WaitGroup
	//    - Запусти горутину, которая:
	//      a) Вызывает fetchWithTimeout
	//      b) Отправляет FetchResult в канал результатов
	//      c) Вызывает wg.Done()
	// 4. Запусти отдельную горутину, которая:
	//    - Ждёт завершения всех загрузок (wg.Wait())
	//    - Закрывает канал результатов
	// 5. Верни канал результатов

	results := make(chan FetchResult)
	wg := &sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			data, err := fetchWithTimeout(u, timeout)
			results <- FetchResult{
				URL:   u,
				Data:  data,
				Error: err,
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== ДЗ 14.1: Select Timeout ===")
	fmt.Println()

	// Тест 1: Одиночный запрос с таймаутом
	fmt.Println("--- Тест 1: Одиночный запрос ---")
	result, err := fetchWithTimeout("https://api.example.com/data", 500*time.Millisecond)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}

	// Тест 2: Короткий таймаут (скорее всего не успеет)
	fmt.Println("\n--- Тест 2: Короткий таймаут (50ms) ---")
	result, err = fetchWithTimeout("https://api.example.com/slow", 50*time.Millisecond)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}

	// Тест 3: Длинный таймаут (скорее всего успеет)
	fmt.Println("\n--- Тест 3: Длинный таймаут (3s) ---")
	result, err = fetchWithTimeout("https://api.example.com/fast", 3*time.Second)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}

	// Тест 4: Параллельные запросы
	fmt.Println("\n--- Тест 4: Параллельные запросы ---")
	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/posts",
		"https://api.example.com/comments",
		"https://slow-api.example.com/data",
	}

	results := fetchMultiple(urls, 1*time.Second)

	successCount := 0
	failCount := 0

	for r := range results {
		if r.Error != nil {
			fmt.Printf("[FAIL] %s: %v\n", r.URL, r.Error)
			failCount++
		} else {
			fmt.Printf("[OK] %s: %s\n", r.URL, r.Data)
			successCount++
		}
	}

	fmt.Printf("\nИтого: %d успешно, %d с ошибкой\n", successCount, failCount)
}
