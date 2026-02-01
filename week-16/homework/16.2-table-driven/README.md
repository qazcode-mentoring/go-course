# ДЗ 16.2: Table-Driven Tests

## Цель
Освоить паттерн table-driven tests — идиоматический способ тестирования в Go. Научиться структурировать тестовые случаи в виде таблицы и использовать `t.Run` для создания subtests.

## Что нужно сделать

В файле `main.go` реализован тип `Calculator` с методами для математических операций. Твоя задача — написать table-driven тесты для всех методов.

### Методы Calculator

1. **Add(a, b float64) float64** — сложение
2. **Subtract(a, b float64) float64** — вычитание
3. **Multiply(a, b float64) float64** — умножение
4. **Divide(a, b float64) (float64, error)** — деление (с обработкой деления на ноль)
5. **Power(base, exp float64) float64** — возведение в степень
6. **Sqrt(x float64) (float64, error)** — квадратный корень (с обработкой отрицательных чисел)

### Структура table-driven теста

```go
func TestCalculator_Add(t *testing.T) {
    calc := NewCalculator()

    tests := []struct {
        name     string   // название тестового случая
        a, b     float64  // входные данные
        expected float64  // ожидаемый результат
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed signs", -2, 5, 3},
        {"zeros", 0, 0, 0},
        {"decimals", 1.5, 2.5, 4.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calc.Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%.2f, %.2f) = %.2f; want %.2f",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Тесты, которые нужно написать

1. **TestCalculator_Add** — минимум 5 случаев
2. **TestCalculator_Subtract** — минимум 4 случая
3. **TestCalculator_Multiply** — минимум 5 случаев (включая умножение на 0 и 1)
4. **TestCalculator_Divide** — минимум 5 случаев (включая деление на 0, проверку ошибки)
5. **TestCalculator_Power** — минимум 4 случая (включая степень 0 и отрицательную степень)
6. **TestCalculator_Sqrt** — минимум 4 случая (включая корень из 0 и отрицательного числа)

### Тестирование ошибок в table-driven tests

```go
func TestCalculator_Divide(t *testing.T) {
    calc := NewCalculator()

    tests := []struct {
        name      string
        a, b      float64
        expected  float64
        expectErr bool  // ожидаем ли ошибку
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.Divide(tt.a, tt.b)

            if tt.expectErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return // не проверяем результат при ошибке
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if result != tt.expected {
                t.Errorf("got %.2f, want %.2f", result, tt.expected)
            }
        })
    }
}
```

## Запуск тестов

```bash
# Запустить все тесты
go test -v

# Запустить конкретный subtest
go test -v -run "TestCalculator_Divide/division_by_zero"

# Запустить все тесты Divide
go test -v -run TestCalculator_Divide
```

## Подсказки

- Используй `float64` осторожно: сравнивай с небольшой погрешностью для дробных чисел
- Для сравнения float можно использовать функцию: `math.Abs(got-want) < 0.0001`
- Имена тестовых случаев должны быть понятными и описывать что тестируется
- `t.Run` создаёт subtest, который можно запускать отдельно

## Критерии выполнения

- [ ] Написаны table-driven тесты для всех 6 методов Calculator
- [ ] Каждый тест содержит минимум указанное количество случаев
- [ ] Используется `t.Run` для каждого тестового случая
- [ ] Тесты методов с ошибками (Divide, Sqrt) проверяют и успешные случаи, и ошибки
- [ ] Все тесты проходят (`go test` без ошибок)
- [ ] Имена тестовых случаев понятные и описательные
