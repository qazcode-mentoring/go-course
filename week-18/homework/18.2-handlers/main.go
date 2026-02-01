package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// User представляет пользователя системы
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

// listUsersHandler обрабатывает GET /users
// Возвращает список всех пользователей
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Установи Content-Type: text/plain; charset=utf-8
	// 2. Пройди по списку users и выведи каждого в формате:
	//    "ID: Name (Age лет)"
	// Пример: "1: Alice (25 лет)"
	_ = r // убрать после реализации
	_ = w // убрать после реализации
}

// getUserHandler обрабатывает GET /users/{id}
// Возвращает пользователя по ID
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи id из пути: r.PathValue("id")
	// 2. Конвертируй в int: strconv.Atoi(idStr)
	// 3. Найди пользователя в slice users
	// 4. Если не найден — http.Error(w, "User not found", http.StatusNotFound)
	// 5. Если найден — выведи информацию:
	//    ID: ...
	//    Name: ...
	//    Age: ...
	_ = r // убрать после реализации
	_ = w // убрать после реализации
}

// createUserHandler обрабатывает POST /users
// Создаёт нового пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Распарси форму: r.ParseForm()
	// 2. Получи данные: r.FormValue("name"), r.FormValue("age")
	// 3. Проверь, что name и age не пустые, иначе — 400 Bad Request
	// 4. Конвертируй age в int
	// 5. Создай нового пользователя с nextID, добавь в users
	// 6. Увеличь nextID
	// 7. Верни статус 201 Created: w.WriteHeader(http.StatusCreated)
	// 8. Выведи информацию о созданном пользователе
	_ = r // убрать после реализации
	_ = w // убрать после реализации
}

// searchHandler обрабатывает GET /search?name=...
// Ищет пользователей по имени (case-insensitive)
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи query параметр: r.URL.Query().Get("name")
	// 2. Если параметр пустой — 400 Bad Request с сообщением "name parameter is required"
	// 3. Найди всех пользователей, чьё имя содержит подстроку (case-insensitive)
	//    Подсказка: strings.Contains(strings.ToLower(user.Name), strings.ToLower(query))
	// 4. Выведи результаты:
	//    "Found N user(s):"
	//    "ID: Name (Age лет)"
	_ = r // убрать после реализации
	_ = w // убрать после реализации
}

func main() {
	// Создаём новый ServeMux (Go 1.22+ с поддержкой методов и параметров)
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики с использованием паттернов Go 1.22+
	// mux.HandleFunc("GET /users", listUsersHandler)
	// mux.HandleFunc("GET /users/{id}", getUserHandler)
	// mux.HandleFunc("POST /users", createUserHandler)
	// mux.HandleFunc("GET /search", searchHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /users      - List all users")
	fmt.Println("  GET  /users/{id} - Get user by ID")
	fmt.Println("  POST /users      - Create new user")
	fmt.Println("  GET  /search     - Search users by name")

	// Запускаем сервер с нашим mux
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	// Убираем предупреждения о неиспользуемых импортах
	_ = strconv.Atoi
	_ = strings.Contains
}
