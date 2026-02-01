package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

// Storage — потокобезопасное хранилище пользователей
type Storage struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewStorage создаёт новое хранилище с тестовыми данными
func NewStorage() *Storage {
	s := &Storage{
		users:  make(map[int]User),
		nextID: 1,
	}

	// Добавляем тестовых пользователей
	s.Create(CreateUserRequest{Name: "Alice", Email: "alice@example.com", Age: 25})
	s.Create(CreateUserRequest{Name: "Bob", Email: "bob@example.com", Age: 30})

	return s
}

// Create добавляет нового пользователя
func (s *Storage) Create(req CreateUserRequest) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:        s.nextID,
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: time.Now(),
	}

	s.users[user.ID] = user
	s.nextID++

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
func (s *Storage) Update(id int, req UpdateUserRequest) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return User{}, false
	}

	// Обновляем поля, сохраняя ID и CreatedAt
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

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

// Глобальное хранилище (для простоты примера)
var storage = NewStorage()

// writeJSON отправляет JSON ответ с указанным статус-кодом
func writeJSON(w http.ResponseWriter, status int, data any) {
	// TODO: реализуй функцию
	// 1. Установи заголовок Content-Type: application/json
	// 2. Установи статус-код: w.WriteHeader(status)
	// 3. Закодируй data в JSON: json.NewEncoder(w).Encode(data)
	_ = w
	_ = status
	_ = data
}

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, status int, message string) {
	// TODO: реализуй функцию
	// Используй writeJSON с ErrorResponse{Error: message}
	_ = w
	_ = status
	_ = message
}

// parseJSON читает JSON из тела запроса в структуру v
func parseJSON(r *http.Request, v any) error {
	// TODO: реализуй функцию
	// 1. Используй json.NewDecoder(r.Body).Decode(v)
	// 2. Верни ошибку, если декодирование не удалось
	_ = r
	_ = v
	return nil
}

// listUsersHandler обрабатывает GET /api/users
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи всех пользователей: users := storage.GetAll()
	// 2. Отправь JSON ответ: writeJSON(w, http.StatusOK, users)
	_ = r
	_ = w
}

// getUserHandler обрабатывает GET /api/users/{id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи id из пути: idStr := r.PathValue("id")
	// 2. Конвертируй в int: id, err := strconv.Atoi(idStr)
	// 3. При ошибке конвертации: writeError(w, http.StatusBadRequest, "invalid user ID")
	// 4. Получи пользователя: user, exists := storage.Get(id)
	// 5. Если не существует: writeError(w, http.StatusNotFound, "user not found")
	// 6. Иначе: writeJSON(w, http.StatusOK, user)
	_ = r
	_ = w
}

// createUserHandler обрабатывает POST /api/users
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Создай переменную var req CreateUserRequest
	// 2. Распарси JSON: if err := parseJSON(r, &req); err != nil { ... }
	// 3. При ошибке парсинга: writeError(w, http.StatusBadRequest, "invalid request body")
	// 4. Проверь, что Name не пустой (strings.TrimSpace)
	// 5. Если пустой: writeError(w, http.StatusBadRequest, "name is required")
	// 6. Создай пользователя: user := storage.Create(req)
	// 7. Верни пользователя: writeJSON(w, http.StatusCreated, user)
	_ = r
	_ = w
}

// updateUserHandler обрабатывает PUT /api/users/{id}
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи и проверь ID из пути
	// 2. Распарси JSON в UpdateUserRequest
	// 3. Проверь, что Name не пустой
	// 4. Обнови пользователя: user, exists := storage.Update(id, req)
	// 5. Если не существует: writeError(w, http.StatusNotFound, "user not found")
	// 6. Иначе: writeJSON(w, http.StatusOK, user)
	_ = r
	_ = w
}

// deleteUserHandler обрабатывает DELETE /api/users/{id}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи и проверь ID из пути
	// 2. Удали пользователя: deleted := storage.Delete(id)
	// 3. Если не существовал: writeError(w, http.StatusNotFound, "user not found")
	// 4. Иначе: w.WriteHeader(http.StatusNoContent) — БЕЗ тела ответа
	_ = r
	_ = w
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики
	// mux.HandleFunc("GET /api/users", listUsersHandler)
	// mux.HandleFunc("GET /api/users/{id}", getUserHandler)
	// mux.HandleFunc("POST /api/users", createUserHandler)
	// mux.HandleFunc("PUT /api/users/{id}", updateUserHandler)
	// mux.HandleFunc("DELETE /api/users/{id}", deleteUserHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /api/users      - List all users")
	fmt.Println("  GET    /api/users/{id} - Get user by ID")
	fmt.Println("  POST   /api/users      - Create new user")
	fmt.Println("  PUT    /api/users/{id} - Update user")
	fmt.Println("  DELETE /api/users/{id} - Delete user")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	// Используем импорты, чтобы код компилировался
	_ = json.NewEncoder
	_ = strconv.Atoi
	_ = strings.TrimSpace
}
