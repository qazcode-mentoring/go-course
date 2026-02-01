# ДЗ 12.3: Generic Filter и Map

## Цель
Научиться создавать обобщенные функции высшего порядка для работы со слайсами.

## Что нужно сделать

Реализовать функциональные утилиты для работы со слайсами:

1. `Filter[T]` - фильтрация элементов по предикату
2. `Map[T, R]` - преобразование элементов
3. `Reduce[T, R]` - свертка слайса в одно значение
4. `Find[T]` - поиск первого подходящего элемента
5. `Any[T]` и `All[T]` - проверка условий

## Сигнатуры

```go
// Filter возвращает новый слайс с элементами, для которых predicate вернул true
func Filter[T any](slice []T, predicate func(T) bool) []T

// Map преобразует каждый элемент слайса с помощью функции transform
// T - исходный тип, R - результирующий тип
func Map[T, R any](slice []T, transform func(T) R) []R

// Reduce сворачивает слайс в одно значение
// initial - начальное значение аккумулятора
// reducer - функция, принимающая аккумулятор и элемент, возвращает новый аккумулятор
func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R

// Find ищет первый элемент, для которого predicate вернул true
// Возвращает (element, true) если найден
// Возвращает (zero, false) если не найден
func Find[T any](slice []T, predicate func(T) bool) (T, bool)

// Any возвращает true, если хотя бы один элемент удовлетворяет predicate
func Any[T any](slice []T, predicate func(T) bool) bool

// All возвращает true, если все элементы удовлетворяют predicate
// Для пустого слайса возвращает true
func All[T any](slice []T, predicate func(T) bool) bool
```

## Пример использования

```go
func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Filter: только четные числа
    evens := Filter(nums, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2, 4, 6, 8, 10]

    // Map: возвести в квадрат
    squares := Map(nums, func(n int) int {
        return n * n
    })
    fmt.Println(squares) // [1, 4, 9, 16, 25, 36, 49, 64, 81, 100]

    // Map с изменением типа: числа в строки
    strs := Map(nums, func(n int) string {
        return fmt.Sprintf("#%d", n)
    })
    fmt.Println(strs) // ["#1", "#2", ..., "#10"]

    // Reduce: сумма всех чисел
    sum := Reduce(nums, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println(sum) // 55

    // Find: найти первое число > 5
    if val, ok := Find(nums, func(n int) bool { return n > 5 }); ok {
        fmt.Println(val) // 6
    }

    // Any: есть ли число > 5?
    fmt.Println(Any(nums, func(n int) bool { return n > 5 })) // true

    // All: все ли числа > 0?
    fmt.Println(All(nums, func(n int) bool { return n > 0 })) // true
}
```

## Подсказки

- `Filter`: создай пустой слайс, добавляй элементы если predicate(element) == true
- `Map`: создай слайс результатов той же длины, заполни преобразованными значениями
- `Reduce`: начни с initial, на каждой итерации обновляй: acc = reducer(acc, element)
- `Any`: верни true при первом совпадении, false если дошел до конца
- `All`: верни false при первом несовпадении, true если дошел до конца

## Практическое применение

Эти функции часто используются при обработке данных:

```go
type User struct {
    Name   string
    Age    int
    Active bool
}

users := []User{...}

// Найти активных пользователей старше 18
adults := Filter(users, func(u User) bool {
    return u.Age >= 18 && u.Active
})

// Получить список имен
names := Map(adults, func(u User) string {
    return u.Name
})

// Посчитать средний возраст
totalAge := Reduce(adults, 0, func(sum int, u User) int {
    return sum + u.Age
})
avgAge := float64(totalAge) / float64(len(adults))
```

## Критерии выполнения

- [ ] Filter корректно фильтрует элементы
- [ ] Map преобразует элементы, включая изменение типа
- [ ] Reduce правильно сворачивает слайс
- [ ] Find находит первый подходящий элемент
- [ ] Any и All корректно проверяют условия
- [ ] Все функции корректно обрабатывают пустой слайс
- [ ] Код компилируется без ошибок
