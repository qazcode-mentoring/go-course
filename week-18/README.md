# Неделя 18: HTTP-серверы в Go

## Теория

Пакет `net/http` — это мощная стандартная библиотека Go для создания HTTP-серверов и клиентов. В Go 1.22 появились значительные улучшения в маршрутизации, которые делают стандартную библиотеку ещё более удобной.

### Основные концепции

#### http.Handler — интерфейс обработчика

Центральный интерфейс в пакете `net/http`:

```go
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

Любой тип, реализующий этот интерфейс, может обрабатывать HTTP-запросы.

#### http.HandlerFunc — функциональный адаптер

Позволяет использовать обычные функции как обработчики:

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP вызывает f(w, r)
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}
```

Пример использования:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

// Можно использовать напрямую
http.Handle("/hello", http.HandlerFunc(helloHandler))

// Или через удобную обёртку
http.HandleFunc("/hello", helloHandler)
```

#### http.ServeMux — маршрутизатор (роутер)

ServeMux сопоставляет URL-паттерны с обработчиками:

```go
mux := http.NewServeMux()
mux.HandleFunc("/", homeHandler)
mux.HandleFunc("/api/users", usersHandler)
```

### Go 1.22+: Улучшенная маршрутизация

В Go 1.22 ServeMux получил поддержку HTTP-методов и параметров в путях:

#### Паттерны с методами

```go
mux := http.NewServeMux()

// Только GET запросы
mux.HandleFunc("GET /users", listUsersHandler)

// Только POST запросы
mux.HandleFunc("POST /users", createUserHandler)

// Только DELETE запросы
mux.HandleFunc("DELETE /users/{id}", deleteUserHandler)
```

#### Параметры в путях (wildcards)

```go
mux := http.NewServeMux()

// {id} — параметр пути
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // получение значения параметра
    fmt.Fprintf(w, "User ID: %s", id)
})

// {path...} — захват оставшегося пути
mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
    path := r.PathValue("path")
    fmt.Fprintf(w, "File path: %s", path)
})
```

#### Приоритеты маршрутов

Более специфичные паттерны имеют приоритет:

```go
mux.HandleFunc("GET /users/{id}", getUser)       // /users/123
mux.HandleFunc("GET /users/me", getCurrentUser)  // /users/me (приоритет!)
```

### http.ResponseWriter

Интерфейс для формирования HTTP-ответа:

```go
type ResponseWriter interface {
    Header() http.Header       // доступ к заголовкам
    Write([]byte) (int, error) // запись тела ответа
    WriteHeader(statusCode int) // установка статус-кода
}
```

Примеры использования:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Установка заголовков (до WriteHeader или Write!)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "value")

    // Статус код (по умолчанию 200)
    w.WriteHeader(http.StatusCreated) // 201

    // Запись тела ответа
    w.Write([]byte(`{"status":"ok"}`))
}
```

### http.Request

Структура входящего запроса:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Метод запроса
    method := r.Method // "GET", "POST", etc.

    // URL и путь
    url := r.URL.Path      // "/users/123"
    query := r.URL.Query() // map[string][]string

    // Параметры пути (Go 1.22+)
    id := r.PathValue("id")

    // Заголовки
    auth := r.Header.Get("Authorization")

    // Тело запроса
    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()

    // Query параметры
    name := r.URL.Query().Get("name")

    // Форма (POST form data)
    r.ParseForm()
    field := r.FormValue("field")
}
```

### http.Server — настраиваемый сервер

Для production-сред используйте структуру Server:

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}

// Запуск сервера
log.Fatal(server.ListenAndServe())
```

### Запуск сервера

#### Простой способ (для разработки)

```go
// Использует DefaultServeMux
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)
```

#### Рекомендуемый способ

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /", handler)

server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}

log.Printf("Server starting on %s", server.Addr)
log.Fatal(server.ListenAndServe())
```

### Работа с JSON

```go
import "encoding/json"

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Отправка JSON
func sendJSON(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// Чтение JSON из запроса
func parseJSON(r *http.Request, v any) error {
    return json.NewDecoder(r.Body).Decode(v)
}

func handler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := parseJSON(r, &user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendJSON(w, user)
}
```

### Middleware — промежуточные обработчики

Middleware — это паттерн для добавления общей логики (логирование, аутентификация, CORS):

```go
// Middleware — функция, оборачивающая Handler
type Middleware func(http.Handler) http.Handler

// Логирование запросов
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Вызов следующего обработчика
        next.ServeHTTP(w, r)

        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// Применение middleware
mux := http.NewServeMux()
mux.HandleFunc("GET /", handler)

// Оборачиваем весь mux
wrappedMux := LoggingMiddleware(mux)

http.ListenAndServe(":8080", wrappedMux)
```

### Цепочка middleware

```go
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    // Применяем в обратном порядке
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

// Использование
handler := Chain(mux,
    LoggingMiddleware,
    RecoveryMiddleware,
    CORSMiddleware,
)
```

### Паттерн ResponseWriter Wrapper

Для middleware часто нужно перехватывать статус код:

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        recorder := &statusRecorder{ResponseWriter: w, status: 200}

        next.ServeHTTP(recorder, r)

        log.Printf("%d %s %s", recorder.status, r.Method, r.URL.Path)
    })
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 18.1 | [simple-server](./homework/18.1-simple-server/) | Hello World HTTP-сервер, основы net/http |
| 18.2 | [handlers](./homework/18.2-handlers/) | Несколько endpoints, методы GET/POST, параметры пути |
| 18.3 | [json-api](./homework/18.3-json-api/) | REST API с JSON: создание и получение ресурсов |
| 18.4 | [middleware](./homework/18.4-middleware/) | Logging middleware, цепочка обработки, recovery |

## Вопросы для самопроверки

Создай файл `answers-18.txt` и напиши ответы на вопросы:

1. В чём разница между `http.Handle` и `http.HandleFunc`? Когда использовать каждый?

2. Почему важно устанавливать заголовки (`w.Header().Set(...)`) ДО вызова `w.WriteHeader()` или `w.Write()`? Что произойдёт, если сделать наоборот?

3. Объясни новый синтаксис паттернов в Go 1.22+ (например, `"GET /users/{id}"`). Как получить значение `{id}` внутри обработчика?

4. Что такое middleware и зачем он нужен? Приведи примеры задач, которые решаются через middleware.

5. Почему для production-сервера рекомендуется создавать `http.Server` с таймаутами вместо использования `http.ListenAndServe(":8080", nil)`?

## Дополнительные материалы

- [Go Doc: net/http](https://pkg.go.dev/net/http)
- [Go 1.22 Release Notes: Enhanced routing patterns](https://go.dev/doc/go1.22#enhanced_routing_patterns)
- [Go by Example: HTTP Servers](https://gobyexample.com/http-servers)
- [Go by Example: HTTP Clients](https://gobyexample.com/http-clients)
- [Go Blog: The Go Blog](https://go.dev/blog/)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)
- [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
