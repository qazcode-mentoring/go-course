package main

import (
	"errors"
	"fmt"
)

// Ошибки
var (
	ErrUserNotFound = errors.New("user not found")
	ErrSendFailed   = errors.New("failed to send message")
)

// User представляет пользователя системы
type User struct {
	ID    int
	Name  string
	Email string
}

// UserRepository — интерфейс для работы с пользователями
type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
}

// MessageSender — интерфейс для отправки сообщений
type MessageSender interface {
	Send(to, subject, body string) error
}

// NotificationService — сервис отправки уведомлений
type NotificationService struct {
	users  UserRepository
	sender MessageSender
}

// NewNotificationService создаёт новый сервис уведомлений
func NewNotificationService(users UserRepository, sender MessageSender) *NotificationService {
	return &NotificationService{
		users:  users,
		sender: sender,
	}
}

// NotifyUser отправляет уведомление пользователю по ID
func (s *NotificationService) NotifyUser(userID int, message string) error {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	subject := fmt.Sprintf("Уведомление для %s", user.Name)
	if err := s.sender.Send(user.Email, subject, message); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

// NotifyByEmail отправляет уведомление пользователю по email
func (s *NotificationService) NotifyByEmail(email, message string) error {
	user, err := s.users.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("getting user by email: %w", err)
	}

	subject := fmt.Sprintf("Уведомление для %s", user.Name)
	if err := s.sender.Send(user.Email, subject, message); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

// BroadcastToUsers отправляет сообщение нескольким пользователям
// Возвращает количество успешно отправленных и неудачных отправок
func (s *NotificationService) BroadcastToUsers(userIDs []int, message string) (sent int, failed int) {
	for _, id := range userIDs {
		err := s.NotifyUser(id, message)
		if err != nil {
			failed++
		} else {
			sent++
		}
	}
	return sent, failed
}

// === Реальные реализации (для демонстрации, не для тестов) ===

// InMemoryUserRepository — простая реализация репозитория в памяти
type InMemoryUserRepository struct {
	users map[int]*User
}

// NewInMemoryUserRepository создаёт репозиторий с тестовыми данными
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: map[int]*User{
			1: {ID: 1, Name: "Алиса", Email: "alice@example.com"},
			2: {ID: 2, Name: "Боб", Email: "bob@example.com"},
			3: {ID: 3, Name: "Чарли", Email: "charlie@example.com"},
		},
	}
}

func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// ConsoleMessageSender — отправщик, который печатает в консоль
type ConsoleMessageSender struct{}

func (s *ConsoleMessageSender) Send(to, subject, body string) error {
	fmt.Printf("=== Отправка сообщения ===\n")
	fmt.Printf("Кому: %s\n", to)
	fmt.Printf("Тема: %s\n", subject)
	fmt.Printf("Текст: %s\n", body)
	fmt.Println()
	return nil
}

func main() {
	fmt.Println("=== NotificationService для мокирования ===")

	// Создаём сервис с реальными реализациями
	repo := NewInMemoryUserRepository()
	sender := &ConsoleMessageSender{}
	service := NewNotificationService(repo, sender)

	// Отправляем уведомление пользователю
	fmt.Println("--- NotifyUser ---")
	if err := service.NotifyUser(1, "Привет! Это тестовое сообщение."); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Отправляем по email
	fmt.Println("--- NotifyByEmail ---")
	if err := service.NotifyByEmail("bob@example.com", "Важное уведомление!"); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Broadcast
	fmt.Println("--- BroadcastToUsers ---")
	sent, failed := service.BroadcastToUsers([]int{1, 2, 99}, "Массовая рассылка")
	fmt.Printf("Отправлено: %d, Неудачно: %d\n", sent, failed)

	fmt.Println("\nТеперь напиши тесты с моками в main_test.go!")
}
