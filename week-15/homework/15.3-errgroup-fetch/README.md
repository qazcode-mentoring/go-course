# ДЗ 15.3: Параллельные запросы с errgroup

## Цель
Научиться использовать пакет `golang.org/x/sync/errgroup` для параллельного выполнения задач с обработкой ошибок.

## Подготовка

Установи пакет errgroup:

```bash
go get golang.org/x/sync/errgroup
```

## Что нужно сделать

Реализовать систему параллельной загрузки данных с использованием errgroup:

1. **`fetchAll`** — параллельная загрузка нескольких URL
2. **`fetchAllWithLimit`** — загрузка с ограничением параллелизма
3. **`aggregateData`** — сбор данных из нескольких источников с агрегацией

## Сигнатуры

```go
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

// fetchAll загружает данные из всех URL параллельно.
// При первой ошибке отменяет остальные запросы.
// Возвращает все успешные результаты или первую ошибку.
func fetchAll(ctx context.Context, urls []string) ([]DataItem, error)

// fetchAllWithLimit загружает данные с ограничением параллелизма.
// limit - максимальное количество одновременных запросов.
// Все запросы выполняются, даже если некоторые завершаются ошибкой.
func fetchAllWithLimit(ctx context.Context, urls []string, limit int) ([]DataItem, []error)

// aggregateData собирает данные из нескольких источников и агрегирует результат.
// Продолжает работу даже при ошибках в отдельных источниках.
func aggregateData(ctx context.Context, sources []string) (AggregatedResult, error)
```

## Теория: errgroup

```go
import "golang.org/x/sync/errgroup"

// Базовое использование
g, ctx := errgroup.WithContext(ctx)

for _, url := range urls {
    url := url // Захват переменной (Go <1.22)
    g.Go(func() error {
        return fetch(ctx, url)
    })
}

if err := g.Wait(); err != nil {
    return err // Первая ошибка
}

// С ограничением параллелизма (Go 1.20+)
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(10) // Максимум 10 горутин одновременно

// С TryGo (неблокирующий запуск)
if g.TryGo(func() error { ... }) {
    // Горутина запущена
} else {
    // Достигнут лимит, горутина не запущена
}
```

### Важные особенности errgroup

1. **Первая ошибка отменяет контекст**: при использовании `WithContext` первая возвращённая ошибка отменит ctx
2. **Wait() возвращает первую ошибку**: остальные ошибки игнорируются
3. **SetLimit() ограничивает параллелизм**: горутины сверх лимита будут ждать
4. **Безопасен для concurrent использования**: можно вызывать Go() из разных горутин

## Пример использования

```go
func main() {
    urls := []string{
        "https://api1.example.com/data",
        "https://api2.example.com/data",
        "https://api3.example.com/data",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Параллельная загрузка
    items, err := fetchAll(ctx, urls)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }

    for _, item := range items {
        fmt.Printf("%s: %d bytes\n", item.URL, item.Size)
    }

    // С ограничением параллелизма
    items, errs := fetchAllWithLimit(ctx, urls, 2)
    fmt.Printf("Загружено: %d, ошибок: %d\n", len(items), len(errs))

    // Агрегация
    result, err := aggregateData(ctx, urls)
    fmt.Printf("Всего: %d элементов, %d байт\n", result.TotalItems, result.TotalSize)
}
```

## Ожидаемый вывод (пример)

```
=== Тест 1: Параллельная загрузка (все успешно) ===
[1] https://api1.example.com: 156 bytes (89ms)
[2] https://api2.example.com: 203 bytes (124ms)
[3] https://api3.example.com: 178 bytes (67ms)
Общее время: 127ms (параллельно!)

=== Тест 2: Параллельная загрузка с ошибкой ===
Первая ошибка: fetch failed for https://api2.example.com

=== Тест 3: Загрузка с лимитом (limit=2) ===
Загружено 8/10 элементов
Ошибок: 2
Общее время: 523ms

=== Тест 4: Агрегация данных ===
Результат:
  Всего элементов: 5
  Общий размер: 1024 bytes
  Успешных: 4, ошибок: 1
  Время: 234ms
```

## Подсказки

- В `fetchAll` используй `errgroup.WithContext` для автоматической отмены при ошибке
- Для записи результатов в слайс используй индекс: `results[i] = item`
- В `fetchAllWithLimit` используй `g.SetLimit(limit)`
- Для сбора всех ошибок (не только первой) используй отдельный слайс с mutex
- В `aggregateData` не используй `errgroup.WithContext` если хочешь продолжить работу после ошибки

## Дополнительно: Паттерн сбора всех ошибок

```go
var (
    mu     sync.Mutex
    errors []error
)

g := new(errgroup.Group)
g.SetLimit(limit)

for i, url := range urls {
    i, url := i, url
    g.Go(func() error {
        item, err := fetch(ctx, url)
        if err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
            return nil // Не останавливаем группу!
        }
        results[i] = item
        return nil
    })
}

g.Wait() // Всегда nil, т.к. мы возвращаем nil
return results, errors
```

## Критерии выполнения

- [ ] `fetchAll` выполняет запросы параллельно
- [ ] `fetchAll` возвращает первую ошибку и отменяет остальные запросы
- [ ] `fetchAllWithLimit` соблюдает ограничение параллелизма
- [ ] `fetchAllWithLimit` собирает все ошибки, не только первую
- [ ] `aggregateData` продолжает работу при ошибках
- [ ] `aggregateData` корректно агрегирует статистику
- [ ] Программа проходит проверку `go run -race`
