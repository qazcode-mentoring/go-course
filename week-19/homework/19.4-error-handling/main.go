package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ==================== Система обработки ошибок ====================

// AppError — типизированная ошибка приложения
type AppError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Status  int               `json:"-"` // не сериализуется в JSON
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"` // оригинальная ошибка для логов
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NotFound создаёт ошибку "ресурс не найден"
func NotFound(resource string) *AppError {
	// TODO: реализуй функцию
	// return &AppError{
	//     Code:    "NOT_FOUND",
	//     Message: fmt.Sprintf("%s not found", resource),
	//     Status:  http.StatusNotFound,
	// }
	_ = resource
	return nil
}

// BadRequest создаёт ошибку "неверный запрос"
func BadRequest(message string) *AppError {
	// TODO: реализуй функцию
	// return &AppError{
	//     Code:    "BAD_REQUEST",
	//     Message: message,
	//     Status:  http.StatusBadRequest,
	// }
	_ = message
	return nil
}

// ValidationError создаёт ошибку валидации с деталями по полям
func ValidationError(details map[string]string) *AppError {
	// TODO: реализуй функцию
	// return &AppError{
	//     Code:    "VALIDATION_ERROR",
	//     Message: "validation failed",
	//     Status:  http.StatusUnprocessableEntity,
	//     Details: details,
	// }
	_ = details
	return nil
}

// InternalError создаёт внутреннюю ошибку сервера
func InternalError(err error) *AppError {
	// TODO: реализуй функцию
	// return &AppError{
	//     Code:    "INTERNAL_ERROR",
	//     Message: "internal server error",
	//     Status:  http.StatusInternalServerError,
	//     Err:     err,
	// }
	_ = err
	return nil
}

// ==================== AppHandler и WrapHandler ====================

// AppHandler — обработчик, возвращающий ошибку
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// WrapHandler оборачивает AppHandler в стандартный http.HandlerFunc
func WrapHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: реализуй функцию
		// 1. Вызови обработчик: err := h(w, r)
		// 2. Если err == nil — ничего не делаем, return
		// 3. Проверь, является ли ошибка AppError: var appErr *AppError
		//    if errors.As(err, &appErr) { ... }
		// 4. Если AppError — логируем оригинальную ошибку (если есть) и отправляем JSON
		// 5. Если неизвестная ошибка — логируем и отправляем INTERNAL_ERROR

		// Временная заглушка — просто вызываем обработчик
		_ = h(w, r)
	}
}

// ==================== Модели данных ====================

// Task представляет задачу
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateTaskRequest — запрос на создание задачи
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// Validate проверяет валидность запроса на создание
func (r CreateTaskRequest) Validate() map[string]string {
	// TODO: реализуй функцию
	// errs := make(map[string]string)
	//
	// title := strings.TrimSpace(r.Title)
	// if title == "" {
	//     errs["title"] = "title is required"
	// } else if len(title) < 3 {
	//     errs["title"] = "title must be at least 3 characters"
	// } else if len(title) > 100 {
	//     errs["title"] = "title must be at most 100 characters"
	// }
	//
	// if r.Priority < 1 || r.Priority > 5 {
	//     errs["priority"] = "priority must be between 1 and 5"
	// }
	//
	// return errs
	return nil
}

// UpdateTaskRequest — запрос на обновление задачи
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
}

// Validate проверяет валидность запроса на обновление
func (r UpdateTaskRequest) Validate() map[string]string {
	// TODO: реализуй функцию
	// errs := make(map[string]string)
	//
	// title := strings.TrimSpace(r.Title)
	// if title == "" {
	//     errs["title"] = "title is required"
	// } else if len(title) < 3 {
	//     errs["title"] = "title must be at least 3 characters"
	// } else if len(title) > 100 {
	//     errs["title"] = "title must be at most 100 characters"
	// }
	//
	// validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
	// if !validStatuses[r.Status] {
	//     errs["status"] = "status must be one of: pending, in_progress, done"
	// }
	//
	// if r.Priority < 1 || r.Priority > 5 {
	//     errs["priority"] = "priority must be between 1 and 5"
	// }
	//
	// return errs
	return nil
}

// ==================== Storage ====================

// TaskStorage — хранилище задач
type TaskStorage struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

// NewTaskStorage создаёт хранилище с тестовыми данными
func NewTaskStorage() *TaskStorage {
	now := time.Now()
	return &TaskStorage{
		tasks: map[int]Task{
			1: {ID: 1, Title: "Learn Go", Description: "Complete Go course", Status: "in_progress", Priority: 5, CreatedAt: now.Add(-48 * time.Hour)},
			2: {ID: 2, Title: "Build REST API", Description: "Create REST API project", Status: "pending", Priority: 4, CreatedAt: now.Add(-24 * time.Hour)},
		},
		nextID: 3,
	}
}

// Get возвращает задачу по ID
func (s *TaskStorage) Get(id int) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	return task, exists
}

// GetAll возвращает все задачи
func (s *TaskStorage) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// Create создаёт новую задачу
func (s *TaskStorage) Create(req CreateTaskRequest) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:          s.nextID,
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
		Status:      "pending",
		Priority:    req.Priority,
		CreatedAt:   time.Now(),
	}

	s.tasks[task.ID] = task
	s.nextID++

	return task
}

// Update обновляет задачу
func (s *TaskStorage) Update(id int, req UpdateTaskRequest) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return Task{}, false
	}

	task.Title = strings.TrimSpace(req.Title)
	task.Description = req.Description
	task.Status = req.Status
	task.Priority = req.Priority

	s.tasks[id] = task
	return task, true
}

// Delete удаляет задачу
func (s *TaskStorage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}

	delete(s.tasks, id)
	return true
}

// Глобальное хранилище
var storage = NewTaskStorage()

// ==================== Вспомогательные функции ====================

// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, status int, data any) {
	// TODO: реализуй функцию
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(status)
	// json.NewEncoder(w).Encode(data)
	_ = w
	_ = status
	_ = data
}

// ==================== Обработчики ====================

// listTasks обрабатывает GET /api/tasks
func listTasks(w http.ResponseWriter, r *http.Request) error {
	// TODO: реализуй обработчик
	// tasks := storage.GetAll()
	// writeJSON(w, http.StatusOK, tasks)
	// return nil
	_ = r
	_ = w
	return nil
}

// getTask обрабатывает GET /api/tasks/{id}
func getTask(w http.ResponseWriter, r *http.Request) error {
	// TODO: реализуй обработчик
	// 1. Получи id из пути: idStr := r.PathValue("id")
	// 2. Конвертируй в int: id, err := strconv.Atoi(idStr)
	// 3. При ошибке конвертации: return BadRequest("invalid task ID")
	// 4. Получи задачу: task, exists := storage.Get(id)
	// 5. Если не существует: return NotFound("task")
	// 6. Отправь ответ: writeJSON(w, http.StatusOK, task)
	// 7. return nil
	_ = r
	_ = w
	return nil
}

// createTask обрабатывает POST /api/tasks
func createTask(w http.ResponseWriter, r *http.Request) error {
	// TODO: реализуй обработчик
	// 1. Распарси JSON: var req CreateTaskRequest
	//    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//        return BadRequest("invalid request body")
	//    }
	// 2. Валидируй: if errs := req.Validate(); len(errs) > 0 {
	//        return ValidationError(errs)
	//    }
	// 3. Создай задачу: task := storage.Create(req)
	// 4. Отправь ответ: writeJSON(w, http.StatusCreated, task)
	// 5. return nil
	_ = r
	_ = w
	return nil
}

// updateTask обрабатывает PUT /api/tasks/{id}
func updateTask(w http.ResponseWriter, r *http.Request) error {
	// TODO: реализуй обработчик
	// 1. Получи и проверь ID из пути
	// 2. Распарси JSON в UpdateTaskRequest
	// 3. Валидируй запрос
	// 4. Обнови задачу: task, exists := storage.Update(id, req)
	// 5. Если не существует: return NotFound("task")
	// 6. Отправь ответ: writeJSON(w, http.StatusOK, task)
	// 7. return nil
	_ = r
	_ = w
	return nil
}

// deleteTask обрабатывает DELETE /api/tasks/{id}
func deleteTask(w http.ResponseWriter, r *http.Request) error {
	// TODO: реализуй обработчик
	// 1. Получи и проверь ID из пути
	// 2. Удали задачу: deleted := storage.Delete(id)
	// 3. Если не существовала: return NotFound("task")
	// 4. Отправь 204: w.WriteHeader(http.StatusNoContent)
	// 5. return nil
	_ = r
	_ = w
	return nil
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики с WrapHandler
	// mux.HandleFunc("GET /api/tasks", WrapHandler(listTasks))
	// mux.HandleFunc("GET /api/tasks/{id}", WrapHandler(getTask))
	// mux.HandleFunc("POST /api/tasks", WrapHandler(createTask))
	// mux.HandleFunc("PUT /api/tasks/{id}", WrapHandler(updateTask))
	// mux.HandleFunc("DELETE /api/tasks/{id}", WrapHandler(deleteTask))

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /api/tasks      - List all tasks")
	fmt.Println("  GET    /api/tasks/{id} - Get task by ID")
	fmt.Println("  POST   /api/tasks      - Create new task")
	fmt.Println("  PUT    /api/tasks/{id} - Update task")
	fmt.Println("  DELETE /api/tasks/{id} - Delete task")
	fmt.Println()
	fmt.Println("Error codes:")
	fmt.Println("  NOT_FOUND        - Resource not found (404)")
	fmt.Println("  BAD_REQUEST      - Invalid request (400)")
	fmt.Println("  VALIDATION_ERROR - Validation failed (422)")
	fmt.Println("  INTERNAL_ERROR   - Internal server error (500)")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	// Используем импорты
	_ = json.NewEncoder
	_ = strconv.Atoi
	_ = strings.TrimSpace
	_ = errors.As
}
