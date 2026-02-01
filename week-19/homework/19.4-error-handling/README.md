# ДЗ 19.4: Единообразная обработка ошибок в API

## Цель

Научиться создавать систему единообразной обработки ошибок в REST API. Освоить паттерн типизированных ошибок (AppError), централизованную обработку и формирование консистентных ответов.

## Что нужно сделать

Реализовать REST API для управления задачами (tasks) с централизованной обработкой ошибок:

1. **AppError** — типизированная ошибка с HTTP-контекстом
2. **Конструкторы ошибок** — NotFound, BadRequest, ValidationError, InternalError
3. **WrapHandler** — обёртка для автоматической обработки ошибок
4. **Валидация запросов** — с детальными сообщениями об ошибках

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | /api/tasks | Список всех задач |
| GET | /api/tasks/{id} | Получить задачу по ID |
| POST | /api/tasks | Создать задачу |
| PUT | /api/tasks/{id} | Обновить задачу |
| DELETE | /api/tasks/{id} | Удалить задачу |

## Структура ошибок

```go
// AppError — типизированная ошибка приложения
type AppError struct {
    Code    string            `json:"code"`              // код ошибки (VALIDATION_ERROR, NOT_FOUND, etc.)
    Message string            `json:"message"`           // сообщение для пользователя
    Status  int               `json:"-"`                 // HTTP статус-код (не сериализуется)
    Details map[string]string `json:"details,omitempty"` // детали (например, ошибки полей)
    Err     error             `json:"-"`                 // оригинальная ошибка (для логов)
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}
```

## Коды ошибок

| Код | HTTP Status | Описание |
|-----|-------------|----------|
| NOT_FOUND | 404 | Ресурс не найден |
| BAD_REQUEST | 400 | Невалидный запрос (неправильный JSON, ID) |
| VALIDATION_ERROR | 422 | Ошибка валидации полей |
| INTERNAL_ERROR | 500 | Внутренняя ошибка сервера |

## Конструкторы ошибок

```go
// NotFound создаёт ошибку "ресурс не найден"
func NotFound(resource string) *AppError {
    return &AppError{
        Code:    "NOT_FOUND",
        Message: fmt.Sprintf("%s not found", resource),
        Status:  http.StatusNotFound,
    }
}

// BadRequest создаёт ошибку "неверный запрос"
func BadRequest(message string) *AppError {
    return &AppError{
        Code:    "BAD_REQUEST",
        Message: message,
        Status:  http.StatusBadRequest,
    }
}

// ValidationError создаёт ошибку валидации с деталями по полям
func ValidationError(details map[string]string) *AppError {
    return &AppError{
        Code:    "VALIDATION_ERROR",
        Message: "validation failed",
        Status:  http.StatusUnprocessableEntity,
        Details: details,
    }
}

// InternalError создаёт внутреннюю ошибку сервера
func InternalError(err error) *AppError {
    return &AppError{
        Code:    "INTERNAL_ERROR",
        Message: "internal server error",
        Status:  http.StatusInternalServerError,
        Err:     err,
    }
}
```

## Паттерн AppHandler

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

        // Обработка AppError
        var appErr *AppError
        if errors.As(err, &appErr) {
            // Логируем оригинальную ошибку, если есть
            if appErr.Err != nil {
                log.Printf("Error [%s]: %v", appErr.Code, appErr.Err)
            }
            writeJSON(w, appErr.Status, appErr)
            return
        }

        // Неизвестная ошибка
        log.Printf("Unexpected error: %v", err)
        writeJSON(w, http.StatusInternalServerError, &AppError{
            Code:    "INTERNAL_ERROR",
            Message: "internal server error",
        })
    }
}
```

## Структуры данных

```go
// Task представляет задачу
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`    // "pending", "in_progress", "done"
    Priority    int       `json:"priority"`  // 1-5
    CreatedAt   time.Time `json:"created_at"`
}

// CreateTaskRequest — запрос на создание задачи
type CreateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Priority    int    `json:"priority"`
}

// Validate проверяет валидность запроса
func (r CreateTaskRequest) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(r.Title) == "" {
        errors["title"] = "title is required"
    } else if len(r.Title) < 3 {
        errors["title"] = "title must be at least 3 characters"
    } else if len(r.Title) > 100 {
        errors["title"] = "title must be at most 100 characters"
    }

    if r.Priority < 1 || r.Priority > 5 {
        errors["priority"] = "priority must be between 1 and 5"
    }

    return errors
}

// UpdateTaskRequest — запрос на обновление задачи
type UpdateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    Priority    int    `json:"priority"`
}

// Validate проверяет валидность запроса
func (r UpdateTaskRequest) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(r.Title) == "" {
        errors["title"] = "title is required"
    } else if len(r.Title) < 3 {
        errors["title"] = "title must be at least 3 characters"
    }

    validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
    if !validStatuses[r.Status] {
        errors["status"] = "status must be one of: pending, in_progress, done"
    }

    if r.Priority < 1 || r.Priority > 5 {
        errors["priority"] = "priority must be between 1 and 5"
    }

    return errors
}
```

## Примеры ответов с ошибками

### 404 Not Found

```json
{
    "code": "NOT_FOUND",
    "message": "task not found"
}
```

### 400 Bad Request

```json
{
    "code": "BAD_REQUEST",
    "message": "invalid task ID"
}
```

### 422 Validation Error

```json
{
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "details": {
        "title": "title must be at least 3 characters",
        "priority": "priority must be between 1 and 5"
    }
}
```

### 500 Internal Server Error

```json
{
    "code": "INTERNAL_ERROR",
    "message": "internal server error"
}
```

## Примеры использования (curl)

```bash
# Список задач
$ curl http://localhost:8080/api/tasks
[{"id":1,"title":"Learn Go","description":"Complete Go course","status":"in_progress","priority":5,"created_at":"..."}]

# Получить задачу
$ curl http://localhost:8080/api/tasks/1
{"id":1,"title":"Learn Go","description":"Complete Go course","status":"in_progress","priority":5,"created_at":"..."}

# Задача не найдена
$ curl -i http://localhost:8080/api/tasks/999
HTTP/1.1 404 Not Found
{"code":"NOT_FOUND","message":"task not found"}

# Невалидный ID
$ curl -i http://localhost:8080/api/tasks/abc
HTTP/1.1 400 Bad Request
{"code":"BAD_REQUEST","message":"invalid task ID"}

# Создание задачи
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"title":"New Task","description":"Task description","priority":3}' \
    http://localhost:8080/api/tasks
HTTP/1.1 201 Created
{"id":3,"title":"New Task","description":"Task description","status":"pending","priority":3,"created_at":"..."}

# Ошибка валидации при создании
$ curl -i -X POST -H "Content-Type: application/json" \
    -d '{"title":"AB","priority":10}' \
    http://localhost:8080/api/tasks
HTTP/1.1 422 Unprocessable Entity
{"code":"VALIDATION_ERROR","message":"validation failed","details":{"title":"title must be at least 3 characters","priority":"priority must be between 1 and 5"}}

# Невалидный JSON
$ curl -i -X POST -H "Content-Type: application/json" \
    -d 'not json' \
    http://localhost:8080/api/tasks
HTTP/1.1 400 Bad Request
{"code":"BAD_REQUEST","message":"invalid request body"}

# Обновление задачи
$ curl -X PUT -H "Content-Type: application/json" \
    -d '{"title":"Updated Task","description":"Updated description","status":"done","priority":4}' \
    http://localhost:8080/api/tasks/1
{"id":1,"title":"Updated Task","description":"Updated description","status":"done","priority":4,"created_at":"..."}

# Ошибка валидации статуса при обновлении
$ curl -i -X PUT -H "Content-Type: application/json" \
    -d '{"title":"Task","status":"invalid","priority":3}' \
    http://localhost:8080/api/tasks/1
HTTP/1.1 422 Unprocessable Entity
{"code":"VALIDATION_ERROR","message":"validation failed","details":{"status":"status must be one of: pending, in_progress, done"}}

# Удаление задачи
$ curl -i -X DELETE http://localhost:8080/api/tasks/1
HTTP/1.1 204 No Content

# Удаление несуществующей задачи
$ curl -i -X DELETE http://localhost:8080/api/tasks/999
HTTP/1.1 404 Not Found
{"code":"NOT_FOUND","message":"task not found"}
```

## Пример обработчика с AppHandler

```go
func getTask(w http.ResponseWriter, r *http.Request) error {
    // Парсинг ID
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return BadRequest("invalid task ID")
    }

    // Поиск задачи
    task, exists := storage.Get(id)
    if !exists {
        return NotFound("task")
    }

    // Успешный ответ
    writeJSON(w, http.StatusOK, task)
    return nil
}

func createTask(w http.ResponseWriter, r *http.Request) error {
    var req CreateTaskRequest

    // Парсинг JSON
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return BadRequest("invalid request body")
    }

    // Валидация
    if errors := req.Validate(); len(errors) > 0 {
        return ValidationError(errors)
    }

    // Создание
    task := storage.Create(req)

    writeJSON(w, http.StatusCreated, task)
    return nil
}

// Регистрация обработчиков
mux.HandleFunc("GET /api/tasks/{id}", WrapHandler(getTask))
mux.HandleFunc("POST /api/tasks", WrapHandler(createTask))
```

## Подсказки

- Используй `errors.As(err, &appErr)` для проверки типа ошибки
- Не забывай логировать оригинальные ошибки для отладки
- Валидация должна собирать ВСЕ ошибки, а не останавливаться на первой
- Status код не должен попадать в JSON (используй `json:"-"`)
- При удалении возвращай 204 No Content (без тела)

## Критерии выполнения

- [ ] Реализована структура AppError с полями code, message, status, details
- [ ] Реализованы конструкторы: NotFound, BadRequest, ValidationError, InternalError
- [ ] WrapHandler корректно обрабатывает AppError и неизвестные ошибки
- [ ] GET /api/tasks/{id} с несуществующим ID возвращает 404 с кодом NOT_FOUND
- [ ] GET /api/tasks/{id} с невалидным ID (abc) возвращает 400 с кодом BAD_REQUEST
- [ ] POST /api/tasks с невалидным JSON возвращает 400
- [ ] POST /api/tasks с ошибками валидации возвращает 422 с details по полям
- [ ] Валидация проверяет все поля и возвращает все ошибки сразу
- [ ] PUT /api/tasks/{id} валидирует status (pending/in_progress/done)
- [ ] Оригинальные ошибки логируются, но не отправляются клиенту
- [ ] Все ответы имеют единый формат
- [ ] Код компилируется без ошибок
