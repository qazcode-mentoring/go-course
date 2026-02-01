# ДЗ 18.3: JSON REST API

## Цель

Научиться создавать REST API, которое принимает и возвращает JSON. Освоить работу с `encoding/json`, установку правильных заголовков и HTTP-статусов.

## Что нужно сделать

Создать REST API для управления книгами:

1. `GET /api/books` — список всех книг (JSON массив)
2. `GET /api/books/{id}` — получение книги по ID (JSON объект)
3. `POST /api/books` — создание новой книги (принимает JSON, возвращает JSON)
4. `DELETE /api/books/{id}` — удаление книги по ID

## Структуры данных

```go
// Book представляет книгу
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
    Error string `json:"error"`
}

// SuccessResponse — ответ об успешной операции
type SuccessResponse struct {
    Message string `json:"message"`
}
```

## Требования к endpoints

### GET /api/books
- Возвращает JSON массив всех книг
- Content-Type: application/json
- Статус: 200 OK

Пример ответа:
```json
[
    {"id": 1, "title": "1984", "author": "George Orwell", "year": 1949},
    {"id": 2, "title": "Brave New World", "author": "Aldous Huxley", "year": 1932}
]
```

### GET /api/books/{id}
- Возвращает JSON объект книги
- Если книга не найдена — 404 Not Found с JSON ошибкой
- Content-Type: application/json

Пример успешного ответа (200 OK):
```json
{"id": 1, "title": "1984", "author": "George Orwell", "year": 1949}
```

Пример ответа с ошибкой (404 Not Found):
```json
{"error": "book not found"}
```

### POST /api/books
- Принимает JSON объект книги (без id)
- Создаёт книгу с автоматическим ID
- Возвращает созданную книгу с ID
- Статус: 201 Created
- Если JSON невалидный — 400 Bad Request с описанием ошибки
- Если title пустой — 400 Bad Request

Пример запроса:
```json
{"title": "The Martian", "author": "Andy Weir", "year": 2011}
```

Пример ответа (201 Created):
```json
{"id": 3, "title": "The Martian", "author": "Andy Weir", "year": 2011}
```

### DELETE /api/books/{id}
- Удаляет книгу по ID
- Если книга не найдена — 404 Not Found
- Статус: 200 OK с сообщением об успехе

Пример успешного ответа (200 OK):
```json
{"message": "book deleted successfully"}
```

## Вспомогательные функции

```go
// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, status int, data any)

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, status int, message string)

// parseJSON читает JSON из тела запроса
func parseJSON(r *http.Request, v any) error
```

## Примеры использования

```bash
# Список книг
$ curl http://localhost:8080/api/books
[{"id":1,"title":"1984","author":"George Orwell","year":1949}]

# Получение книги
$ curl http://localhost:8080/api/books/1
{"id":1,"title":"1984","author":"George Orwell","year":1949}

# Книга не найдена
$ curl -i http://localhost:8080/api/books/999
HTTP/1.1 404 Not Found
Content-Type: application/json
{"error":"book not found"}

# Создание книги
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"title":"Dune","author":"Frank Herbert","year":1965}' \
    http://localhost:8080/api/books
{"id":2,"title":"Dune","author":"Frank Herbert","year":1965}

# Невалидный JSON
$ curl -X POST -H "Content-Type: application/json" \
    -d 'not json' \
    http://localhost:8080/api/books
{"error":"invalid JSON"}

# Удаление книги
$ curl -X DELETE http://localhost:8080/api/books/1
{"message":"book deleted successfully"}
```

## Подсказки

- Для сериализации: `json.NewEncoder(w).Encode(data)`
- Для десериализации: `json.NewDecoder(r.Body).Decode(&data)`
- Не забудь `defer r.Body.Close()` после чтения тела
- Всегда устанавливай `Content-Type: application/json` ДО записи тела
- Используй `w.WriteHeader(status)` для установки статус-кода ДО записи тела
- Для поиска по slice используй цикл или функцию `slices.IndexFunc`

## Порядок вызовов при формировании ответа

```go
// ПРАВИЛЬНЫЙ порядок:
w.Header().Set("Content-Type", "application/json")  // 1. Заголовки
w.WriteHeader(http.StatusCreated)                    // 2. Статус-код
json.NewEncoder(w).Encode(data)                      // 3. Тело

// НЕПРАВИЛЬНО — заголовки после Write игнорируются:
json.NewEncoder(w).Encode(data)
w.Header().Set("Content-Type", "application/json")  // Уже поздно!
```

## Критерии выполнения

- [ ] GET /api/books возвращает JSON массив книг
- [ ] GET /api/books/{id} возвращает книгу или 404 с JSON ошибкой
- [ ] POST /api/books создаёт книгу, возвращает 201 с JSON
- [ ] POST /api/books с невалидным JSON возвращает 400
- [ ] DELETE /api/books/{id} удаляет книгу или возвращает 404
- [ ] Все ответы имеют Content-Type: application/json
- [ ] Код компилируется без ошибок
