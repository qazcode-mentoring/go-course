# ДЗ 19.2: Параметры пути и вложенные ресурсы

## Цель

Научиться работать с вложенными ресурсами в REST API и множественными параметрами пути. Освоить проектирование иерархических URL-структур.

## Что нужно сделать

Реализовать REST API для управления пользователями и их постами (вложенный ресурс):

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | /api/users | Список всех пользователей |
| GET | /api/users/{userId} | Получить пользователя по ID |
| GET | /api/users/{userId}/posts | Список постов пользователя |
| GET | /api/users/{userId}/posts/{postId} | Получить конкретный пост пользователя |
| POST | /api/users/{userId}/posts | Создать пост для пользователя |
| DELETE | /api/users/{userId}/posts/{postId} | Удалить пост пользователя |

## Структуры данных

```go
// User представляет пользователя
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// Post представляет пост пользователя
type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

// CreatePostRequest — запрос на создание поста
type CreatePostRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
    Error string `json:"error"`
}
```

## Требования к endpoints

### GET /api/users/{userId}/posts

- Возвращает JSON массив постов указанного пользователя
- Если пользователь не найден — 404 Not Found с `{"error": "user not found"}`
- Если у пользователя нет постов — возвращает пустой массив `[]`

Пример запроса:
```
GET /api/users/1/posts
```

Пример ответа (200 OK):
```json
[
    {"id": 1, "user_id": 1, "title": "First Post", "content": "Hello world!", "created_at": "2024-01-15T10:00:00Z"},
    {"id": 2, "user_id": 1, "title": "Second Post", "content": "Another post", "created_at": "2024-01-15T11:00:00Z"}
]
```

### GET /api/users/{userId}/posts/{postId}

- Возвращает конкретный пост пользователя
- Проверяет, что пользователь существует — иначе 404 `{"error": "user not found"}`
- Проверяет, что пост существует — иначе 404 `{"error": "post not found"}`
- Проверяет, что пост принадлежит пользователю — иначе 404 `{"error": "post not found"}`

Пример запроса:
```
GET /api/users/1/posts/2
```

Пример ответа (200 OK):
```json
{"id": 2, "user_id": 1, "title": "Second Post", "content": "Another post", "created_at": "2024-01-15T11:00:00Z"}
```

### POST /api/users/{userId}/posts

- Создаёт новый пост для указанного пользователя
- Если пользователь не найден — 404 Not Found
- Если title пустой — 400 Bad Request `{"error": "title is required"}`
- ID поста генерируется автоматически
- user_id берётся из URL
- created_at устанавливается автоматически
- Возвращает созданный пост со статусом 201 Created

Пример запроса:
```
POST /api/users/1/posts
Content-Type: application/json

{"title": "New Post", "content": "Post content here"}
```

Пример ответа (201 Created):
```json
{"id": 3, "user_id": 1, "title": "New Post", "content": "Post content here", "created_at": "2024-01-15T12:00:00Z"}
```

### DELETE /api/users/{userId}/posts/{postId}

- Удаляет пост пользователя
- Проверяет существование пользователя — иначе 404
- Проверяет существование поста — иначе 404
- Проверяет принадлежность поста пользователю — иначе 404
- Возвращает 204 No Content при успехе

## Примеры использования (curl)

```bash
# Список пользователей
$ curl http://localhost:8080/api/users
[{"id":1,"username":"alice","email":"alice@example.com"},{"id":2,"username":"bob","email":"bob@example.com"}]

# Получить пользователя
$ curl http://localhost:8080/api/users/1
{"id":1,"username":"alice","email":"alice@example.com"}

# Посты пользователя Alice
$ curl http://localhost:8080/api/users/1/posts
[{"id":1,"user_id":1,"title":"Alice First Post","content":"Hello from Alice!","created_at":"..."}]

# Посты несуществующего пользователя
$ curl -i http://localhost:8080/api/users/999/posts
HTTP/1.1 404 Not Found
{"error":"user not found"}

# Конкретный пост пользователя
$ curl http://localhost:8080/api/users/1/posts/1
{"id":1,"user_id":1,"title":"Alice First Post","content":"Hello from Alice!","created_at":"..."}

# Пост другого пользователя (пост 2 принадлежит Bob)
$ curl -i http://localhost:8080/api/users/1/posts/2
HTTP/1.1 404 Not Found
{"error":"post not found"}

# Создать пост
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"title":"New Post","content":"Content here"}' \
    http://localhost:8080/api/users/1/posts
HTTP/1.1 201 Created
{"id":3,"user_id":1,"title":"New Post","content":"Content here","created_at":"..."}

# Создать пост для несуществующего пользователя
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"title":"Post","content":"Content"}' \
    http://localhost:8080/api/users/999/posts
HTTP/1.1 404 Not Found
{"error":"user not found"}

# Удалить пост
$ curl -i -X DELETE http://localhost:8080/api/users/1/posts/1
HTTP/1.1 204 No Content

# Удалить чужой пост
$ curl -i -X DELETE http://localhost:8080/api/users/1/posts/2
HTTP/1.1 404 Not Found
{"error":"post not found"}
```

## Важные моменты

### Множественные параметры пути

В Go 1.22+ можно использовать несколько параметров:

```go
mux.HandleFunc("GET /api/users/{userId}/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
    userId := r.PathValue("userId")
    postId := r.PathValue("postId")
    // ...
})
```

### Проверка принадлежности ресурса

При работе с вложенными ресурсами важно проверять, что ресурс действительно принадлежит родительскому:

```go
// Пост должен принадлежать указанному пользователю
post, exists := postStorage.Get(postId)
if !exists || post.UserID != userId {
    writeError(w, http.StatusNotFound, "post not found")
    return
}
```

### Порядок проверок

1. Проверить существование родительского ресурса (user)
2. Проверить существование дочернего ресурса (post)
3. Проверить принадлежность дочернего родительскому

## Подсказки

- Для получения параметров: `r.PathValue("userId")`, `r.PathValue("postId")`
- Конвертация: `strconv.Atoi(idStr)`
- При проверках валидности — отдавай понятные ошибки
- Для фильтрации постов по user_id используй цикл или функцию filter

## Критерии выполнения

- [ ] GET /api/users возвращает список пользователей
- [ ] GET /api/users/{userId} возвращает пользователя или 404
- [ ] GET /api/users/{userId}/posts возвращает посты пользователя
- [ ] GET /api/users/{userId}/posts для несуществующего user возвращает 404
- [ ] GET /api/users/{userId}/posts/{postId} возвращает пост
- [ ] GET /api/users/{userId}/posts/{postId} проверяет принадлежность поста
- [ ] POST /api/users/{userId}/posts создаёт пост со статусом 201
- [ ] POST для несуществующего user возвращает 404
- [ ] DELETE /api/users/{userId}/posts/{postId} удаляет пост, возвращает 204
- [ ] DELETE проверяет принадлежность поста пользователю
- [ ] Все JSON ответы имеют Content-Type: application/json
- [ ] Код компилируется без ошибок
