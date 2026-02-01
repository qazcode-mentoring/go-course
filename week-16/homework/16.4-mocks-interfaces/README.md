# ДЗ 16.4: Мокирование через интерфейсы

## Цель
Научиться писать тестируемый код с использованием интерфейсов и создавать mock-объекты для изоляции тестов от внешних зависимостей.

## Что нужно сделать

В файле `main.go` реализован сервис уведомлений `NotificationService`, который зависит от внешних сервисов: репозитория пользователей и отправщика сообщений. Твоя задача — написать тесты с использованием mock-объектов.

### Структура кода

```go
// Интерфейсы (контракты)
type UserRepository interface {
    GetByID(id int) (*User, error)
    GetByEmail(email string) (*User, error)
}

type MessageSender interface {
    Send(to, subject, body string) error
}

// Сервис, который нужно тестировать
type NotificationService struct {
    users  UserRepository
    sender MessageSender
}

func (s *NotificationService) NotifyUser(userID int, message string) error
func (s *NotificationService) NotifyByEmail(email, message string) error
func (s *NotificationService) BroadcastToUsers(userIDs []int, message string) (sent int, failed int)
```

### Создание Mock-объектов

```go
// MockUserRepository — мок для UserRepository
type MockUserRepository struct {
    users map[int]*User
    err   error // ошибка, которую будет возвращать мок
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    user, ok := m.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
}

// MockMessageSender — мок для MessageSender
type MockMessageSender struct {
    sentMessages []SentMessage // записываем что отправили
    err          error         // ошибка при отправке
}

type SentMessage struct {
    To      string
    Subject string
    Body    string
}

func (m *MockMessageSender) Send(to, subject, body string) error {
    if m.err != nil {
        return m.err
    }
    m.sentMessages = append(m.sentMessages, SentMessage{to, subject, body})
    return nil
}
```

### Тесты, которые нужно написать

1. **TestNotificationService_NotifyUser**:
   - Успешная отправка уведомления существующему пользователю
   - Пользователь не найден (UserRepository возвращает ошибку)
   - Ошибка отправки сообщения (MessageSender возвращает ошибку)
   - Проверь, что сообщение отправлено правильному получателю

2. **TestNotificationService_NotifyByEmail**:
   - Успешная отправка по email
   - Пользователь с таким email не найден
   - Ошибка отправки

3. **TestNotificationService_BroadcastToUsers**:
   - Успешная отправка всем пользователям
   - Часть пользователей не найдена (проверь счётчики sent/failed)
   - Ошибка отправки для некоторых (проверь счётчики)
   - Пустой список пользователей

## Паттерн AAA в тестах

Структурируй тесты по паттерну Arrange-Act-Assert:

```go
func TestNotificationService_NotifyUser_Success(t *testing.T) {
    // Arrange: подготовка
    mockRepo := &MockUserRepository{
        users: map[int]*User{
            1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
        },
    }
    mockSender := &MockMessageSender{}

    service := NewNotificationService(mockRepo, mockSender)

    // Act: действие
    err := service.NotifyUser(1, "Hello!")

    // Assert: проверка
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(mockSender.sentMessages) != 1 {
        t.Fatalf("expected 1 message sent, got %d", len(mockSender.sentMessages))
    }

    msg := mockSender.sentMessages[0]
    if msg.To != "alice@example.com" {
        t.Errorf("message sent to %q, want %q", msg.To, "alice@example.com")
    }
}
```

## Проверка вызовов Mock

Моки позволяют проверить:
- Был ли вызван метод
- С какими аргументами
- Сколько раз

```go
// После теста проверяем, что Send был вызван правильно
if len(mockSender.sentMessages) != 1 {
    t.Errorf("expected Send to be called once, called %d times",
        len(mockSender.sentMessages))
}

if mockSender.sentMessages[0].To != expectedEmail {
    t.Errorf("Send called with wrong email")
}
```

## Подсказки

- Моки создаются с нужным состоянием перед каждым тестом
- Используй поле `err` в моке для симуляции ошибок
- Записывай вызовы в слайс для последующей проверки
- Каждый тест должен быть независимым — создавай новые моки

## Критерии выполнения

- [ ] Созданы `MockUserRepository` и `MockMessageSender`
- [ ] Написаны тесты для `NotifyUser` (минимум 4 случая)
- [ ] Написаны тесты для `NotifyByEmail` (минимум 3 случая)
- [ ] Написаны тесты для `BroadcastToUsers` (минимум 4 случая)
- [ ] Тесты проверяют не только результат, но и что моки были вызваны правильно
- [ ] Тесты покрывают успешные сценарии и обработку ошибок
- [ ] Все тесты проходят (`go test -v`)
