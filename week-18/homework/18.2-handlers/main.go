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
		fmt.Fprintf(w, "%d: %s (%d лет)\n", user.ID, user.Name, user.Age)
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
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	for _, user := range users {
		if user.ID == id {
			fmt.Fprintf(w,
				"ID: %d\n"+
					"Name: %s\n"+
					"Age: %d\n",
				user.ID, user.Name, user.Age,
			)
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
		http.Error(w, "Form parsing error", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	ageStr := r.FormValue("age")

	if name == "" || ageStr == "" {
		http.Error(w, "Name and age are required", http.StatusBadRequest)
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Age must be a number", http.StatusBadRequest)
		return
	}

	user := User{
		ID:   nextID,
		Name: name,
		Age:  age,
	}
	users = append(users, user)

	nextID++

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created:\n ID: %d\n Name:%s\n Age:%d\n", user.ID, user.Name, user.Age)
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

	var foundUsers []User

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
			foundUsers = append(foundUsers, user)
		}
	}

	if len(foundUsers) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Found %d user(s):\n", len(foundUsers))

	for _, user := range foundUsers {
		fmt.Fprintf(w, "%d: %s (%d лет)\n", user.ID, user.Name, user.Age)
	}
}

func main() {
	// Создаём новый ServeMux (Go 1.22+ с поддержкой методов и параметров)
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики с использованием паттернов Go 1.22+
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
