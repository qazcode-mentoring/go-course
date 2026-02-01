# ДЗ 12.1: Generic Min/Max

## Цель
Научиться создавать обобщенные функции с ограничениями типов (type constraints).

## Что нужно сделать

1. Реализовать функцию `Min[T]` - возвращает минимальное из двух значений
2. Реализовать функцию `Max[T]` - возвращает максимальное из двух значений
3. Реализовать функцию `MinSlice[T]` - возвращает минимальный элемент слайса
4. Реализовать функцию `MaxSlice[T]` - возвращает максимальный элемент слайса
5. Реализовать функцию `Clamp[T]` - ограничивает значение в заданном диапазоне

## Сигнатуры

```go
// Min возвращает меньшее из двух значений
func Min[T cmp.Ordered](a, b T) T

// Max возвращает большее из двух значений
func Max[T cmp.Ordered](a, b T) T

// MinSlice возвращает минимальный элемент слайса и true
// Если слайс пустой, возвращает нулевое значение и false
func MinSlice[T cmp.Ordered](slice []T) (T, bool)

// MaxSlice возвращает максимальный элемент слайса и true
// Если слайс пустой, возвращает нулевое значение и false
func MaxSlice[T cmp.Ordered](slice []T) (T, bool)

// Clamp ограничивает значение v в диапазоне [min, max]
// Если v < min, возвращает min
// Если v > max, возвращает max
// Иначе возвращает v
func Clamp[T cmp.Ordered](v, min, max T) T
```

## Ограничение cmp.Ordered

```go
import "cmp"

// cmp.Ordered включает все типы, поддерживающие операторы сравнения:
// int, int8, int16, int32, int64
// uint, uint8, uint16, uint32, uint64, uintptr
// float32, float64
// string
```

## Пример использования

```go
func main() {
    // Min/Max для разных типов
    fmt.Println(Min(3, 5))           // 3
    fmt.Println(Max(3.14, 2.71))     // 3.14
    fmt.Println(Min("apple", "banana")) // "apple"

    // MinSlice/MaxSlice
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
    if min, ok := MinSlice(nums); ok {
        fmt.Println("Min:", min) // Min: 1
    }
    if max, ok := MaxSlice(nums); ok {
        fmt.Println("Max:", max) // Max: 9
    }

    // Пустой слайс
    empty := []int{}
    if _, ok := MinSlice(empty); !ok {
        fmt.Println("Слайс пустой!")
    }

    // Clamp
    fmt.Println(Clamp(5, 0, 10))   // 5 (в диапазоне)
    fmt.Println(Clamp(-5, 0, 10))  // 0 (меньше min)
    fmt.Println(Clamp(15, 0, 10))  // 10 (больше max)
}
```

## Подсказки

- Импортируй пакет `cmp` для использования ограничения `cmp.Ordered`
- Для пустого слайса используй `var zero T` для получения нулевого значения типа
- Функция `Clamp` - это комбинация `Min` и `Max`

## Критерии выполнения

- [ ] Функции Min и Max работают с int, float64, string
- [ ] MinSlice и MaxSlice корректно обрабатывают пустой слайс
- [ ] Clamp правильно ограничивает значения в диапазоне
- [ ] Код компилируется без ошибок
- [ ] Все тестовые примеры в main() выводят ожидаемые результаты
