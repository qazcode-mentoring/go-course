# ДЗ 14.1: Select Timeout

## Цель
Научиться использовать `select` для реализации таймаутов и отмены операций.

## Что нужно сделать

Реализовать функцию `fetchWithTimeout`, которая эмулирует загрузку данных с таймаутом, и функцию `fetchMultiple`, которая загружает данные из нескольких источников параллельно.

## Сигнатуры

```go
// fetchWithTimeout эмулирует загрузку данных с указанного URL.
// Возвращает данные или ошибку по таймауту.
// Параметры:
//   - url: адрес для "загрузки" (используется для эмуляции)
//   - timeout: максимальное время ожидания
// Возвращает:
//   - данные в виде строки
//   - ошибку, если превышен таймаут
func fetchWithTimeout(url string, timeout time.Duration) (string, error)

// fetchMultiple загружает данные из нескольких URL параллельно.
// Возвращает результаты по мере их получения через канал.
// Если загрузка URL не успевает за индивидуальный таймаут - возвращается ошибка для этого URL.
// Параметры:
//   - urls: список URL для загрузки
//   - timeout: таймаут для каждой отдельной загрузки
// Возвращает:
//   - канал с результатами (FetchResult)
func fetchMultiple(urls []string, timeout time.Duration) <-chan FetchResult
```

## Теория: Select для таймаутов

`select` позволяет ожидать несколько каналов одновременно. В сочетании с `time.After` это даёт простой способ реализации таймаутов:

```go
select {
case result := <-resultChan:
    // Получили результат вовремя
    return result, nil
case <-time.After(timeout):
    // Время вышло
    return "", errors.New("timeout")
}
```

## Эмуляция загрузки

Для эмуляции "загрузки" используй `simulateFetch`, которая уже реализована в шаблоне. Она возвращает канал, из которого можно прочитать результат через случайное время.

```go
// Уже реализована в шаблоне
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
```

## Пример использования

```go
func main() {
    rand.Seed(time.Now().UnixNano())

    // Одиночный запрос с таймаутом
    fmt.Println("=== Одиночный запрос ===")
    result, err := fetchWithTimeout("https://api.example.com/data", 500*time.Millisecond)
    if err != nil {
        fmt.Println("Ошибка:", err)
    } else {
        fmt.Println("Результат:", result)
    }

    // Параллельные запросы
    fmt.Println("\n=== Параллельные запросы ===")
    urls := []string{
        "https://api.example.com/users",
        "https://api.example.com/posts",
        "https://api.example.com/comments",
        "https://slow-api.example.com/data",
    }

    results := fetchMultiple(urls, 1*time.Second)
    for r := range results {
        if r.Error != nil {
            fmt.Printf("[FAIL] %s: %v\n", r.URL, r.Error)
        } else {
            fmt.Printf("[OK] %s: %s\n", r.URL, r.Data)
        }
    }
}
```

## Ожидаемый вывод (пример)

```
=== Одиночный запрос ===
Результат: Data from https://api.example.com/data (took 234ms)

=== Параллельные запросы ===
[OK] https://api.example.com/posts: Data from https://api.example.com/posts (took 156ms)
[OK] https://api.example.com/users: Data from https://api.example.com/users (took 423ms)
[FAIL] https://slow-api.example.com/data: timeout exceeded
[OK] https://api.example.com/comments: Data from https://api.example.com/comments (took 891ms)
```

Порядок вывода и конкретные значения будут меняться при каждом запуске.

## Подсказки

- Используй `select` с `case <-time.After(timeout)` для реализации таймаута
- В `fetchMultiple` запускай отдельную горутину для каждого URL
- Не забудь закрыть результирующий канал, когда все загрузки завершены
- Используй `sync.WaitGroup` для отслеживания завершения всех горутин

## Критерии выполнения

- [ ] Функция `fetchWithTimeout` корректно возвращает данные или ошибку таймаута
- [ ] Функция `fetchMultiple` запускает загрузки параллельно
- [ ] Результаты приходят по мере готовности (не ждём все сразу)
- [ ] Канал закрывается после получения всех результатов
- [ ] Программа проходит проверку `go run -race`
