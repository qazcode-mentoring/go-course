# ДЗ 16.3: Покрытие кода и Subtests

## Цель
Научиться анализировать покрытие кода тестами, писать тесты для достижения высокого покрытия и эффективно использовать subtests для группировки тестов.

## Что нужно сделать

В файле `main.go` реализован валидатор данных пользователя с несколькими методами. Твоя задача — написать тесты с покрытием не менее 80%.

### Структура и методы

```go
type User struct {
    Username string
    Email    string
    Age      int
    Password string
}

type Validator struct{}

func (v *Validator) ValidateUsername(username string) error
func (v *Validator) ValidateEmail(email string) error
func (v *Validator) ValidateAge(age int) error
func (v *Validator) ValidatePassword(password string) error
func (v *Validator) ValidateUser(user User) []error
```

### Правила валидации

1. **Username**:
   - Длина от 3 до 20 символов
   - Только буквы, цифры и подчёркивание
   - Не может начинаться с цифры

2. **Email**:
   - Содержит ровно один символ @
   - До @ минимум 1 символ
   - После @ минимум 3 символа и точка

3. **Age**:
   - От 0 до 150

4. **Password**:
   - Минимум 8 символов
   - Содержит хотя бы одну цифру
   - Содержит хотя бы одну заглавную букву
   - Содержит хотя бы одну строчную букву

### Организация тестов с Subtests

Используй `t.Run` для группировки связанных тестов:

```go
func TestValidator_ValidateUsername(t *testing.T) {
    v := &Validator{}

    t.Run("valid usernames", func(t *testing.T) {
        validNames := []string{"user", "user_123", "User_Name"}
        for _, name := range validNames {
            t.Run(name, func(t *testing.T) {
                if err := v.ValidateUsername(name); err != nil {
                    t.Errorf("expected %q to be valid, got error: %v", name, err)
                }
            })
        }
    })

    t.Run("invalid usernames", func(t *testing.T) {
        tests := []struct {
            name   string
            input  string
            reason string
        }{
            {"too short", "ab", "менее 3 символов"},
            {"starts with digit", "1user", "начинается с цифры"},
            {"invalid chars", "user@name", "содержит @"},
        }

        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                err := v.ValidateUsername(tt.input)
                if err == nil {
                    t.Errorf("expected error for %q (%s)", tt.input, tt.reason)
                }
            })
        }
    })
}
```

## Работа с покрытием

```bash
# Запустить тесты с показом процента покрытия
go test -cover

# Создать профиль покрытия
go test -coverprofile=coverage.out

# Посмотреть покрытие по функциям
go tool cover -func=coverage.out

# Открыть визуальный отчёт в браузере
go tool cover -html=coverage.out
```

Визуальный отчёт покажет:
- Зелёным — покрытые строки
- Красным — непокрытые строки
- Серым — строки без исполняемого кода

## Тесты, которые нужно написать

1. **TestValidator_ValidateUsername** — группы:
   - valid usernames (минимум 3 случая)
   - too short (1-2 символа)
   - too long (более 20 символов)
   - starts with digit
   - invalid characters

2. **TestValidator_ValidateEmail** — группы:
   - valid emails
   - missing @
   - multiple @
   - empty before @
   - invalid after @

3. **TestValidator_ValidateAge** — группы:
   - valid ages
   - negative age
   - too old (>150)
   - edge cases (0, 150)

4. **TestValidator_ValidatePassword** — группы:
   - valid passwords
   - too short
   - missing digit
   - missing uppercase
   - missing lowercase

5. **TestValidator_ValidateUser** — проверь:
   - Полностью валидный User (нет ошибок)
   - User с одной ошибкой
   - User с несколькими ошибками (все должны вернуться)

## Подсказки

- Используй вложенные `t.Run` для логической группировки
- Для проверки конкретной ошибки используй `errors.Is`
- Чтобы покрыть все ветки if/else, нужны тесты для каждого условия
- Обрати внимание на граничные случаи (edge cases)

## Критерии выполнения

- [ ] Покрытие кода не менее 80% (`go test -cover`)
- [ ] Тесты организованы с помощью subtests (`t.Run`)
- [ ] Каждая функция валидации имеет отдельную тестовую функцию
- [ ] Тесты покрывают успешные и неуспешные случаи
- [ ] Тесты покрывают граничные случаи
- [ ] `TestValidator_ValidateUser` проверяет сбор множественных ошибок
