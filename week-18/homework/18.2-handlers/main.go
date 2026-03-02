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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, user := range users {
		_, err := fmt.Fprintf(w, "%d: %s (%d лет)\n", user.ID, user.Name, user.Age)
		if err != nil {
			log.Printf("failed to write: %v", err)
			return
		}
	}
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
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Failed to convert str to int", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.ID == id {
			_, err = fmt.Fprintf(w, "%d: %s (%d лет)", user.ID, user.Name, user.Age)
			if err != nil {
				http.Error(w, "Failed to display", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name, ageStr := r.FormValue("name"), r.FormValue("age")
	if name == "" || ageStr == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Failed to convert str to int", http.StatusBadRequest)
		return
	}

	user := User{
		ID:   nextID,
		Name: name,
		Age:  age,
	}
	users = append(users, user)
	nextID++

	w.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprintf(w, "Created: %d: %s (%d лет)", user.ID, user.Name, user.Age)
	if err != nil {
		http.Error(w, "Failed to display", http.StatusInternalServerError)
		return
	}

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
	query := r.URL.Query().Get("name")
	if query == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	count := 0
	var foundUsers []User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
			count++
			foundUsers = append(foundUsers, user)
		}
	}

	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err := fmt.Fprintf(w, "Found %d user(s)\n", count)
	if err != nil {
		http.Error(w, "Failed to display", http.StatusInternalServerError)
		return
	}

	for _, user := range foundUsers {
		_, err = fmt.Fprintf(w, "%d: %s (%d лет)\n", user.ID, user.Name, user.Age)
		if err != nil {
			http.Error(w, "Failed to display", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Создаём новый ServeMux (Go 1.22+ с поддержкой методов и параметров)
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики с использованием паттернов Go 1.22+
	// mux.HandleFunc("GET /users", listUsersHandler)
	// mux.HandleFunc("GET /users/{id}", getUserHandler)
	// mux.HandleFunc("POST /users", createUserHandler)
	// mux.HandleFunc("GET /search", searchHandler)

	mux.HandleFunc("GET /users", listUsersHandler)
	mux.HandleFunc("GET /users/{id}", getUserHandler)
	mux.HandleFunc("POST /users", createUserHandler)
	mux.HandleFunc("GET /search", searchHandler)

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
}
