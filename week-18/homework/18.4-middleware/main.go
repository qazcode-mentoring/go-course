package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware — тип функции middleware
type Middleware func(http.Handler) http.Handler

// statusRecorder — обёртка для ResponseWriter, сохраняющая статус-код
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader перехватывает статус-код перед записью
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует информацию о каждом запросе
// Формат: "METHOD /path STATUS TIME"
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: реализуй middleware
		// 1. Запомни время начала: start := time.Now()
		// 2. Создай statusRecorder: rec := &statusRecorder{ResponseWriter: w, status: 200}
		// 3. Вызови следующий handler: next.ServeHTTP(rec, r)
		// 4. Залогируй результат: log.Printf("%s %s %d %v", r.Method, r.URL.Path, rec.status, time.Since(start))

		// Временная заглушка — просто передаём управление
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware перехватывает паники и возвращает 500
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: реализуй middleware
		// 1. Добавь defer с recover():
		//    defer func() {
		//        if err := recover(); err != nil {
		//            log.Printf("Panic recovered: %v", err)
		//            w.Header().Set("Content-Type", "application/json")
		//            w.WriteHeader(http.StatusInternalServerError)
		//            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		//        }
		//    }()
		// 2. Вызови следующий handler: next.ServeHTTP(w, r)

		// Временная заглушка — просто передаём управление
		next.ServeHTTP(w, r)
	})
}

// Chain объединяет middleware в цепочку
// Chain(handler, A, B, C) эквивалентно A(B(C(handler)))
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// TODO: реализуй функцию
	// Применяй middleware в обратном порядке:
	// for i := len(middlewares) - 1; i >= 0; i-- {
	//     handler = middlewares[i](handler)
	// }
	// return handler

	// Временная заглушка
	return handler
}

// --- Тестовые обработчики ---

// homeHandler — простой обработчик для корневого пути
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, World!")
	_ = r // используется в сигнатуре
}

// dataHandler — возвращает JSON данные
func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := map[string]any{
		"message":   "Hello from API",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = r // используется в сигнатуре
}

// panicHandler — вызывает панику для тестирования recovery
func panicHandler(w http.ResponseWriter, r *http.Request) {
	_ = w // не используется — паника раньше
	_ = r // используется в сигнатуре
	panic("test panic")
}

// slowHandler — медленный обработчик для тестирования логирования времени
func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Done!")
	_ = r // используется в сигнатуре
}

// notFoundHandler — возвращает 404 для тестирования логирования статусов
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
	_ = r // используется в сигнатуре
}

func main() {
	mux := http.NewServeMux()

	// Регистрация обработчиков
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /api/data", dataHandler)
	mux.HandleFunc("GET /panic", panicHandler)
	mux.HandleFunc("GET /slow", slowHandler)
	mux.HandleFunc("GET /notfound", notFoundHandler)

	// TODO: оберни mux в цепочку middleware
	// handler := Chain(mux, LoggingMiddleware, RecoveryMiddleware)
	handler := mux // временная заглушка

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET /          - Home page")
	fmt.Println("  GET /api/data  - JSON API endpoint")
	fmt.Println("  GET /panic     - Triggers panic (test recovery)")
	fmt.Println("  GET /slow      - Slow endpoint (100ms)")
	fmt.Println("  GET /notfound  - Returns 404")
	fmt.Println()
	fmt.Println("Middleware chain: Logging -> Recovery -> Handler")

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}
