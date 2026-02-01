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

// UserStorage — хранилище пользователей
type UserStorage struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// PostStorage — хранилище постов
type PostStorage struct {
	mu     sync.RWMutex
	posts  map[int]Post
	nextID int
}

// NewUserStorage создаёт хранилище с тестовыми пользователями
func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: map[int]User{
			1: {ID: 1, Username: "alice", Email: "alice@example.com"},
			2: {ID: 2, Username: "bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

// NewPostStorage создаёт хранилище с тестовыми постами
func NewPostStorage() *PostStorage {
	now := time.Now()
	return &PostStorage{
		posts: map[int]Post{
			1: {ID: 1, UserID: 1, Title: "Alice First Post", Content: "Hello from Alice!", CreatedAt: now},
			2: {ID: 2, UserID: 2, Title: "Bob First Post", Content: "Hello from Bob!", CreatedAt: now},
		},
		nextID: 3,
	}
}

// Get возвращает пользователя по ID
func (s *UserStorage) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[id]
	return user, exists
}

// GetAll возвращает всех пользователей
func (s *UserStorage) GetAll() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

// Exists проверяет существование пользователя
func (s *UserStorage) Exists(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.users[id]
	return exists
}

// Get возвращает пост по ID
func (s *PostStorage) Get(id int) (Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, exists := s.posts[id]
	return post, exists
}

// GetByUserID возвращает все посты пользователя
func (s *PostStorage) GetByUserID(userID int) []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	posts := make([]Post, 0)
	for _, p := range s.posts {
		if p.UserID == userID {
			posts = append(posts, p)
		}
	}
	return posts
}

// Create создаёт новый пост
func (s *PostStorage) Create(userID int, req CreatePostRequest) Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := Post{
		ID:        s.nextID,
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	s.posts[post.ID] = post
	s.nextID++

	return post
}

// Delete удаляет пост
func (s *PostStorage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; !exists {
		return false
	}

	delete(s.posts, id)
	return true
}

// Глобальные хранилища
var userStorage = NewUserStorage()
var postStorage = NewPostStorage()

// writeJSON отправляет JSON ответ
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

// parseJSON читает JSON из тела запроса
func parseJSON(r *http.Request, v any) error {
	// TODO: реализуй функцию
	_ = r
	_ = v
	return nil
}

// parseIntParam парсит параметр пути как int
func parseIntParam(r *http.Request, name string) (int, error) {
	// TODO: реализуй функцию
	// 1. Получи значение: r.PathValue(name)
	// 2. Конвертируй: strconv.Atoi(...)
	// 3. Верни результат или ошибку
	_ = r
	_ = name
	return 0, nil
}

// --- Обработчики пользователей ---

// listUsersHandler обрабатывает GET /api/users
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// Верни список всех пользователей
	_ = r
	_ = w
}

// getUserHandler обрабатывает GET /api/users/{userId}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи userId из пути
	// 2. Найди пользователя или верни 404
	_ = r
	_ = w
}

// --- Обработчики постов ---

// listUserPostsHandler обрабатывает GET /api/users/{userId}/posts
func listUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи userId из пути: r.PathValue("userId")
	// 2. Конвертируй в int
	// 3. Проверь, что пользователь существует: userStorage.Exists(userID)
	// 4. Если не существует: writeError(w, http.StatusNotFound, "user not found")
	// 5. Получи посты: posts := postStorage.GetByUserID(userID)
	// 6. Верни посты: writeJSON(w, http.StatusOK, posts)
	_ = r
	_ = w
}

// getUserPostHandler обрабатывает GET /api/users/{userId}/posts/{postId}
func getUserPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи userId и postId из пути
	// 2. Проверь, что пользователь существует — иначе 404 "user not found"
	// 3. Получи пост: post, exists := postStorage.Get(postID)
	// 4. Проверь, что пост существует — иначе 404 "post not found"
	// 5. Проверь, что post.UserID == userID — иначе 404 "post not found"
	// 6. Верни пост
	_ = r
	_ = w
}

// createUserPostHandler обрабатывает POST /api/users/{userId}/posts
func createUserPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи userId из пути
	// 2. Проверь, что пользователь существует — иначе 404
	// 3. Распарси JSON в CreatePostRequest
	// 4. Проверь, что title не пустой — иначе 400 "title is required"
	// 5. Создай пост: post := postStorage.Create(userID, req)
	// 6. Верни пост: writeJSON(w, http.StatusCreated, post)
	_ = r
	_ = w
}

// deleteUserPostHandler обрабатывает DELETE /api/users/{userId}/posts/{postId}
func deleteUserPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи userId и postId из пути
	// 2. Проверь, что пользователь существует — иначе 404
	// 3. Получи пост и проверь его существование — иначе 404
	// 4. Проверь, что post.UserID == userID — иначе 404
	// 5. Удали пост: postStorage.Delete(postID)
	// 6. Верни 204: w.WriteHeader(http.StatusNoContent)
	_ = r
	_ = w
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики
	// Пользователи
	// mux.HandleFunc("GET /api/users", listUsersHandler)
	// mux.HandleFunc("GET /api/users/{userId}", getUserHandler)

	// Посты пользователей (вложенные ресурсы)
	// mux.HandleFunc("GET /api/users/{userId}/posts", listUserPostsHandler)
	// mux.HandleFunc("GET /api/users/{userId}/posts/{postId}", getUserPostHandler)
	// mux.HandleFunc("POST /api/users/{userId}/posts", createUserPostHandler)
	// mux.HandleFunc("DELETE /api/users/{userId}/posts/{postId}", deleteUserPostHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /api/users                         - List all users")
	fmt.Println("  GET    /api/users/{userId}                - Get user by ID")
	fmt.Println("  GET    /api/users/{userId}/posts          - List user's posts")
	fmt.Println("  GET    /api/users/{userId}/posts/{postId} - Get specific post")
	fmt.Println("  POST   /api/users/{userId}/posts          - Create post for user")
	fmt.Println("  DELETE /api/users/{userId}/posts/{postId} - Delete user's post")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	// Используем импорты
	_ = json.NewEncoder
	_ = strconv.Atoi
	_ = strings.TrimSpace
}
