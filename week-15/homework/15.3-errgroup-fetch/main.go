package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
	// Раскомментируй после реализации
	// "golang.org/x/sync/errgroup"
)

// DataItem представляет загруженные данные
type DataItem struct {
	URL      string
	Data     string
	Size     int
	Duration time.Duration
}

// AggregatedResult содержит агрегированные данные
type AggregatedResult struct {
	TotalItems   int
	TotalSize    int
	TotalTime    time.Duration
	SuccessCount int
	ErrorCount   int
	Items        []DataItem
}

// simulateFetch эмулирует загрузку данных из URL.
// 80% успех, 20% ошибка.
func simulateFetch(ctx context.Context, url string) (DataItem, error) {
	// Случайная задержка от 50ms до 300ms
	delay := time.Duration(50+rand.IntN(250)) * time.Millisecond

	select {
	case <-time.After(delay):
		// 20% шанс ошибки
		if rand.IntN(5) == 0 {
			return DataItem{}, fmt.Errorf("fetch failed for %s", url)
		}

		// Случайный размер данных
		size := 100 + rand.IntN(200)
		data := fmt.Sprintf("Data from %s (%d bytes)", url, size)

		return DataItem{
			URL:      url,
			Data:     data,
			Size:     size,
			Duration: delay,
		}, nil

	case <-ctx.Done():
		return DataItem{}, ctx.Err()
	}
}

// fetchAll загружает данные из всех URL параллельно.
// При первой ошибке отменяет остальные запросы.
// Возвращает все успешные результаты или первую ошибку.
func fetchAll(ctx context.Context, urls []string) ([]DataItem, error) {
	// TODO: реализуй функцию
	// 1. Создай errgroup с контекстом: g, ctx := errgroup.WithContext(ctx)
	// 2. Создай слайс результатов: results := make([]DataItem, len(urls))
	// 3. Для каждого URL запусти горутину через g.Go:
	//    - Захвати переменные i, url := i, url
	//    - Вызови simulateFetch(ctx, url)
	//    - При успехе: results[i] = item
	//    - Верни ошибку (или nil при успехе)
	// 4. Дождись завершения: err := g.Wait()
	// 5. При ошибке верни nil, err
	// 6. При успехе верни results, nil

	_ = ctx  // удали после реализации
	_ = urls // удали после реализации

	return nil, fmt.Errorf("not implemented")
}

// fetchAllWithLimit загружает данные с ограничением параллелизма.
// limit - максимальное количество одновременных запросов.
// Все запросы выполняются, даже если некоторые завершаются ошибкой.
func fetchAllWithLimit(ctx context.Context, urls []string, limit int) ([]DataItem, []error) {
	// TODO: реализуй функцию
	// 1. Создай errgroup БЕЗ контекста: g := new(errgroup.Group)
	// 2. Установи лимит: g.SetLimit(limit)
	// 3. Создай слайс результатов и слайс ошибок
	// 4. Создай mutex для защиты слайсов: var mu sync.Mutex
	// 5. Для каждого URL запусти горутину через g.Go:
	//    - Вызови simulateFetch(ctx, url)
	//    - При успехе: добавь в results (под mutex)
	//    - При ошибке: добавь в errors (под mutex)
	//    - ВАЖНО: возвращай nil, чтобы не останавливать группу!
	// 6. g.Wait() - игнорируй возвращаемое значение
	// 7. Верни results, errors

	_ = ctx   // удали после реализации
	_ = urls  // удали после реализации
	_ = limit // удали после реализации

	// Заглушка для sync.Mutex чтобы импорт не ломался
	var mu sync.Mutex
	_ = mu

	return nil, nil
}

// aggregateData собирает данные из нескольких источников и агрегирует результат.
// Продолжает работу даже при ошибках в отдельных источниках.
func aggregateData(ctx context.Context, sources []string) (AggregatedResult, error) {
	// TODO: реализуй функцию
	// 1. Вызови fetchAllWithLimit с лимитом 3
	// 2. Создай AggregatedResult
	// 3. Заполни поля:
	//    - TotalItems = len(items)
	//    - TotalSize = сумма Size всех items
	//    - SuccessCount = len(items)
	//    - ErrorCount = len(errs)
	//    - Items = items
	// 4. Верни результат
	//
	// Бонус: добавь TotalTime как максимальное время из всех Duration

	_ = ctx     // удали после реализации
	_ = sources // удали после реализации

	return AggregatedResult{}, fmt.Errorf("not implemented")
}

func main() {
	fmt.Println("=== ДЗ 15.3: errgroup Fetch ===")
	fmt.Println()

	urls := []string{
		"https://api1.example.com/data",
		"https://api2.example.com/data",
		"https://api3.example.com/data",
		"https://api4.example.com/data",
		"https://api5.example.com/data",
	}

	ctx := context.Background()

	// Тест 1: Параллельная загрузка (все успешно - несколько попыток)
	fmt.Println("--- Тест 1: Параллельная загрузка ---")
	start := time.Now()
	items, err := fetchAll(ctx, urls)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		for i, item := range items {
			fmt.Printf("[%d] %s: %d bytes (%v)\n",
				i+1, item.URL, item.Size, item.Duration.Round(time.Millisecond))
		}
		fmt.Printf("Общее время: %v (параллельно!)\n", elapsed.Round(time.Millisecond))
	}

	// Тест 2: Загрузка с лимитом
	fmt.Println("\n--- Тест 2: Загрузка с лимитом (limit=2) ---")

	moreURLs := make([]string, 10)
	for i := range moreURLs {
		moreURLs[i] = fmt.Sprintf("https://api%d.example.com/data", i+1)
	}

	start = time.Now()
	items, errs := fetchAllWithLimit(ctx, moreURLs, 2)
	elapsed = time.Since(start)

	fmt.Printf("Загружено: %d/%d\n", len(items), len(moreURLs))
	fmt.Printf("Ошибок: %d\n", len(errs))
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("  - %v\n", e)
		}
	}
	fmt.Printf("Общее время: %v\n", elapsed.Round(time.Millisecond))

	// Тест 3: Агрегация данных
	fmt.Println("\n--- Тест 3: Агрегация данных ---")
	sources := []string{
		"https://source1.example.com",
		"https://source2.example.com",
		"https://source3.example.com",
		"https://source4.example.com",
		"https://source5.example.com",
		"https://source6.example.com",
	}

	start = time.Now()
	result, err := aggregateData(ctx, sources)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Результат агрегации:")
		fmt.Printf("  Всего элементов: %d\n", result.TotalItems)
		fmt.Printf("  Общий размер: %d bytes\n", result.TotalSize)
		fmt.Printf("  Успешных: %d\n", result.SuccessCount)
		fmt.Printf("  Ошибок: %d\n", result.ErrorCount)
		fmt.Printf("  Время: %v\n", elapsed.Round(time.Millisecond))
	}

	// Тест 4: Отмена контекста во время загрузки
	fmt.Println("\n--- Тест 4: Отмена контекста во время загрузки ---")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start = time.Now()
	items, err = fetchAll(ctx, urls)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("Ошибка через %v: %v\n", elapsed.Round(time.Millisecond), err)
	} else {
		fmt.Printf("Неожиданный успех: %d items\n", len(items))
	}

	fmt.Println("\n=== Тесты завершены ===")
}
