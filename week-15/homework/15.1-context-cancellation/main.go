package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// TaskResult содержит результат выполнения задачи
type TaskResult struct {
	TaskID         int           // ID задачи
	CompletedSteps int           // Сколько шагов выполнено
	TotalSteps     int           // Сколько шагов было запланировано
	Duration       time.Duration // Сколько времени заняла задача
	Error          error         // Ошибка (nil если успех, ctx.Err() если отменено)
}

// SearchResult содержит результат поиска
type SearchResult struct {
	Source   string        // Источник, где найден результат
	Data     string        // Найденные данные
	Duration time.Duration // Время поиска
}

// longRunningTask эмулирует долгую операцию, которая выполняется steps шагов.
// Каждый шаг занимает stepDuration времени.
// Операция должна проверять контекст и завершаться досрочно при отмене.
// Возвращает количество выполненных шагов и ошибку (ctx.Err() при отмене).
func longRunningTask(ctx context.Context, taskID int, steps int, stepDuration time.Duration) (int, error) {
	// TODO: реализуй функцию
	// 1. Создай счётчик выполненных шагов
	// 2. Используй цикл for i := range steps
	// 3. На каждой итерации:
	//    - Проверь, не отменён ли контекст (select с ctx.Done() и default)
	//    - Если отменён - верни количество выполненных шагов и ctx.Err()
	//    - Если нет - подожди stepDuration и увеличь счётчик
	// 4. После завершения всех шагов верни их количество и nil

	completedSteps := 0
	for i := 0; i < steps; i++ {
		select {
		case <-ctx.Done():
			return completedSteps, ctx.Err()
		default:
			time.Sleep(stepDuration)
			completedSteps++
		}
	}

	return completedSteps, nil
}

// processWithCancellation запускает несколько задач параллельно.
// Если контекст отменён - все задачи должны завершиться.
// Возвращает результаты всех задач через канал.
func processWithCancellation(ctx context.Context, taskCount int) <-chan TaskResult {
	// TODO: реализуй функцию
	// 1. Создай канал для результатов
	// 2. Создай WaitGroup для отслеживания горутин
	// 3. Для каждой задачи (for i := range taskCount):
	//    - Добавь 1 к WaitGroup
	//    - Запусти горутину, которая:
	//      a) Засекает время начала
	//      b) Вызывает longRunningTask с 5 шагами по 100ms
	//      c) Отправляет TaskResult в канал
	//      d) Вызывает wg.Done()
	// 4. Запусти горутину, которая ждёт WaitGroup и закрывает канал
	// 5. Верни канал

	results := make(chan TaskResult)
	wg := &sync.WaitGroup{}

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			start := time.Now()
			completedSteps, err := longRunningTask(ctx, ii, 5, 100*time.Millisecond)
			results <- TaskResult{
				TaskID:         ii,
				CompletedSteps: completedSteps,
				TotalSteps:     5,
				Duration:       time.Since(start),
				Error:          err,
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// simulateSearch эмулирует поиск в источнике.
// Занимает случайное время от 50ms до 300ms.
func simulateSearch(ctx context.Context, source string) (string, error) {
	// Случайная задержка
	delay := time.Duration(50+rand.IntN(250)) * time.Millisecond

	select {
	case <-time.After(delay):
		// 20% шанс ошибки
		if rand.IntN(5) == 0 {
			return "", fmt.Errorf("search failed in %s", source)
		}
		return fmt.Sprintf("Result from %s", source), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// searchFirst выполняет поиск параллельно в нескольких источниках.
// Возвращает первый успешный результат и отменяет остальные поиски.
// Если все поиски завершились ошибкой - возвращает ошибку.
func searchFirst(ctx context.Context, query string, sources []string) (SearchResult, error) {
	// TODO: реализуй функцию
	// 1. Создай дочерний контекст с отменой: ctx, cancel := context.WithCancel(ctx)
	// 2. Не забудь defer cancel()
	// 3. Создай канал для результатов (буферизированный на len(sources))
	// 4. Для каждого источника запусти горутину:
	//    - Засеки время начала
	//    - Вызови simulateSearch
	//    - Отправь результат в канал (используй select с ctx.Done())
	// 5. Собери результаты:
	//    - Жди len(sources) результатов
	//    - При первом успешном: вызови cancel() и верни результат
	//    - Если все ошибки - верни ошибку "all searches failed"
	//
	// Подсказка: используй структуру для передачи через канал:
	type searchResponse struct {
		result SearchResult
		err    error
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	results := make(chan searchResponse, len(sources))

	for _, source := range sources {
		go func(src string) {
			start := time.Now()
			searchRes, err := simulateSearch(ctx, src)
			select {
			case <-ctx.Done():
				return
			case results <- searchResponse{
				result: SearchResult{
					Source:   src,
					Data:     searchRes,
					Duration: time.Since(start),
				},
				err: err,
			}:
			}
		}(source)
	}

	for i := 0; i < len(sources); i++ {
		select {
		case res := <-results:
			if res.err == nil {
				cancel()
				return res.result, nil
			}
		case <-ctx.Done():
			return SearchResult{}, ctx.Err()
		}
	}
	return SearchResult{}, fmt.Errorf("all searches failed")
}

func main() {
	fmt.Println("=== ДЗ 15.1: Context Cancellation ===")
	fmt.Println()

	// Тест 1: Задача без отмены (должна завершиться успешно)
	fmt.Println("--- Тест 1: Задача без отмены ---")
	ctx := context.Background()
	start := time.Now()
	completed, err := longRunningTask(ctx, 1, 5, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Задача завершена: %d/5 шагов за %v\n", completed, time.Since(start).Round(time.Millisecond))
	}

	// Тест 2: Задача с отменой через 250ms
	fmt.Println("\n--- Тест 2: Задача с отменой через 250ms ---")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(250 * time.Millisecond)
		fmt.Println("Отменяем задачу!")
		cancel()
	}()

	start = time.Now()
	completed, err = longRunningTask(ctx, 1, 10, 100*time.Millisecond)
	fmt.Printf("Задача: %d/10 шагов за %v", completed, time.Since(start).Round(time.Millisecond))
	if err != nil {
		fmt.Printf(" (ошибка: %v)\n", err)
	} else {
		fmt.Println()
	}

	// Тест 3: Несколько задач с отменой
	fmt.Println("\n--- Тест 3: Несколько задач с отменой ---")
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("Отменяем все задачи!")
		cancel()
	}()

	results := processWithCancellation(ctx, 3)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range results {
			if result.Error != nil {
				fmt.Printf("Задача %d: %d/%d шагов (%v)\n",
					result.TaskID, result.CompletedSteps, result.TotalSteps, result.Error)
			} else {
				fmt.Printf("Задача %d: завершена успешно за %v\n",
					result.TaskID, result.Duration.Round(time.Millisecond))
			}
		}
	}()
	wg.Wait()

	// Тест 4: Поиск первого результата
	fmt.Println("\n--- Тест 4: Поиск первого результата ---")
	sources := []string{"source-1", "source-2", "source-3", "source-4", "source-5"}

	ctx = context.Background()
	start = time.Now()
	searchResult, err := searchFirst(ctx, "test query", sources)
	if err != nil {
		fmt.Printf("Ошибка поиска: %v\n", err)
	} else {
		fmt.Printf("Найдено в '%s': %q за %v\n",
			searchResult.Source, searchResult.Data, searchResult.Duration.Round(time.Millisecond))
	}

	// Тест 5: Поиск с таймаутом контекста
	fmt.Println("\n--- Тест 5: Поиск с коротким таймаутом (10ms) ---")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	searchResult, err = searchFirst(ctx, "test query", sources)
	if err != nil {
		fmt.Printf("Ошибка поиска: %v\n", err)
	} else {
		fmt.Printf("Найдено: %s\n", searchResult.Data)
	}

	fmt.Println("\n=== Тесты завершены ===")
}
