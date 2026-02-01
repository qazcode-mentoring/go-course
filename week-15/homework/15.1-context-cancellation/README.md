# ДЗ 15.1: Отмена операций через Context

## Цель
Научиться использовать `context.Context` для отмены долгих операций и корректного завершения горутин.

## Что нужно сделать

Реализовать систему фоновых задач с возможностью отмены:

1. **`longRunningTask`** — имитация долгой операции, которая периодически проверяет контекст
2. **`processWithCancellation`** — запуск нескольких задач с возможностью отмены всех при отмене одной
3. **`searchFirst`** — поиск среди нескольких источников, возврат первого результата и отмена остальных

## Сигнатуры

```go
// longRunningTask эмулирует долгую операцию, которая выполняется steps шагов.
// Каждый шаг занимает stepDuration времени.
// Операция должна проверять контекст и завершаться досрочно при отмене.
// Возвращает количество выполненных шагов и ошибку (ctx.Err() при отмене).
func longRunningTask(ctx context.Context, taskID int, steps int, stepDuration time.Duration) (int, error)

// processWithCancellation запускает несколько задач параллельно.
// Если контекст отменён - все задачи должны завершиться.
// Возвращает результаты всех задач через канал.
func processWithCancellation(ctx context.Context, taskCount int) <-chan TaskResult

// searchFirst выполняет поиск параллельно в нескольких источниках.
// Возвращает первый успешный результат и отменяет остальные поиски.
// Если все поиски завершились ошибкой - возвращает ошибку.
func searchFirst(ctx context.Context, query string, sources []string) (SearchResult, error)
```

## Структуры

```go
// TaskResult содержит результат выполнения задачи
type TaskResult struct {
    TaskID        int           // ID задачи
    CompletedSteps int          // Сколько шагов выполнено
    TotalSteps    int           // Сколько шагов было запланировано
    Duration      time.Duration // Сколько времени заняла задача
    Error         error         // Ошибка (nil если успех, ctx.Err() если отменено)
}

// SearchResult содержит результат поиска
type SearchResult struct {
    Source   string        // Источник, где найден результат
    Data     string        // Найденные данные
    Duration time.Duration // Время поиска
}
```

## Теория: Проверка отмены контекста

Есть несколько способов проверить, отменён ли контекст:

```go
// Способ 1: через select (неблокирующий)
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // Контекст не отменён, продолжаем
}

// Способ 2: проверка ctx.Err() (неблокирующий)
if ctx.Err() != nil {
    return ctx.Err()
}

// Способ 3: блокирующее ожидание
<-ctx.Done() // Заблокируется до отмены контекста
return ctx.Err()
```

## Пример использования

```go
func main() {
    // Создаём контекст с отменой
    ctx, cancel := context.WithCancel(context.Background())

    // Запускаем обработку
    results := processWithCancellation(ctx, 5)

    // Через 500ms отменяем все задачи
    go func() {
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Отменяем все задачи!")
        cancel()
    }()

    // Собираем результаты
    for result := range results {
        if result.Error != nil {
            fmt.Printf("Задача %d: отменена после %d/%d шагов\n",
                result.TaskID, result.CompletedSteps, result.TotalSteps)
        } else {
            fmt.Printf("Задача %d: завершена успешно за %v\n",
                result.TaskID, result.Duration)
        }
    }
}
```

## Ожидаемый вывод (пример)

```
=== Тест 1: Задача без отмены ===
Задача завершена: 5/5 шагов за 500ms

=== Тест 2: Отмена через 250ms ===
Отменяем все задачи!
Задача 0: отменена после 2/5 шагов (context canceled)
Задача 1: отменена после 2/5 шагов (context canceled)
Задача 2: отменена после 2/5 шагов (context canceled)

=== Тест 3: Поиск первого результата ===
Найдено в источнике 'source-2': "Result from source-2" за 150ms
```

## Подсказки

- В `longRunningTask` проверяй контекст на каждой итерации цикла
- Используй `time.NewTicker` для эмуляции шагов с периодической проверкой
- В `searchFirst` создай дочерний контекст с `context.WithCancel` и отмени его после получения первого результата
- Не забывай вызывать `cancel()` через `defer`

## Критерии выполнения

- [ ] `longRunningTask` корректно завершается при отмене контекста
- [ ] `longRunningTask` возвращает правильное количество выполненных шагов
- [ ] `processWithCancellation` запускает задачи параллельно
- [ ] `processWithCancellation` корректно закрывает канал результатов
- [ ] `searchFirst` возвращает первый результат и отменяет остальные поиски
- [ ] Программа проходит проверку `go run -race`
