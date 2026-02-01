# ДЗ 18.2: HTTP Handlers

## Цель

Научиться работать с разными HTTP-методами (GET, POST), параметрами пути и query-параметрами. Использовать новый синтаксис маршрутов Go 1.22+ с указанием метода и параметрами `{param}`.

## Что нужно сделать

Создать HTTP-сервер с несколькими endpoints для управления списком пользователей (in-memory):

1. `GET /users` — список всех пользователей
2. `GET /users/{id}` — получение пользователя по ID
3. `POST /users` — создание нового пользователя
4. `GET /search?name=...` — поиск пользователя по имени

## Структура данных

```go
type User struct {
    ID   int
    Name string
    Age  int
}

// Хранилище пользователей (in-memory)
var users = []User{
    {ID: 1, Name: "Alice", Age: 25},
    {ID: 2, Name: "Bob", Age: 30},
    {ID: 3, Name: "Charlie", Age: 35},
}
var nextID = 4
```

## Сигнатуры функций

```go
// GET /users — возвращает список всех пользователей
func listUsersHandler(w http.ResponseWriter, r *http.Request)

// GET /users/{id} — возвращает пользователя по ID
func getUserHandler(w http.ResponseWriter, r *http.Request)

// POST /users — создаёт нового пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request)

// GET /search?name=... — поиск по имени
func searchHandler(w http.ResponseWriter, r *http.Request)
```

## Требования к endpoints

### GET /users
- Возвращает список всех пользователей
- Формат: по одной строке на пользователя
- Пример: `1: Alice (25 лет)`

### GET /users/{id}
- Возвращает информацию о пользователе с указанным ID
- Если пользователь не найден — статус 404 и сообщение "User not found"
- Параметр `{id}` получается через `r.PathValue("id")`

### POST /users
- Создаёт нового пользователя
- Данные передаются в теле запроса как form data: `name=...&age=...`
- Возвращает статус 201 Created и информацию о созданном пользователе
- Если name или age не указаны — статус 400 Bad Request

### GET /search?name=...
- Ищет пользователей, чьё имя содержит указанную подстроку (case-insensitive)
- Query-параметр получается через `r.URL.Query().Get("name")`
- Если параметр name не указан — статус 400 Bad Request
- Возвращает список найденных пользователей

## Примеры использования

```bash
# Список пользователей
$ curl http://localhost:8080/users
1: Alice (25 лет)
2: Bob (30 лет)
3: Charlie (35 лет)

# Получение пользователя по ID
$ curl http://localhost:8080/users/2
ID: 2
Name: Bob
Age: 30

# Пользователь не найден
$ curl -i http://localhost:8080/users/999
HTTP/1.1 404 Not Found
User not found

# Создание пользователя
$ curl -X POST -d "name=Diana&age=28" http://localhost:8080/users
Created user:
ID: 4
Name: Diana
Age: 28

# Поиск по имени
$ curl "http://localhost:8080/search?name=ali"
Found 1 user(s):
1: Alice (25 лет)
```

## Подсказки

- Go 1.22+ позволяет указывать метод в паттерне: `mux.HandleFunc("GET /users", handler)`
- Параметры пути: `r.PathValue("id")` (Go 1.22+)
- Query параметры: `r.URL.Query().Get("name")`
- Form data: `r.ParseForm()` и затем `r.FormValue("name")`
- Конвертация строки в число: `strconv.Atoi(str)`
- Поиск подстроки: `strings.Contains(strings.ToLower(s), strings.ToLower(sub))`

## Регистрация маршрутов (Go 1.22+)

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users", listUsersHandler)
mux.HandleFunc("GET /users/{id}", getUserHandler)
mux.HandleFunc("POST /users", createUserHandler)
mux.HandleFunc("GET /search", searchHandler)
```

## Критерии выполнения

- [ ] GET /users возвращает список пользователей
- [ ] GET /users/{id} возвращает пользователя или 404
- [ ] POST /users создаёт пользователя и возвращает 201
- [ ] POST /users без данных возвращает 400
- [ ] GET /search?name=... ищет пользователей по имени (case-insensitive)
- [ ] GET /search без параметра возвращает 400
- [ ] Используется http.NewServeMux() с паттернами Go 1.22+
- [ ] Код компилируется без ошибок
