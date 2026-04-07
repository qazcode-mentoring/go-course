package main

import (
	"errors"
	"strings"
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
	if m.err != nil {
		return nil, m.err
	}

	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
	// TODO: Реализуй метод
	// Аналогично GetByID, но ищи в m.usersByEmail
	if m.err != nil {
		return nil, m.err
	}

	for _, user := range m.usersByEmail {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
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
	if m.err != nil {
		return m.err
	}

	send := SentMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	m.sentMessages = append(m.sentMessages, send)
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
	err := service.NotifyUser(1, "hello")

	// Assert
	// TODO: Проверь, что ошибки нет
	// TODO: Проверь, что было отправлено 1 сообщение
	// TODO: Проверь, что сообщение отправлено на правильный email

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	sent := mockSender.GetSentMessages()
	if len(sent) != 1 {
		t.Fatalf("expected 1 message to be sent, got: %d", len(sent))
	}

	msg := sent[0]
	if msg.To != "alice@example.com" {
		t.Errorf("expected email to be alice@example.com, got: %q", msg.To)
	}
}

func TestNotificationService_NotifyUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository() // пустой репозиторий
	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO: Вызови service.NotifyUser для несуществующего пользователя
	err := service.NotifyUser(2, "some message")
	// Assert
	// TODO: Проверь, что вернулась ошибка
	// TODO: Проверь, что ошибка содержит ErrUserNotFound (используй errors.Is)
	// TODO: Проверь, что сообщения НЕ были отправлены
	if err == nil {
		t.Fatal("expected an error because user does not exist, but got nil")
	}

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	sent := mockSender.GetSentMessages()
	if len(sent) != 0 {
		t.Errorf("expected 0 message to be sent, got: %d", len(sent))
	}
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
	err := service.NotifyUser(1, "Test")
	// Assert
	// TODO: Проверь, что вернулась ошибка
	// TODO: Пользователь найден, но отправка не удалась
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected user to be found, but got %v", err)
	}

	expectedErr := "network error"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error to contain %q, but got: %q", expectedErr, err.Error())
	}
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
	err := service.NotifyUser(1, message)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert
	// TODO: Проверь содержимое отправленного сообщения
	// - To должен быть "charlie@example.com"
	// - Subject должен содержать имя пользователя
	// - Body должен быть равен message
	sent := mockSender.GetSentMessages()

	if len(sent) != 1 {
		t.Fatalf("expected 1 message to be sent, got %d", len(sent))
	}

	msg := sent[0]
	if msg.To != "charlie@example.com" {
		t.Errorf("To: expected %q, got %q", "charlie@example.com", msg.To)
	}

	if msg.Subject != "Уведомление для Charlie" {
		t.Errorf("Subject: expected %q, got %q", "Charlie", msg.Subject)
	}

	if msg.Body != message {
		t.Errorf("Body: expected %q, got %q", message, msg.Body)
	}
}

// === Тесты для NotifyByEmail ===

func TestNotificationService_NotifyByEmail_Success(t *testing.T) {
	// TODO: Arrange — создай моки с пользователем
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Dam", Email: "dkairzhanov@beeline.kz"})

	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)
	// TODO: Act — вызови NotifyByEmail
	message := "Привет"
	err := service.NotifyByEmail("dkairzhanov@beeline.kz", message)
	// TODO: Assert — проверь успешную отправку
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sent := mockSender.GetSentMessages()

	if len(sent) != 1 {
		t.Fatalf("expected 1 message to be sent, got %d", len(sent))
	}

	msg := sent[0]
	if msg.To != "dkairzhanov@beeline.kz" {
		t.Errorf("expected email to be dkairzhanov@beeline.kz, got: %q", msg.To)
	}
}

func TestNotificationService_NotifyByEmail_UserNotFound(t *testing.T) {
	// TODO: Arrange — создай моки БЕЗ пользователя с искомым email
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Islam", Email: "iuzakpai@beeline.kz"})
	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// TODO: Act — вызови NotifyByEmail с несуществующим email
	err := service.NotifyByEmail("damirkairzhanov4@gmail.com", "hi")
	// TODO: Assert — проверь, что вернулась ошибка
	if err == nil {
		t.Fatal("expected no error because email does not exist, but go nil")
	}

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected error ErrUserNotFound, got: %v", err)
	}

	sent := mockSender.GetSentMessages()
	if len(sent) != 0 {
		t.Errorf("expected 0 message to be sent, got: %d", len(sent))
	}
}

func TestNotificationService_NotifyByEmail_SendError(t *testing.T) {
	// TODO: Arrange — создай моки, настрой ошибку в sender
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "Damir", Email: "damir@gmail.com"})

	mockSender := NewMockMessageSender()
	expectedErr := "network error"
	mockSender.SetError(errors.New(expectedErr))

	service := NewNotificationService(mockRepo, mockSender)

	// TODO: Act — вызови NotifyByEmail
	err := service.NotifyByEmail("damir@gmail.com", "hello")
	// TODO: Assert — проверь, что вернулась ошибка отправки
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected user to be found, but got %v", err)
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error contain %q, but got: %q", expectedErr, err.Error())
	}
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
	// TODO:
	sent, failed := service.BroadcastToUsers([]int{1, 2, 3}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 3, failed == 0
	// TODO: Проверь, что отправлено 3 сообщения

	if sent != 3 {
		t.Errorf("expected 3 sent, got: %d", sent)
	}

	if failed != 0 {
		t.Errorf("expected 0 failed, got: %d", failed)
	}

	sentMsg := mockSender.GetSentMessages()
	if len(sentMsg) != 3 {
		t.Errorf("expected 3 messages, got: %d", len(sentMsg))
	}
}

func TestNotificationService_BroadcastToUsers_SomeNotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockRepo.AddUser(&User{ID: 1, Name: "User1", Email: "user1@example.com"})
	// ID 2 и 3 не существуют

	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO:
	sent, failed := service.BroadcastToUsers([]int{1, 2, 3}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 1, failed == 2
	if sent != 1 {
		t.Fatalf("expected sent 1 message but got: %d", sent)
	}

	if failed != 2 {
		t.Errorf("expected 2 failed but got: %d", failed)
	}
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
	// TODO:
	sent, failed := service.BroadcastToUsers([]int{1, 2}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 0, failed == 2
	if sent != 0 {
		t.Fatalf("expected 0 sent message, but got: %d", sent)
	}

	if failed != 2 {
		t.Errorf("expected 2 failed, but got: %d", failed)
	}
}

func TestNotificationService_BroadcastToUsers_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockSender := NewMockMessageSender()
	service := NewNotificationService(mockRepo, mockSender)

	// Act
	// TODO:
	sent, failed := service.BroadcastToUsers([]int{}, "Broadcast!")

	// Assert
	// TODO: Проверь, что sent == 0, failed == 0
	// TODO: Проверь, что сообщения не отправлялись
	if sent != 0 {
		t.Fatalf("expected 0 sent messages, but got: %d", sent)
	}

	if failed != 0 {
		t.Fatalf("expected 0 failed, but got: %d", sent)
	}

	sentMsg := mockSender.GetSentMessages()
	if len(sentMsg) != 0 {
		t.Errorf("expected no messages, but got: %d", len(sentMsg))
	}
}
