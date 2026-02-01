# ДЗ 18.4: Middleware

## Цель

Научиться создавать middleware (промежуточные обработчики) для HTTP-серверов. Освоить паттерн цепочки middleware, логирование запросов и обработку паник (recovery).

## Что нужно сделать

1. Реализовать **LoggingMiddleware** — логирует метод, путь, статус-код и время выполнения
2. Реализовать **RecoveryMiddleware** — перехватывает паники и возвращает 500 ошибку
3. Реализовать функцию **Chain** — объединяет middleware в цепочку
4. Применить middleware к тестовым обработчикам

## Концепция Middleware

Middleware — это функция, которая оборачивает handler и добавляет дополнительную логику:

```go
type Middleware func(http.Handler) http.Handler
```

## Сигнатуры функций

```go
// LoggingMiddleware логирует информацию о каждом запросе
// Формат: "METHOD /path STATUS TIME"
// Пример: "GET /api/users 200 1.234ms"
func LoggingMiddleware(next http.Handler) http.Handler

// RecoveryMiddleware перехватывает паники и возвращает 500
// При панике логирует ошибку и возвращает JSON: {"error": "internal server error"}
func RecoveryMiddleware(next http.Handler) http.Handler

// Chain объединяет middleware в цепочку
// Middleware применяются в порядке передачи:
// Chain(handler, A, B, C) -> A(B(C(handler)))
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler
```

## StatusRecorder — обёртка для ResponseWriter

Для логирования статус-кода нужна обёртка:

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}
```

## Тестовые обработчики

```go
// GET / — обычный обработчик
func homeHandler(w http.ResponseWriter, r *http.Request)

// GET /api/data — возвращает JSON
func dataHandler(w http.ResponseWriter, r *http.Request)

// GET /panic — вызывает панику (для тестирования recovery)
func panicHandler(w http.ResponseWriter, r *http.Request)

// GET /slow — медленный обработчик (sleep 100ms)
func slowHandler(w http.ResponseWriter, r *http.Request)
```

## Требования

### LoggingMiddleware
- Логирует в stdout перед и/или после обработки запроса
- Формат лога: `METHOD /path STATUS DURATION`
- Пример: `GET /api/data 200 2.345ms`
- Используй `time.Since(start)` для измерения времени
- Статус-код получай через statusRecorder

### RecoveryMiddleware
- Использует `defer` и `recover()` для перехвата паники
- При панике:
  - Логирует ошибку: `Panic recovered: ERROR`
  - Устанавливает статус 500
  - Возвращает JSON: `{"error": "internal server error"}`
- Если паники нет — просто передаёт управление следующему handler

### Chain
- Принимает handler и список middleware
- Возвращает handler, обёрнутый всеми middleware
- Middleware применяются в порядке передачи (первый — внешний)

## Примеры использования

```bash
# Обычный запрос
$ curl http://localhost:8080/
# Лог сервера: GET / 200 156.084µs
Hello, World!

# Запрос к API
$ curl http://localhost:8080/api/data
# Лог сервера: GET /api/data 200 1.234ms
{"message":"Hello from API","timestamp":"2024-01-15T10:30:00Z"}

# Медленный запрос
$ curl http://localhost:8080/slow
# Лог сервера: GET /slow 200 100.234ms
Done!

# Запрос, вызывающий панику
$ curl http://localhost:8080/panic
# Лог сервера: Panic recovered: test panic
# Лог сервера: GET /panic 500 234.567µs
{"error":"internal server error"}
```

## Порядок применения middleware

```go
// Запрос проходит через middleware в порядке: Logging -> Recovery -> Handler
// Ответ идёт обратно: Handler -> Recovery -> Logging
handler := Chain(mux, LoggingMiddleware, RecoveryMiddleware)

// Эквивалентно:
handler := LoggingMiddleware(RecoveryMiddleware(mux))
```

## Подсказки

- Используй `defer` для выполнения кода после обработки запроса
- `recover()` работает только в defer-функции
- Время начала: `start := time.Now()`
- Длительность: `time.Since(start)`
- Для логирования: `log.Printf("FORMAT", args...)`
- Не забудь вызвать `next.ServeHTTP(w, r)` для передачи управления

## Пример структуры RecoveryMiddleware

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Обработка паники
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

## Критерии выполнения

- [ ] LoggingMiddleware логирует метод, путь, статус-код и время
- [ ] RecoveryMiddleware перехватывает паники и возвращает 500
- [ ] RecoveryMiddleware логирует информацию о панике
- [ ] Chain корректно объединяет middleware
- [ ] /panic не крашит сервер, а возвращает JSON ошибку
- [ ] Все обработчики логируются
- [ ] Код компилируется без ошибок
