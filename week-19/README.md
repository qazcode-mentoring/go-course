# Неделя 19: REST API в Go

## Теория

REST (Representational State Transfer) — это архитектурный стиль для создания веб-сервисов. В этом модуле мы научимся проектировать и реализовывать RESTful API на Go с использованием стандартной библиотеки `net/http` (Go 1.22+).

### Основные принципы REST

#### 1. Ресурсы и их идентификация

REST оперирует ресурсами, идентифицируемыми через URL:

```
/users           — коллекция пользователей
/users/123       — конкретный пользователь
/users/123/posts — посты пользователя 123
```

#### 2. HTTP-методы (глаголы)

| Метод | Операция | Описание | Идемпотентность |
|-------|----------|----------|-----------------|
| GET | Read | Получение ресурса | Да |
| POST | Create | Создание нового ресурса | Нет |
| PUT | Update/Replace | Полная замена ресурса | Да |
| PATCH | Update/Modify | Частичное обновление | Нет |
| DELETE | Delete | Удаление ресурса | Да |

```go
mux := http.NewServeMux()

// CRUD для пользователей
mux.HandleFunc("GET /users", listUsers)       // список
mux.HandleFunc("GET /users/{id}", getUser)    // получить одного
mux.HandleFunc("POST /users", createUser)     // создать
mux.HandleFunc("PUT /users/{id}", updateUser) // обновить полностью
mux.HandleFunc("PATCH /users/{id}", patchUser)// обновить частично
mux.HandleFunc("DELETE /users/{id}", deleteUser) // удалить
```

#### 3. HTTP-статус коды

Правильные статус-коды — важная часть REST API:

**2xx — Успех:**
```go
http.StatusOK                  // 200 — успешный GET, PUT, PATCH, DELETE
http.StatusCreated             // 201 — успешный POST (ресурс создан)
http.StatusNoContent           // 204 — успешный DELETE (нет тела ответа)
```

**4xx — Ошибки клиента:**
```go
http.StatusBadRequest          // 400 — невалидные данные
http.StatusUnauthorized        // 401 — не авторизован
http.StatusForbidden           // 403 — доступ запрещён
http.StatusNotFound            // 404 — ресурс не найден
http.StatusMethodNotAllowed    // 405 — метод не поддерживается
http.StatusConflict            // 409 — конфликт (например, дубликат)
http.StatusUnprocessableEntity // 422 — ошибка валидации
```

**5xx — Ошибки сервера:**
```go
http.StatusInternalServerError // 500 — внутренняя ошибка
http.StatusServiceUnavailable  // 503 — сервис недоступен
```

### Параметры запроса

#### Параметры пути (Path Parameters)

Используются для идентификации конкретного ресурса:

```go
// Go 1.22+ — встроенная поддержка
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // "123"
    // ...
})

// Вложенные ресурсы
mux.HandleFunc("GET /users/{userId}/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
    userId := r.PathValue("userId")
    postId := r.PathValue("postId")
    // ...
})
```

#### Query Parameters

Используются для фильтрации, сортировки и пагинации:

```go
// GET /users?status=active&sort=name&page=2&limit=10
func listUsers(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    // Получение одного значения
    status := query.Get("status")     // "active" или ""
    sort := query.Get("sort")         // "name" или ""

    // Получение всех значений (для повторяющихся параметров)
    // GET /users?tag=go&tag=rest
    tags := query["tag"]              // []string{"go", "rest"}

    // Значение по умолчанию
    page := query.Get("page")
    if page == "" {
        page = "1"
    }

    // Конвертация в число
    pageNum, err := strconv.Atoi(page)
    if err != nil {
        pageNum = 1
    }
}
```

### Структура REST API ответов

#### Успешные ответы

```go
// Один ресурс
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Коллекция с пагинацией
type UsersResponse struct {
    Users      []User `json:"users"`
    Total      int    `json:"total"`
    Page       int    `json:"page"`
    PerPage    int    `json:"per_page"`
    TotalPages int    `json:"total_pages"`
}
```

#### Ответы с ошибками

```go
// Единый формат ошибки
type ErrorResponse struct {
    Error   string `json:"error"`             // краткое сообщение
    Code    string `json:"code,omitempty"`    // код ошибки
    Details any    `json:"details,omitempty"` // дополнительная информация
}

// Пример использования
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Ошибки валидации
type ValidationError struct {
    Error  string            `json:"error"`
    Fields map[string]string `json:"fields"` // поле -> сообщение
}
```

### Валидация входных данных

```go
// CreateUserRequest — запрос на создание пользователя
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// Validate проверяет валидность данных
func (r CreateUserRequest) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(r.Name) == "" {
        errors["name"] = "name is required"
    } else if len(r.Name) < 2 {
        errors["name"] = "name must be at least 2 characters"
    }

    if strings.TrimSpace(r.Email) == "" {
        errors["email"] = "email is required"
    } else if !strings.Contains(r.Email, "@") {
        errors["email"] = "invalid email format"
    }

    if r.Age < 0 || r.Age > 150 {
        errors["age"] = "age must be between 0 and 150"
    }

    return errors
}

// Использование в обработчике
func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    if errors := req.Validate(); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnprocessableEntity)
        json.NewEncoder(w).Encode(ValidationError{
            Error:  "validation failed",
            Fields: errors,
        })
        return
    }

    // Создание пользователя...
}
```

### Паттерн обработки ошибок в API

#### Типизированные ошибки приложения

```go
// AppError — ошибка приложения с HTTP-контекстом
type AppError struct {
    Code    string // код ошибки для клиента
    Message string // сообщение для пользователя
    Status  int    // HTTP статус-код
    Err     error  // оригинальная ошибка (для логов)
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Конструкторы для частых ошибок
func NotFound(resource string) *AppError {
    return &AppError{
        Code:    "NOT_FOUND",
        Message: fmt.Sprintf("%s not found", resource),
        Status:  http.StatusNotFound,
    }
}

func BadRequest(message string) *AppError {
    return &AppError{
        Code:    "BAD_REQUEST",
        Message: message,
        Status:  http.StatusBadRequest,
    }
}

func InternalError(err error) *AppError {
    return &AppError{
        Code:    "INTERNAL_ERROR",
        Message: "internal server error",
        Status:  http.StatusInternalServerError,
        Err:     err,
    }
}
```

#### Обработчик с возвратом ошибок

```go
// AppHandler — обработчик, возвращающий ошибку
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// WrapHandler оборачивает AppHandler в стандартный http.HandlerFunc
func WrapHandler(h AppHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := h(w, r)
        if err == nil {
            return
        }

        // Преобразуем ошибку в AppError
        var appErr *AppError
        if errors.As(err, &appErr) {
            writeJSON(w, appErr.Status, ErrorResponse{
                Error: appErr.Message,
                Code:  appErr.Code,
            })
            // Логируем оригинальную ошибку
            if appErr.Err != nil {
                log.Printf("Error: %v", appErr.Err)
            }
            return
        }

        // Неизвестная ошибка
        log.Printf("Unexpected error: %v", err)
        writeJSON(w, http.StatusInternalServerError, ErrorResponse{
            Error: "internal server error",
            Code:  "INTERNAL_ERROR",
        })
    }
}

// Использование
func getUser(w http.ResponseWriter, r *http.Request) error {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        return BadRequest("invalid user ID")
    }

    user, exists := users[id]
    if !exists {
        return NotFound("user")
    }

    writeJSON(w, http.StatusOK, user)
    return nil
}

// Регистрация
mux.HandleFunc("GET /users/{id}", WrapHandler(getUser))
```

### Пагинация

```go
// PaginationParams содержит параметры пагинации
type PaginationParams struct {
    Page    int
    PerPage int
}

// ParsePagination извлекает параметры пагинации из запроса
func ParsePagination(r *http.Request, defaultPerPage, maxPerPage int) PaginationParams {
    query := r.URL.Query()

    page, _ := strconv.Atoi(query.Get("page"))
    if page < 1 {
        page = 1
    }

    perPage, _ := strconv.Atoi(query.Get("per_page"))
    if perPage < 1 {
        perPage = defaultPerPage
    }
    if perPage > maxPerPage {
        perPage = maxPerPage
    }

    return PaginationParams{Page: page, PerPage: perPage}
}

// Paginate возвращает срез данных для указанной страницы
func Paginate[T any](items []T, params PaginationParams) []T {
    start := (params.Page - 1) * params.PerPage
    if start >= len(items) {
        return []T{}
    }

    end := start + params.PerPage
    if end > len(items) {
        end = len(items)
    }

    return items[start:end]
}
```

### Фильтрация

```go
// UserFilter содержит параметры фильтрации пользователей
type UserFilter struct {
    Status string   // фильтр по статусу
    Role   string   // фильтр по роли
    Search string   // поиск по имени/email
    Tags   []string // фильтр по тегам
}

// ParseUserFilter извлекает фильтры из запроса
func ParseUserFilter(r *http.Request) UserFilter {
    query := r.URL.Query()
    return UserFilter{
        Status: query.Get("status"),
        Role:   query.Get("role"),
        Search: query.Get("search"),
        Tags:   query["tag"], // множественные значения
    }
}

// FilterUsers применяет фильтры к списку пользователей
func FilterUsers(users []User, filter UserFilter) []User {
    result := make([]User, 0)

    for _, u := range users {
        if filter.Status != "" && u.Status != filter.Status {
            continue
        }
        if filter.Role != "" && u.Role != filter.Role {
            continue
        }
        if filter.Search != "" {
            search := strings.ToLower(filter.Search)
            if !strings.Contains(strings.ToLower(u.Name), search) &&
               !strings.Contains(strings.ToLower(u.Email), search) {
                continue
            }
        }
        result = append(result, u)
    }

    return result
}
```

### In-Memory Storage

Для учебных целей используем простое хранилище в памяти:

```go
// Storage — потокобезопасное хранилище
type Storage struct {
    mu     sync.RWMutex
    users  map[int]User
    nextID int
}

// NewStorage создаёт новое хранилище
func NewStorage() *Storage {
    return &Storage{
        users:  make(map[int]User),
        nextID: 1,
    }
}

// Create добавляет нового пользователя
func (s *Storage) Create(user User) User {
    s.mu.Lock()
    defer s.mu.Unlock()

    user.ID = s.nextID
    s.nextID++
    s.users[user.ID] = user
    return user
}

// Get возвращает пользователя по ID
func (s *Storage) Get(id int) (User, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, exists := s.users[id]
    return user, exists
}

// GetAll возвращает всех пользователей
func (s *Storage) GetAll() []User {
    s.mu.RLock()
    defer s.mu.RUnlock()

    users := make([]User, 0, len(s.users))
    for _, u := range s.users {
        users = append(users, u)
    }
    return users
}

// Update обновляет пользователя
func (s *Storage) Update(id int, user User) (User, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[id]; !exists {
        return User{}, false
    }

    user.ID = id
    s.users[id] = user
    return user, true
}

// Delete удаляет пользователя
func (s *Storage) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[id]; !exists {
        return false
    }

    delete(s.users, id)
    return true
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 19.1 | [crud-api](./homework/19.1-crud-api/) | Полный CRUD API для пользователей с in-memory storage |
| 19.2 | [path-params](./homework/19.2-path-params/) | Вложенные ресурсы и параметры пути (/users/{id}/posts/{postId}) |
| 19.3 | [query-params](./homework/19.3-query-params/) | Query parameters: фильтрация, сортировка, пагинация |
| 19.4 | [error-handling](./homework/19.4-error-handling/) | Единообразная обработка ошибок с типизированными AppError |

## Вопросы для самопроверки

Создай файл `answers-19.txt` и напиши ответы на вопросы:

1. Какой HTTP-статус код следует вернуть при успешном создании ресурса через POST? А при успешном удалении через DELETE, если не возвращаем тело ответа?

2. Чем отличаются параметры пути (path parameters) от query parameters? Приведи примеры, когда использовать каждый тип.

3. Что такое идемпотентность HTTP-методов? Какие методы идемпотентны, какие нет? Почему это важно?

4. Почему в REST API рекомендуется использовать единый формат ошибок? Какие поля должен содержать JSON-ответ с ошибкой?

5. Как реализовать пагинацию для списка ресурсов? Какие параметры обычно используются и какая информация возвращается в ответе?

## Дополнительные материалы

- [Go Doc: net/http](https://pkg.go.dev/net/http)
- [Go 1.22: Enhanced routing patterns](https://go.dev/doc/go1.22#enhanced_routing_patterns)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [REST API Tutorial](https://restfulapi.net/)
- [Best Practices for REST API Design](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
- [Google API Design Guide](https://cloud.google.com/apis/design)
- [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- [JSON:API Specification](https://jsonapi.org/)
