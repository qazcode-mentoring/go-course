# ДЗ 19.1: CRUD API для пользователей

## Цель

Научиться создавать полноценный CRUD (Create, Read, Update, Delete) REST API для управления ресурсами. Освоить работу с HTTP-методами, статус-кодами и JSON-сериализацией.

## Что нужно сделать

Реализовать REST API для управления пользователями с in-memory хранилищем:

| Метод | Endpoint | Описание | Статус успеха |
|-------|----------|----------|---------------|
| GET | /api/users | Список всех пользователей | 200 OK |
| GET | /api/users/{id} | Получить пользователя по ID | 200 OK / 404 Not Found |
| POST | /api/users | Создать пользователя | 201 Created |
| PUT | /api/users/{id} | Полностью обновить пользователя | 200 OK / 404 Not Found |
| DELETE | /api/users/{id} | Удалить пользователя | 204 No Content / 404 Not Found |

## Структуры данных

```go
// User представляет пользователя
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest — запрос на создание пользователя
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// UpdateUserRequest — запрос на обновление пользователя
type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
    Error string `json:"error"`
}
```

## Требования к endpoints

### GET /api/users

- Возвращает JSON массив всех пользователей
- Если пользователей нет — возвращает пустой массив `[]`
- Статус: 200 OK

Пример ответа:
```json
[
    {"id": 1, "name": "Alice", "email": "alice@example.com", "age": 25, "created_at": "2024-01-15T10:30:00Z"},
    {"id": 2, "name": "Bob", "email": "bob@example.com", "age": 30, "created_at": "2024-01-15T11:00:00Z"}
]
```

### GET /api/users/{id}

- Возвращает JSON объект пользователя
- Если пользователь не найден — 404 Not Found с JSON ошибкой

Пример успешного ответа (200 OK):
```json
{"id": 1, "name": "Alice", "email": "alice@example.com", "age": 25, "created_at": "2024-01-15T10:30:00Z"}
```

Пример ошибки (404 Not Found):
```json
{"error": "user not found"}
```

### POST /api/users

- Принимает JSON объект CreateUserRequest
- Создаёт пользователя с автоматическим ID и временем создания
- Возвращает созданного пользователя
- Статус: 201 Created
- Если JSON невалидный — 400 Bad Request
- Если name пустой — 400 Bad Request

Пример запроса:
```json
{"name": "Charlie", "email": "charlie@example.com", "age": 28}
```

Пример ответа (201 Created):
```json
{"id": 3, "name": "Charlie", "email": "charlie@example.com", "age": 28, "created_at": "2024-01-15T12:00:00Z"}
```

### PUT /api/users/{id}

- Принимает JSON объект UpdateUserRequest
- Полностью заменяет данные пользователя (кроме ID и CreatedAt)
- Если пользователь не найден — 404 Not Found
- Статус: 200 OK

Пример запроса:
```json
{"name": "Alice Updated", "email": "alice.new@example.com", "age": 26}
```

Пример ответа (200 OK):
```json
{"id": 1, "name": "Alice Updated", "email": "alice.new@example.com", "age": 26, "created_at": "2024-01-15T10:30:00Z"}
```

### DELETE /api/users/{id}

- Удаляет пользователя по ID
- Если пользователь не найден — 404 Not Found
- Статус: 204 No Content (без тела ответа)

## Вспомогательные функции

```go
// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, status int, data any)

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, status int, message string)

// parseJSON читает JSON из тела запроса
func parseJSON(r *http.Request, v any) error
```

## Примеры использования (curl)

```bash
# Список пользователей (изначально с тестовыми данными)
$ curl http://localhost:8080/api/users
[{"id":1,"name":"Alice","email":"alice@example.com","age":25,"created_at":"..."},{"id":2,"name":"Bob","email":"bob@example.com","age":30,"created_at":"..."}]

# Получить пользователя
$ curl http://localhost:8080/api/users/1
{"id":1,"name":"Alice","email":"alice@example.com","age":25,"created_at":"..."}

# Пользователь не найден
$ curl -i http://localhost:8080/api/users/999
HTTP/1.1 404 Not Found
Content-Type: application/json
{"error":"user not found"}

# Создать пользователя
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"name":"Charlie","email":"charlie@example.com","age":28}' \
    http://localhost:8080/api/users
HTTP/1.1 201 Created
{"id":3,"name":"Charlie","email":"charlie@example.com","age":28,"created_at":"..."}

# Создать с невалидным JSON
$ curl -X POST -H "Content-Type: application/json" \
    -d 'not json' \
    http://localhost:8080/api/users
{"error":"invalid request body"}

# Создать без имени
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"email":"test@example.com"}' \
    http://localhost:8080/api/users
{"error":"name is required"}

# Обновить пользователя
$ curl -X PUT -H "Content-Type: application/json" \
    -d '{"name":"Alice Smith","email":"alice.smith@example.com","age":26}' \
    http://localhost:8080/api/users/1
{"id":1,"name":"Alice Smith","email":"alice.smith@example.com","age":26,"created_at":"..."}

# Удалить пользователя
$ curl -i -X DELETE http://localhost:8080/api/users/2
HTTP/1.1 204 No Content

# Проверить, что пользователь удалён
$ curl -i http://localhost:8080/api/users/2
HTTP/1.1 404 Not Found
{"error":"user not found"}
```

## Подсказки

- Для получения ID из пути: `r.PathValue("id")`
- Для конвертации строки в int: `strconv.Atoi(idStr)`
- Для текущего времени: `time.Now()`
- Порядок записи ответа: Header -> WriteHeader -> Body
- При DELETE с 204 No Content не записывай тело ответа
- Используй `sync.RWMutex` для потокобезопасного доступа к map

## Критерии выполнения

- [ ] GET /api/users возвращает JSON массив пользователей
- [ ] GET /api/users/{id} возвращает пользователя или 404
- [ ] POST /api/users создаёт пользователя, возвращает 201
- [ ] POST с пустым name возвращает 400
- [ ] PUT /api/users/{id} обновляет пользователя или возвращает 404
- [ ] DELETE /api/users/{id} удаляет пользователя, возвращает 204
- [ ] DELETE несуществующего возвращает 404
- [ ] Все JSON ответы имеют Content-Type: application/json
- [ ] ID генерируется автоматически при создании
- [ ] CreatedAt устанавливается автоматически при создании
- [ ] Код компилируется без ошибок
