# ДЗ 15.2: HTTP-подобный запрос с Timeout

## Цель
Научиться использовать `context.WithTimeout` и `context.WithDeadline` для ограничения времени выполнения операций.

## Что нужно сделать

Реализовать HTTP-подобный клиент с таймаутами на разных уровнях:

1. **`fetchWithTimeout`** — выполнение запроса с таймаутом
2. **`fetchWithRetry`** — запрос с повторными попытками и общим таймаутом
3. **`fetchSequential`** — последовательные запросы с общим дедлайном

## Сигнатуры

```go
// Response содержит результат запроса
type Response struct {
    URL        string
    StatusCode int
    Body       string
    Duration   time.Duration
}

// fetchWithTimeout выполняет запрос с таймаутом.
// Если запрос не успевает за timeout - возвращает ошибку.
func fetchWithTimeout(ctx context.Context, url string, timeout time.Duration) (Response, error)

// fetchWithRetry выполняет запрос с повторными попытками.
// maxRetries - максимальное количество попыток
// retryDelay - задержка между попытками
// Общий таймаут контролируется через ctx.
// Возвращает результат первого успешного запроса или последнюю ошибку.
func fetchWithRetry(ctx context.Context, url string, maxRetries int, retryDelay time.Duration) (Response, error)

// fetchSequential выполняет последовательные запросы к нескольким URL.
// Все запросы должны уложиться в общий таймаут ctx.
// Если один запрос не успевает - остальные не выполняются.
// Возвращает все успешные ответы и первую ошибку (если была).
func fetchSequential(ctx context.Context, urls []string, perRequestTimeout time.Duration) ([]Response, error)
```

## Теория: Timeout vs Deadline

```go
// WithTimeout - задаёт относительное время
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// Контекст отменится через 5 секунд от текущего момента

// WithDeadline - задаёт абсолютное время
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
// Контекст отменится в конкретный момент времени
```

**Когда что использовать:**
- `WithTimeout` — для одиночных операций с фиксированным временем
- `WithDeadline` — когда есть общий дедлайн на несколько операций

## Эмуляция HTTP-запроса

Для эмуляции запроса используй `simulateHTTPRequest`, которая уже реализована в шаблоне:

```go
// Возвращает случайный результат:
// - 70% успех (статус 200)
// - 20% ошибка сервера (статус 500)
// - 10% таймаут (очень долгий запрос)
func simulateHTTPRequest(ctx context.Context, url string) (Response, error)
```

## Пример использования

```go
func main() {
    // Запрос с таймаутом
    ctx := context.Background()
    resp, err := fetchWithTimeout(ctx, "https://api.example.com/data", 500*time.Millisecond)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Printf("Ответ: %d, %s\n", resp.StatusCode, resp.Body)
    }

    // Запрос с повторами
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    resp, err = fetchWithRetry(ctx, "https://flaky-api.example.com/data", 3, 200*time.Millisecond)
    if err != nil {
        fmt.Printf("Все попытки неуспешны: %v\n", err)
    } else {
        fmt.Printf("Успех: %s\n", resp.Body)
    }

    // Последовательные запросы с общим дедлайном
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel = context.WithDeadline(context.Background(), deadline)
    defer cancel()

    urls := []string{
        "https://api1.example.com",
        "https://api2.example.com",
        "https://api3.example.com",
    }

    responses, err := fetchSequential(ctx, urls, 500*time.Millisecond)
    fmt.Printf("Получено %d ответов\n", len(responses))
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    }
}
```

## Ожидаемый вывод (пример)

```
=== Тест 1: Запрос с таймаутом ===
[OK] Ответ 200 за 234ms: Response from https://api.example.com

=== Тест 2: Запрос с коротким таймаутом ===
[FAIL] context deadline exceeded

=== Тест 3: Запрос с повторами ===
Попытка 1: ошибка 500, повтор через 200ms...
Попытка 2: успех!
Ответ: Response from https://flaky-api.example.com

=== Тест 4: Последовательные запросы ===
[1/3] https://api1.example.com: 200 OK (156ms)
[2/3] https://api2.example.com: 200 OK (423ms)
[3/3] https://api3.example.com: deadline exceeded
Получено 2/3 ответов
```

## Подсказки

- В `fetchWithTimeout` создай дочерний контекст с `context.WithTimeout`
- В `fetchWithRetry` используй цикл с `time.After` для задержки между попытками
- Проверяй `ctx.Err()` перед каждой попыткой — если родительский контекст отменён, новые попытки бессмысленны
- В `fetchSequential` проверяй оставшееся время через `ctx.Deadline()` перед каждым запросом
- Используй `errors.Is(err, context.DeadlineExceeded)` для проверки типа ошибки

## Критерии выполнения

- [ ] `fetchWithTimeout` возвращает ошибку при превышении таймаута
- [ ] `fetchWithRetry` делает указанное количество попыток
- [ ] `fetchWithRetry` прекращает попытки при отмене родительского контекста
- [ ] `fetchSequential` останавливается при исчерпании общего дедлайна
- [ ] Все функции корректно обрабатывают отмену контекста
- [ ] Программа проходит проверку `go run -race`
