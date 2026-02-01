# Неделя 16: Тестирование в Go

## Теория

Тестирование — неотъемлемая часть разработки на Go. В стандартной библиотеке есть мощный пакет `testing`, который покрывает большинство потребностей. В этом модуле мы изучим основы тестирования, table-driven tests, subtests, покрытие кода и мокирование через интерфейсы.

### Пакет testing: основы

Go имеет встроенную поддержку тестов. Тестовые файлы должны:
- Называться `*_test.go`
- Находиться в том же пакете (или `package_test` для black-box тестирования)
- Содержать функции с сигнатурой `func TestXxx(t *testing.T)`

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

#### Запуск тестов

```bash
# Запустить все тесты в текущем пакете
go test

# Запустить с подробным выводом
go test -v

# Запустить конкретный тест
go test -run TestAdd

# Запустить тесты во всех подпакетах
go test ./...

# Запустить с покрытием
go test -cover

# Сгенерировать отчёт о покрытии
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Методы *testing.T

```go
func TestExample(t *testing.T) {
    // Логирование (видно только с -v или при провале)
    t.Log("информационное сообщение")
    t.Logf("форматированное: %d", 42)

    // Провал теста с продолжением
    t.Error("тест провален, но продолжаем")
    t.Errorf("провал: получили %d, ожидали %d", got, want)

    // Немедленный провал и остановка теста
    t.Fatal("критическая ошибка, останавливаемся")
    t.Fatalf("критическая: %v", err)

    // Пропуск теста
    t.Skip("пропускаем этот тест")
    t.Skipf("пропускаем: %s", reason)

    // Пометить тест как провальный, но продолжить
    t.Fail()      // помечает провал
    t.FailNow()   // помечает провал и останавливает
}
```

### Table-Driven Tests

Паттерн table-driven tests — это идиоматический способ тестирования в Go. Вместо написания множества похожих тестов, мы создаём таблицу с тестовыми случаями.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zeros", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

#### Преимущества table-driven tests

1. **Легко добавлять новые случаи** — просто добавь строку в таблицу
2. **Единообразие** — все случаи тестируются одинаково
3. **Читаемость** — легко понять, что тестируется
4. **Изоляция** — каждый случай запускается как subtest

### Subtests (t.Run)

`t.Run` позволяет создавать вложенные тесты (subtests):

```go
func TestMath(t *testing.T) {
    t.Run("Addition", func(t *testing.T) {
        if Add(2, 2) != 4 {
            t.Error("2 + 2 should equal 4")
        }
    })

    t.Run("Subtraction", func(t *testing.T) {
        if Sub(5, 3) != 2 {
            t.Error("5 - 3 should equal 2")
        }
    })
}
```

Запуск конкретного subtest:
```bash
go test -run TestMath/Addition
```

### Покрытие кода (Code Coverage)

```bash
# Показать процент покрытия
go test -cover

# Создать профиль покрытия
go test -coverprofile=coverage.out

# Просмотреть в браузере
go tool cover -html=coverage.out

# Показать покрытие по функциям
go tool cover -func=coverage.out
```

Пример вывода:
```
github.com/user/project/math.go:5:   Add      100.0%
github.com/user/project/math.go:10:  Divide   75.0%
total:                               (statements) 87.5%
```

### Мокирование через интерфейсы

В Go мокирование обычно делается через интерфейсы. Это один из ключевых паттернов для написания тестируемого кода.

#### Принцип: программируй на интерфейсах

```go
// repository.go
type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
}

// service.go
type UserService struct {
    repo UserRepository // зависимость через интерфейс
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetByID(id)
}
```

#### Создание мока

```go
// service_test.go
type MockUserRepository struct {
    users map[int]*User
    err   error
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func (m *MockUserRepository) Save(user *User) error {
    if m.err != nil {
        return m.err
    }
    m.users[user.ID] = user
    return nil
}

func TestUserService_GetUser(t *testing.T) {
    // Arrange: создаём мок с тестовыми данными
    mock := &MockUserRepository{
        users: map[int]*User{
            1: {ID: 1, Name: "Alice"},
        },
    }
    service := NewUserService(mock)

    // Act: вызываем тестируемый метод
    user, err := service.GetUser(1)

    // Assert: проверяем результат
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("got name %q, want %q", user.Name, "Alice")
    }
}
```

### Библиотека testify

Хотя стандартный пакет `testing` достаточен для большинства случаев, библиотека `testify` предоставляет удобные функции:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWithTestify(t *testing.T) {
    // assert — продолжает выполнение при провале
    assert.Equal(t, 5, Add(2, 3), "сумма должна быть 5")
    assert.NotNil(t, user, "пользователь не должен быть nil")
    assert.True(t, isValid, "должно быть валидным")
    assert.Error(t, err, "ожидается ошибка")
    assert.NoError(t, err, "ошибки быть не должно")

    // require — останавливает тест при провале (как t.Fatal)
    require.NoError(t, err, "критическая ошибка")
    require.NotNil(t, user) // без user дальше нет смысла
}
```

Установка:
```bash
go get github.com/stretchr/testify
```

### Паттерн AAA (Arrange-Act-Assert)

Структурируй тесты по паттерну AAA:

```go
func TestCalculator_Divide(t *testing.T) {
    // Arrange: подготовка
    calc := NewCalculator()
    a, b := 10, 2

    // Act: действие
    result, err := calc.Divide(a, b)

    // Assert: проверка
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(%d, %d) = %d; want %d", a, b, result, 5)
    }
}
```

### Тестирование ошибок

```go
func TestDivide_ByZero(t *testing.T) {
    _, err := Divide(10, 0)

    if err == nil {
        t.Fatal("expected error for division by zero")
    }

    // Проверка конкретной ошибки
    if !errors.Is(err, ErrDivisionByZero) {
        t.Errorf("got error %v, want %v", err, ErrDivisionByZero)
    }
}
```

### Хелперы (t.Helper)

Используй `t.Helper()` для вспомогательных функций:

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper() // сообщает, что это хелпер
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}

func TestSomething(t *testing.T) {
    result := Compute()
    assertEqual(t, result, 42) // ошибка укажет на эту строку
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 16.1 | [basic-testing](./homework/16.1-basic-testing/) | Первые тесты: Test* функции, t.Error/t.Fatal |
| 16.2 | [table-driven](./homework/16.2-table-driven/) | Table-driven tests для Calculator |
| 16.3 | [test-coverage](./homework/16.3-test-coverage/) | Тесты с покрытием, subtests (t.Run) |
| 16.4 | [mocks-interfaces](./homework/16.4-mocks-interfaces/) | Мокирование через интерфейсы |

## Вопросы для самопроверки

Создай файл `answers-16.txt` и напиши ответы на вопросы:

1. В чём разница между `t.Error()` и `t.Fatal()`? Когда следует использовать каждый из них?

2. Почему table-driven tests считаются идиоматичным подходом в Go? Какие преимущества они дают по сравнению с отдельными тестовыми функциями?

3. Что такое покрытие кода (code coverage) и почему 100% покрытие не гарантирует отсутствие багов?

4. Как интерфейсы помогают в тестировании? Объясни принцип "accept interfaces, return structs" в контексте мокирования.

5. Для чего нужен метод `t.Helper()` и как он влияет на вывод ошибок в тестах?

## Дополнительные материалы

- [Go Blog: Using subtests and sub-benchmarks](https://go.dev/blog/subtests)
- [Go Doc: testing package](https://pkg.go.dev/testing)
- [Go by Example: Testing](https://gobyexample.com/testing)
- [Go by Example: Table-driven tests](https://gobyexample.com/testing-and-benchmarking)
- [Dave Cheney: Writing table driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testify on GitHub](https://github.com/stretchr/testify)
- [Go Wiki: Table-driven tests](https://go.dev/wiki/TableDrivenTests)
- [Practical Go: Test fixtures](https://peter.bourgon.org/go-best-practices-2016/#testing)
