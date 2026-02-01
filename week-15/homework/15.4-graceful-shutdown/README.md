# ДЗ 15.4: Graceful Shutdown

## Цель
Научиться корректно завершать работу сервиса: останавливать приём новых задач, дожидаться завершения текущих операций и освобождать ресурсы.

## Что нужно сделать

Реализовать сервис обработки задач с корректным завершением:

1. **`Worker`** — воркер, который обрабатывает задачи и корректно завершается
2. **`WorkerPool`** — пул воркеров с graceful shutdown
3. **`Service`** — сервис с фоновыми задачами и корректным завершением

## Сигнатуры

```go
// Task представляет задачу для обработки
type Task struct {
    ID       int
    Data     string
    Priority int
}

// Result представляет результат обработки задачи
type Result struct {
    TaskID   int
    Output   string
    Duration time.Duration
    Error    error
}

// Worker обрабатывает задачи из канала
type Worker struct {
    ID int
    // Добавь нужные поля
}

// NewWorker создаёт нового воркера
func NewWorker(id int) *Worker

// Start запускает воркера. Воркер читает задачи из jobs и отправляет результаты в results.
// При отмене контекста воркер завершает текущую задачу и выходит.
func (w *Worker) Start(ctx context.Context, jobs <-chan Task, results chan<- Result)

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
    // Добавь нужные поля
}

// NewWorkerPool создаёт пул из n воркеров
func NewWorkerPool(n int) *WorkerPool

// Start запускает пул воркеров
func (p *WorkerPool) Start(ctx context.Context)

// Submit добавляет задачу в очередь. Возвращает false если пул остановлен.
func (p *WorkerPool) Submit(task Task) bool

// Results возвращает канал с результатами
func (p *WorkerPool) Results() <-chan Result

// Shutdown корректно останавливает пул
// Ждёт завершения всех задач или до истечения таймаута
func (p *WorkerPool) Shutdown(timeout time.Duration) error

// Service представляет сервис с фоновыми задачами
type Service struct {
    // Добавь нужные поля
}

// NewService создаёт новый сервис
func NewService() *Service

// Start запускает сервис и все фоновые задачи
func (s *Service) Start(ctx context.Context) error

// Shutdown корректно останавливает сервис
func (s *Service) Shutdown(ctx context.Context) error
```

## Теория: Graceful Shutdown

### Обработка сигналов ОС

```go
import (
    "os"
    "os/signal"
    "syscall"
)

// Способ 1: signal.NotifyContext (Go 1.16+)
ctx, stop := signal.NotifyContext(context.Background(),
    syscall.SIGINT,  // Ctrl+C
    syscall.SIGTERM, // kill
)
defer stop()

// Способ 2: Канал сигналов
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

select {
case <-sigChan:
    fmt.Println("Получен сигнал завершения")
}
```

### Паттерн Graceful Shutdown

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(),
        syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Запускаем сервис
    svc := NewService()
    go svc.Start(ctx)

    // Ждём сигнала
    <-ctx.Done()
    fmt.Println("Начинаем graceful shutdown...")

    // Даём время на завершение
    shutdownCtx, cancel := context.WithTimeout(
        context.Background(), 30*time.Second)
    defer cancel()

    if err := svc.Shutdown(shutdownCtx); err != nil {
        fmt.Printf("Ошибка при завершении: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Сервис остановлен корректно")
}
```

### Паттерн Worker с корректным завершением

```go
func worker(ctx context.Context, jobs <-chan Job, results chan<- Result) {
    for {
        select {
        case <-ctx.Done():
            // Контекст отменён - выходим
            return

        case job, ok := <-jobs:
            if !ok {
                // Канал закрыт - выходим
                return
            }

            // Обрабатываем задачу
            result := processJob(ctx, job)

            // Отправляем результат (с проверкой контекста)
            select {
            case results <- result:
            case <-ctx.Done():
                return
            }
        }
    }
}
```

## Пример использования

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(),
        syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Создаём пул воркеров
    pool := NewWorkerPool(4)
    pool.Start(ctx)

    // Добавляем задачи
    for i := range 10 {
        pool.Submit(Task{ID: i, Data: fmt.Sprintf("task-%d", i)})
    }

    // Получаем результаты в отдельной горутине
    go func() {
        for result := range pool.Results() {
            fmt.Printf("Задача %d: %s\n", result.TaskID, result.Output)
        }
    }()

    // Ждём Ctrl+C
    <-ctx.Done()
    fmt.Println("\nПолучен сигнал завершения")

    // Корректно останавливаем
    if err := pool.Shutdown(5 * time.Second); err != nil {
        fmt.Printf("Ошибка: %v\n", err)
    } else {
        fmt.Println("Пул остановлен корректно")
    }
}
```

## Ожидаемый вывод (пример)

```
=== Сервис запущен ===
Worker 1: обработка задачи 0
Worker 2: обработка задачи 1
Worker 3: обработка задачи 2
Worker 1: задача 0 завершена
Worker 4: обработка задачи 3
^C
=== Получен сигнал завершения ===
Ожидание завершения текущих задач...
Worker 2: задача 1 завершена
Worker 3: задача 2 завершена (последняя)
Worker 4: задача 3 отменена (контекст)
Освобождение ресурсов...
=== Сервис остановлен корректно ===
```

## Подсказки

- Используй `sync.WaitGroup` для отслеживания активных воркеров
- Закрывай канал задач для сигнала воркерам о завершении
- В `Shutdown` сначала закрой канал задач, потом жди WaitGroup
- Используй `select` с `ctx.Done()` при отправке в каналы
- Для таймаута shutdown используй `context.WithTimeout`

## Критерии выполнения

- [ ] Worker корректно завершается при отмене контекста
- [ ] Worker завершает текущую задачу перед выходом
- [ ] WorkerPool запускает указанное количество воркеров
- [ ] WorkerPool.Submit возвращает false после закрытия пула
- [ ] WorkerPool.Shutdown дожидается завершения задач
- [ ] WorkerPool.Shutdown возвращает ошибку при таймауте
- [ ] Service корректно освобождает ресурсы
- [ ] Программа проходит проверку `go run -race`

## Бонус

Добавь в Service:
- Периодическую фоновую задачу (например, health check каждые 5 секунд)
- Метрики: количество обработанных задач, среднее время обработки
- Логирование с уровнями (info, debug)
