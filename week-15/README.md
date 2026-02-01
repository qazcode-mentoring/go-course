# Неделя 15: Продвинутая конкурентность (Context, errgroup, Graceful Shutdown)

## Теория

На прошлых неделях мы изучили горутины, каналы, select, WaitGroup и Mutex. Эта неделя посвящена продвинутым паттернам конкурентности, которые используются в production-коде: управление жизненным циклом операций через context, обработка ошибок в группах горутин и корректное завершение сервисов.

### Пакет context: управление отменой и таймаутами

`context.Context` — это стандартный механизм Go для:
- Передачи сигналов отмены операций
- Установки дедлайнов и таймаутов
- Передачи request-scoped данных между функциями

#### Создание контекста

```go
// Пустой контекст (корень дерева контекстов)
ctx := context.Background()

// Контекст с отменой
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // ВАЖНО: всегда вызывай cancel!

// Контекст с таймаутом
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Контекст с дедлайном
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Контекст со значениями (используй осторожно!)
ctx := context.WithValue(parentCtx, "requestID", "abc-123")
```

#### Проверка отмены контекста

```go
func doWork(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            // Контекст отменён
            return ctx.Err() // context.Canceled или context.DeadlineExceeded
        default:
            // Продолжаем работу
        }

        // Выполняем часть работы...
    }
}
```

#### Правила работы с context

1. **Всегда передавай context первым параметром**: `func DoSomething(ctx context.Context, arg1 string) error`
2. **Не храни context в структурах**: передавай его явно в методы
3. **Всегда вызывай cancel()**: используй `defer cancel()` сразу после создания
4. **Не передавай nil context**: используй `context.TODO()` если пока не знаешь какой контекст использовать

### errgroup: группы горутин с обработкой ошибок

Пакет `golang.org/x/sync/errgroup` упрощает работу с группами горутин, особенно когда важна обработка ошибок.

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]string, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]string, len(urls))

    for i, url := range urls {
        i, url := i, url // Захват переменных (Go <1.22)
        g.Go(func() error {
            data, err := fetch(ctx, url)
            if err != nil {
                return err
            }
            results[i] = data
            return nil
        })
    }

    // Ждём завершения всех горутин
    // Если хотя бы одна вернёт ошибку - ctx будет отменён
    if err := g.Wait(); err != nil {
        return nil, err
    }

    return results, nil
}
```

#### Преимущества errgroup над WaitGroup

| WaitGroup | errgroup |
|-----------|----------|
| Нет обработки ошибок | Собирает первую ошибку |
| Ручное управление Add/Done | Автоматическое |
| Нет связи с context | Интегрирован с context |
| Все горутины работают до конца | При ошибке отменяет остальные |

#### errgroup с ограничением параллелизма

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(10) // Максимум 10 параллельных горутин

for _, task := range tasks {
    task := task
    g.Go(func() error {
        return processTask(ctx, task)
    })
}

if err := g.Wait(); err != nil {
    return err
}
```

### Graceful Shutdown: корректное завершение сервиса

Graceful shutdown — это паттерн корректного завершения работы сервиса:
1. Перестаём принимать новые запросы
2. Дожидаемся завершения текущих операций
3. Освобождаем ресурсы (закрываем соединения с БД, файлы и т.д.)
4. Завершаем процесс

#### Обработка сигналов ОС

```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Создаём контекст, который отменится при получении сигнала
    ctx, stop := signal.NotifyContext(context.Background(),
        syscall.SIGINT,  // Ctrl+C
        syscall.SIGTERM, // kill
    )
    defer stop()

    // Запускаем сервер
    go runServer(ctx)

    // Ждём сигнала
    <-ctx.Done()
    fmt.Println("Получен сигнал завершения...")

    // Даём время на завершение операций
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := gracefulShutdown(shutdownCtx); err != nil {
        log.Printf("Ошибка при завершении: %v", err)
    }
}
```

#### Graceful shutdown HTTP-сервера

```go
srv := &http.Server{Addr: ":8080", Handler: mux}

go func() {
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
}()

// Ждём сигнала завершения
<-ctx.Done()

// Корректно останавливаем сервер
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
    log.Printf("HTTP server shutdown error: %v", err)
}
```

### Паттерн: Worker с отменой

```go
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: завершение по отмене контекста\n", id)
            return
        case job, ok := <-jobs:
            if !ok {
                fmt.Printf("Worker %d: канал задач закрыт\n", id)
                return
            }

            // Обрабатываем задачу с учётом контекста
            result, err := processJob(ctx, job)
            if err != nil {
                if ctx.Err() != nil {
                    return // Контекст отменён
                }
                // Другая ошибка
            }

            select {
            case results <- result:
            case <-ctx.Done():
                return
            }
        }
    }
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 15.1 | [context-cancellation](./homework/15.1-context-cancellation/) | Отмена операций через context |
| 15.2 | [context-timeout](./homework/15.2-context-timeout/) | HTTP-подобный запрос с timeout |
| 15.3 | [errgroup-fetch](./homework/15.3-errgroup-fetch/) | Параллельные запросы с errgroup |
| 15.4 | [graceful-shutdown](./homework/15.4-graceful-shutdown/) | Корректное завершение сервиса |

## Вопросы для самопроверки

Создай файл `answers-15.txt` и напиши ответы на вопросы:

1. В чём разница между `context.WithCancel`, `context.WithTimeout` и `context.WithDeadline`? Когда использовать каждый из них?

2. Почему важно всегда вызывать `cancel()` даже если контекст уже завершился? Что произойдёт, если не вызвать?

3. Чем `errgroup` лучше комбинации `sync.WaitGroup` + сбор ошибок вручную? В каких случаях `WaitGroup` всё же предпочтительнее?

4. Опиши последовательность действий при graceful shutdown HTTP-сервера. Почему нельзя просто вызвать `os.Exit(0)`?

5. Как передать request ID через цепочку вызовов функций? Какие есть альтернативы использованию `context.WithValue`?

## Дополнительные материалы

- [Go Blog: Context](https://go.dev/blog/context)
- [Go Blog: Contexts and structs](https://go.dev/blog/context-and-structs)
- [Go Doc: context package](https://pkg.go.dev/context)
- [Go Doc: errgroup package](https://pkg.go.dev/golang.org/x/sync/errgroup)
- [Go by Example: Context](https://gobyexample.com/context)
- [Go by Example: Signals](https://gobyexample.com/signals)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [GopherCon 2019: How to Gracefully Close Channels](https://www.youtube.com/watch?v=Q9nZwBVrDGA)
- [Dave Cheney: Context is for cancelation](https://dave.cheney.net/2017/01/26/context-is-for-cancelation)
