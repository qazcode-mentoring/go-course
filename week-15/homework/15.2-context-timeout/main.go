package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

// Response содержит результат запроса
type Response struct {
	URL        string
	StatusCode int
	Body       string
	Duration   time.Duration
}

// simulateHTTPRequest эмулирует HTTP-запрос.
// Возвращает случайный результат:
// - 70% успех (статус 200)
// - 20% ошибка сервера (статус 500)
// - 10% таймаут (очень долгий запрос)
func simulateHTTPRequest(ctx context.Context, url string) (Response, error) {
	// Определяем тип ответа
	outcome := rand.IntN(10)

	var delay time.Duration
	var statusCode int
	var body string

	switch {
	case outcome < 7: // 70% - успех
		delay = time.Duration(50+rand.IntN(200)) * time.Millisecond
		statusCode = 200
		body = fmt.Sprintf("Response from %s", url)
	case outcome < 9: // 20% - ошибка сервера
		delay = time.Duration(50+rand.IntN(100)) * time.Millisecond
		statusCode = 500
		body = "Internal Server Error"
	default: // 10% - таймаут
		delay = time.Duration(2000+rand.IntN(1000)) * time.Millisecond
		statusCode = 0
		body = ""
	}

	start := time.Now()

	select {
	case <-time.After(delay):
		if statusCode == 0 {
			return Response{}, errors.New("connection timeout")
		}
		return Response{
			URL:        url,
			StatusCode: statusCode,
			Body:       body,
			Duration:   time.Since(start),
		}, nil
	case <-ctx.Done():
		return Response{}, ctx.Err()
	}
}

// fetchWithTimeout выполняет запрос с таймаутом.
// Если запрос не успевает за timeout - возвращает ошибку.
func fetchWithTimeout(ctx context.Context, url string, timeout time.Duration) (Response, error) {
	// TODO: реализуй функцию
	// 1. Создай дочерний контекст с таймаутом: ctx, cancel := context.WithTimeout(ctx, timeout)
	// 2. Не забудь defer cancel()
	// 3. Вызови simulateHTTPRequest с новым контекстом
	// 4. Проверь StatusCode - если 500, верни ошибку "server error: 500"
	// 5. Верни результат
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := simulateHTTPRequest(ctx, url)
	if err != nil {
		return Response{}, err
	}

	if response.StatusCode == 500 {
		return Response{}, fmt.Errorf("server error: 500, %w", err)
	}

	return response, nil
}

// fetchWithRetry выполняет запрос с повторными попытками.
// maxRetries - максимальное количество попыток
// retryDelay - задержка между попытками
// Общий таймаут контролируется через ctx.
// Возвращает результат первого успешного запроса или последнюю ошибку.
func fetchWithRetry(ctx context.Context, url string, maxRetries int, retryDelay time.Duration) (Response, error) {
	// TODO: реализуй функцию
	// 1. Создай переменную для хранения последней ошибки
	// 2. Цикл for attempt := range maxRetries:
	//    a) Проверь ctx.Err() - если контекст отменён, верни ошибку
	//    b) Вызови fetchWithTimeout с таймаутом 500ms
	//    c) Если успех - верни результат
	//    d) Сохрани ошибку
	//    e) Если не последняя попытка - подожди retryDelay с проверкой контекста:
	//       select {
	//       case <-time.After(retryDelay):
	//       case <-ctx.Done():
	//           return Response{}, ctx.Err()
	//       }
	// 3. Верни последнюю ошибку

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if ctx.Err() != nil {
			return Response{}, ctx.Err()
		}

		response, err := fetchWithTimeout(ctx, url, 500*time.Millisecond)
		if response.StatusCode == 200 {
			return response, nil
		}

		lastErr = err

		select {
		case <-time.After(retryDelay):
		case <-ctx.Done():
			return Response{}, ctx.Err()
		}
	}

	return Response{}, lastErr
}

// fetchSequential выполняет последовательные запросы к нескольким URL.
// Все запросы должны уложиться в общий таймаут ctx.
// Если один запрос не успевает - остальные не выполняются.
// Возвращает все успешные ответы и первую ошибку (если была).
func fetchSequential(ctx context.Context, urls []string, perRequestTimeout time.Duration) ([]Response, error) {
	// TODO: реализуй функцию
	// 1. Создай слайс для результатов
	// 2. Для каждого URL (for i, url := range urls):
	//    a) Проверь ctx.Err() - если контекст отменён, верни результаты и ошибку
	//    b) Опционально: проверь оставшееся время через deadline, _ := ctx.Deadline()
	//    c) Вызови fetchWithTimeout
	//    d) Если ошибка - верни собранные результаты и ошибку
	//    e) Добавь успешный результат в слайс
	// 3. Верни все результаты и nil
	var responses []Response
	for _, url := range urls {
		if ctx.Err() != nil {
			return responses, ctx.Err()
		}

		deadline, _ := ctx.Deadline()
		fmt.Println(deadline)
		response, err := fetchWithTimeout(ctx, url, perRequestTimeout)
		if response.StatusCode == 200 {
			responses = append(responses, response)
		} else {
			return responses, err
		}

	}

	return responses, nil
}

func main() {
	fmt.Println("=== ДЗ 15.2: Context Timeout ===")
	fmt.Println()

	// Тест 1: Успешный запрос с таймаутом
	fmt.Println("--- Тест 1: Запрос с таймаутом (500ms) ---")
	ctx := context.Background()
	resp, err := fetchWithTimeout(ctx, "https://api.example.com", 500*time.Millisecond)
	if err != nil {
		fmt.Printf("[FAIL] %v\n", err)
	} else {
		fmt.Printf("[OK] Статус %d за %v: %s\n",
			resp.StatusCode, resp.Duration.Round(time.Millisecond), resp.Body)
	}

	// Тест 2: Запрос с очень коротким таймаутом
	fmt.Println("\n--- Тест 2: Запрос с коротким таймаутом (10ms) ---")
	resp, err = fetchWithTimeout(ctx, "https://api.example.com", 10*time.Millisecond)
	if err != nil {
		fmt.Printf("[FAIL] %v\n", err)
	} else {
		fmt.Printf("[OK] Статус %d: %s\n", resp.StatusCode, resp.Body)
	}

	// Тест 3: Запрос с повторами
	fmt.Println("\n--- Тест 3: Запрос с повторами (max 5, delay 100ms) ---")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	resp, err = fetchWithRetry(ctx, "https://flaky-api.example.com", 5, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("[FAIL] после %v: %v\n", time.Since(start).Round(time.Millisecond), err)
	} else {
		fmt.Printf("[OK] Успех за %v: %s\n", time.Since(start).Round(time.Millisecond), resp.Body)
	}

	// Тест 4: Запрос с повторами и коротким общим таймаутом
	fmt.Println("\n--- Тест 4: Повторы с коротким общим таймаутом (300ms) ---")
	ctx, cancel = context.WithTimeout(context.Background(), 300*time.Millisecond)

	start = time.Now()
	resp, err = fetchWithRetry(ctx, "https://slow-api.example.com", 10, 50*time.Millisecond)
	if err != nil {
		fmt.Printf("[FAIL] после %v: %v\n", time.Since(start).Round(time.Millisecond), err)
	} else {
		fmt.Printf("[OK] Успех за %v: %s\n", time.Since(start).Round(time.Millisecond), resp.Body)
	}
	cancel()

	// Тест 5: Последовательные запросы
	fmt.Println("\n--- Тест 5: Последовательные запросы (deadline 2s) ---")
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel = context.WithDeadline(context.Background(), deadline)
	defer cancel()

	urls := []string{
		"https://api1.example.com",
		"https://api2.example.com",
		"https://api3.example.com",
		"https://api4.example.com",
	}

	start = time.Now()
	responses, err := fetchSequential(ctx, urls, 400*time.Millisecond)

	for i, resp := range responses {
		fmt.Printf("[%d/%d] %s: %d (%v)\n",
			i+1, len(urls), resp.URL, resp.StatusCode, resp.Duration.Round(time.Millisecond))
	}

	if err != nil {
		fmt.Printf("Прервано: %v\n", err)
	}
	fmt.Printf("Итого: %d/%d за %v\n", len(responses), len(urls), time.Since(start).Round(time.Millisecond))

	// Тест 6: Последовательные запросы с очень коротким дедлайном
	fmt.Println("\n--- Тест 6: Последовательные запросы (deadline 200ms) ---")
	ctx, cancel = context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start = time.Now()
	responses, err = fetchSequential(ctx, urls, 100*time.Millisecond)
	fmt.Printf("Получено %d/%d ответов за %v\n",
		len(responses), len(urls), time.Since(start).Round(time.Millisecond))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	fmt.Println("\n=== Тесты завершены ===")
}
