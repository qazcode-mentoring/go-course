package main

import (
	"errors"
	"testing"
)

// SentMessage записывает информацию об отправленном сообщении
type SentMessage struct {
	To      string
	Subject string
	Body    string
}

// MockUserRepository — мок для UserRepository
type MockUserRepository struct {
	users        map[int]*User
	usersByEmail map[string]*User
	err          error // ошибка, которую будет возвращать мок
}

// NewMockUserRepository создаёт мок репозитория
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:        make(map[int]*User),
		usersByEmail: make(map[string]*User),
	}
}

// AddUser добавляет пользователя в мок
func (m *MockUserRepository) AddUser(user *User) {
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
}

// SetError устанавливает ошибку, которую мок будет возвращать
func (m *MockUserRepository) SetError(err error) {
	m.err = err
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	// TODO: Реализуй метод
	// Если m.err != nil, верни nil, m.err
	// Иначе ищи пользователя в m.users
	// Если не найден, верни nil, ErrUserNotFound
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
	// TODO: Реализуй метод
	// Аналогично GetByID, но ищи в m.usersByEmail
	return nil, nil
}

// MockMessageSender — мок для MessageSender
type MockMessageSender struct {
	sentMessages []SentMessage
	err          error // ошибка при отправке
}

// NewMockMessageSender создаёт мок отправщика
func NewMockMessageSender() *MockMessageSender {
	return &MockMessageSender{
		sentMessages: make([]SentMessage, 0),
	}
}

// SetError устанавливает ошибку, которую мок будет возвращать
func (m *MockMessageSender) SetError(err error) {
	m.err = err
}

// GetSentMessages возвращает список отправленных сообщений
func (m *MockMessageSender) GetSentMessages() []SentMessage {
	return m.sentMessages
}

func (m *MockMessageSender) Send(to, subject, body string) error {
	// TODO: Реализуй метод
	// Если m.err != nil, верни m.err
	// Иначе добавь сообщение в m.sentMessages и верни nil
	return nil
}

// === Тесты для NotifyUser ===

func TestNotificationService_NotifyUser_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Alice", Email: "alice@example.com"})

	mockSender := NewMockMessageSender()

	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: Вызови service.NotifyUser(1, "Hello!")

	// Assert
	// TODO: Проверь, что ошибки нет
	// TODO: Проверь, что было отправлено 1 сообщение
	// TODO: Проверь, что сообщение отправлено на правильный email

	_ = service
}

func TestNotificationService_NotifyUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository() // пустой репозиторий
	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: Вызови service.NotifyUser для несуществующего пользователя

	// Assert
	// TODO: Проверь, что вернулась ошибка
	// TODO: Проверь, что ошибка содержит ErrUserNotFound (используй errors.Is)
	// TODO: Проверь, что сообщения НЕ были отправлены

	_ = service
}

func TestNotificationService_NotifyUser_SendError(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Bob", Email: "bob@example.com"})

	mockSender := NewMockMessageSender()
	mockSender.SetError(errors.New("network error"))

	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: Вызови service.NotifyUser(1, "Test")

	// Assert
	// TODO: Проверь, что вернулась ошибка
	// TODO: Пользователь найден, но отправка не удалась

	_ = service
}

func TestNotificationService_NotifyUser_CorrectMessageContent(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Charlie", Email: "charlie@example.com"})

	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	message := "Important notification!"
	// TODO: Вызови service.NotifyUser(1, message)

	// Assert
	// TODO: Проверь содержимое отправленного сообщения
	// - To должен быть "charlie@example.com"
	// - Subject должен содержать имя пользователя
	// - Body должен быть равен message

	_ = service
	_ = message
}

// === Тесты для NotifyByEmail ===

func TestNotificationService_NotifyByEmail_Success(t *testing.T) {
	// TODO: Arrange — создай моки с пользователем

	// TODO: Act — вызови NotifyByEmail

	// TODO: Assert — проверь успешную отправку
}

func TestNotificationService_NotifyByEmail_UserNotFound(t *testing.T) {
	// TODO: Arrange — создай моки БЕЗ пользователя с искомым email

	// TODO: Act — вызови NotifyByEmail с несуществующим email

	// TODO: Assert — проверь, что вернулась ошибка
}

func TestNotificationService_NotifyByEmail_SendError(t *testing.T) {
	// TODO: Arrange — создай моки, настрой ошибку в sender

	// TODO: Act — вызови NotifyByEmail

	// TODO: Assert — проверь, что вернулась ошибка отправки
}

// === Тесты для BroadcastToUsers ===

func TestNotificationService_BroadcastToUsers_AllSuccess(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "User1", Email: "user1@example.com"})
	mockRepo.AddUser(&User{ID: 2, Name: "User2", Email: "user2@example.com"})
	mockRepo.AddUser(&User{ID: 3, Name: "User3", Email: "user3@example.com"})

	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: sent, failed := service.BroadcastToUsers([]int{1, 2, 3}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 3, failed == 0
	// TODO: Проверь, что отправлено 3 сообщения

	_ = service
}

func TestNotificationService_BroadcastToUsers_SomeNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "User1", Email: "user1@example.com"})
	// ID 2 и 3 не существуют

	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: sent, failed := service.BroadcastToUsers([]int{1, 2, 3}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 1, failed == 2

	_ = service
}

func TestNotificationService_BroadcastToUsers_SendErrors(t *testing.T) {
	// TODO: Этот тест сложнее, так как нужно симулировать ошибку только для некоторых отправок
	// Можно упростить: настроить ошибку для всех отправок

	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "User1", Email: "user1@example.com"})
	mockRepo.AddUser(&User{ID: 2, Name: "User2", Email: "user2@example.com"})

	mockSender := NewMockMessageSender()
	mockSender.SetError(errors.New("send failed"))

	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: sent, failed := service.BroadcastToUsers([]int{1, 2}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 0, failed == 2

	_ = service
}

func TestNotificationService_BroadcastToUsers_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: sent, failed := service.BroadcastToUsers([]int{}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 0, failed == 0
	// TODO: Проверь, что сообщения не отправлялись

	_ = service
}
