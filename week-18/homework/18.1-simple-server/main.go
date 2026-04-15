package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler обрабатывает запросы к корневому пути "/"
// Должен возвращать "Hello, World!"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Установи Content-Type: text/plain
	// 2. Верни текст "Hello, World!"
	// Подсказка: используй fmt.Fprintf(w, "Hello, World!")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, world!")
}

// healthHandler обрабатывает запросы к "/health"
// Используется для проверки работоспособности сервера (health check)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Установи Content-Type: text/plain
	// 2. Верни текст "OK"
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "OK")
}

// infoHandler обрабатывает запросы к "/info"
// Возвращает информацию о текущем запросе
func infoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Установи Content-Type: text/plain
	// 2. Выведи информацию о запросе:
	//    - r.Method — HTTP метод
	//    - r.URL.Path — путь запроса
	//    - r.Header.Get("User-Agent") — User-Agent клиента
	//    - r.RemoteAddr — адрес клиента
	//
	// Формат вывода:
	// Method: GET
	// Path: /info
	// User-Agent: curl/8.1.2
	// Remote Address: 127.0.0.1:52341
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w,
		"Method: %s\n"+
			"Path: %s\nUser-Agent: %s\n"+
			"Remote Address: %s\n",
		r.Method, r.URL.Path, r.Header.Get("User-Agent"), r.RemoteAddr)
}

func main() {
	// TODO: зарегистрируй обработчики с помощью http.HandleFunc
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/info", infoHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET /       - Hello World")
	fmt.Println("  GET /health - Health check")
	fmt.Println("  GET /info   - Request info")

	// TODO: запусти сервер с помощью http.ListenAndServe
	log.Fatal(http.ListenAndServe(addr, nil))
	log.Println("Server implementation not complete")
}
