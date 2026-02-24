package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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

// Хранилище книг (in-memory)
var books = []Book{
	{ID: 1, Title: "1984", Author: "George Orwell", Year: 1949},
	{ID: 2, Title: "Brave New World", Author: "Aldous Huxley", Year: 1932},
}
var nextBookID = 3

// writeJSON отправляет JSON ответ с указанным статус-кодом
func writeJSON(w http.ResponseWriter, status int, data any) {
	// TODO: реализуй функцию
	// 1. Установи заголовок Content-Type: application/json
	// 2. Установи статус-код: w.WriteHeader(status)
	// 3. Закодируй data в JSON: json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode to json", http.StatusBadRequest)
		return
	}
}

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, status int, message string) {
	// TODO: реализуй функцию
	// Используй writeJSON с ErrorResponse{Error: message}
	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, status, ErrorResponse{Error: message})
}

// parseJSON читает JSON из тела запроса в структуру v
func parseJSON(r *http.Request, v any) error {
	// TODO: реализуй функцию
	// 1. Используй json.NewDecoder(r.Body).Decode(v)
	// 2. Верни ошибку, если декодирование не удалось
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}

// listBooksHandler обрабатывает GET /api/books
// Возвращает список всех книг в формате JSON
func listBooksHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// Используй writeJSON для отправки списка books

	w.Header().Set("Content-Type", "application/json")

	writeJSON(w, http.StatusOK, books)

}

// getBookHandler обрабатывает GET /api/books/{id}
// Возвращает книгу по ID или 404, если не найдена
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи id из пути: r.PathValue("id")
	// 2. Конвертируй в int: strconv.Atoi(idStr)
	// 3. Найди книгу в slice books
	// 4. Если не найдена — writeError(w, http.StatusNotFound, "book not found")
	// 5. Если найдена — writeJSON(w, http.StatusOK, book)

	idStr := r.PathValue("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Id is empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to convert str to int")
		return
	}

	for _, book := range books {
		if book.ID == id {
			writeJSON(w, http.StatusOK, book)
			return
		}
	}

	writeError(w, http.StatusNotFound, "book not found")
}

// createBookHandler обрабатывает POST /api/books
// Создаёт новую книгу из JSON тела запроса
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Создай переменную var book Book
	// 2. Распарси JSON: if err := parseJSON(r, &book); err != nil
	// 3. При ошибке парсинга — writeError(w, http.StatusBadRequest, "invalid JSON")
	// 4. Проверь, что Title не пустой
	// 5. Присвой book.ID = nextBookID, увеличь nextBookID
	// 6. Добавь книгу в slice: books = append(books, book)
	// 7. Верни книгу: writeJSON(w, http.StatusCreated, book)
	var book Book
	if err := parseJSON(r, &book); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if book.Title == "" {
		writeError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	book.ID = nextBookID
	nextBookID++

	books = append(books, book)
	writeJSON(w, http.StatusCreated, book)
}

// deleteBookHandler обрабатывает DELETE /api/books/{id}
// Удаляет книгу по ID
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи id из пути
	// 2. Найди индекс книги в slice
	// 3. Если не найдена — writeError(w, http.StatusNotFound, "book not found")
	// 4. Удали из slice: books = append(books[:index], books[index+1:]...)
	// 5. Верни успех: writeJSON(w, http.StatusOK, SuccessResponse{Message: "book deleted successfully"})

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to convert str to int")
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			writeJSON(w, http.StatusOK, SuccessResponse{Message: "book deleted successfully"})
			return
		}
	}

	writeError(w, http.StatusNotFound, "book not found")
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики
	// mux.HandleFunc("GET /api/books", listBooksHandler)
	// mux.HandleFunc("GET /api/books/{id}", getBookHandler)
	// mux.HandleFunc("POST /api/books", createBookHandler)
	// mux.HandleFunc("DELETE /api/books/{id}", deleteBookHandler)

	mux.HandleFunc("GET /api/books", listBooksHandler)
	mux.HandleFunc("GET /api/books/{id}", getBookHandler)
	mux.HandleFunc("POST /api/books", createBookHandler)
	mux.HandleFunc("DELETE /api/books/{id}", deleteBookHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /api/books      - List all books")
	fmt.Println("  GET    /api/books/{id} - Get book by ID")
	fmt.Println("  POST   /api/books      - Create new book")
	fmt.Println("  DELETE /api/books/{id} - Delete book")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
