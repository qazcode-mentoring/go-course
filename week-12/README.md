# Неделя 12: Generics (Дженерики)

## Теория

**Ссылка:** https://www.gocat.dev/tour/

Пройди уроки 84-89:
84. Параметры типа
85. Типы-ограничения (constraints)
86. Встроенные ограничения
87. Упражнение: дженерики

## Что такое дженерики?

Дженерики (generics) - это возможность писать код, который работает с разными типами данных, без дублирования логики. Они появились в Go 1.18 и стали важным инструментом для создания переиспользуемого кода.

### Синтаксис

```go
// Функция с параметром типа T
func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Использование
fmt.Println(Min(1, 2))       // int
fmt.Println(Min(1.5, 2.5))   // float64
fmt.Println(Min("a", "b"))   // string
```

### Ограничения типов (constraints)

```go
// any - любой тип (alias для interface{})
func Print[T any](v T) { fmt.Println(v) }

// comparable - типы, которые можно сравнивать через == и !=
func Contains[T comparable](slice []T, val T) bool { ... }

// cmp.Ordered - типы с операторами <, >, <=, >=
func Max[T cmp.Ordered](a, b T) T { ... }

// Собственное ограничение
type Number interface {
    int | int32 | int64 | float32 | float64
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 12.1 | [generic-min-max](./homework/12.1-generic-min-max/) | Min/Max функции для ordered типов |
| 12.2 | [generic-stack](./homework/12.2-generic-stack/) | Универсальный стек Stack[T] |
| 12.3 | [generic-filter-map](./homework/12.3-generic-filter-map/) | Filter и Map для слайсов |
| 12.4 | [generic-cache](./homework/12.4-generic-cache/) | Кэш Cache[K, V] с TTL |

## Вопросы для самопроверки

Создай файл `answers-12.txt` и напиши ответы на вопросы:

1. Чем дженерики отличаются от использования `interface{}`? Какие преимущества они дают?
2. Объясни разницу между ограничениями `any`, `comparable` и `cmp.Ordered`.
3. Как создать собственное ограничение типа (type constraint)?
4. Можно ли использовать дженерики с методами структур? Приведи пример.
5. Когда стоит использовать дженерики, а когда лучше обойтись без них?

## Дополнительные материалы

- [Go Blog: An Introduction To Generics](https://go.dev/blog/intro-generics)
- [Go Blog: When To Use Generics](https://go.dev/blog/when-generics)
- [Go by Example: Generics](https://gobyexample.com/generics)
- [Tutorial: Getting started with generics](https://go.dev/doc/tutorial/generics)
- [Пакет cmp](https://pkg.go.dev/cmp)
- [Пакет slices](https://pkg.go.dev/slices)
- [Пакет maps](https://pkg.go.dev/maps)
