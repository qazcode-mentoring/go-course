# ДЗ 12.2: Generic Stack

## Цель
Научиться создавать обобщенные структуры данных с методами.

## Что нужно сделать

Реализовать универсальный стек `Stack[T]` со следующими методами:

1. `Push(value T)` - добавить элемент на вершину стека
2. `Pop() (T, bool)` - извлечь элемент с вершины стека
3. `Peek() (T, bool)` - посмотреть элемент на вершине без извлечения
4. `Len() int` - вернуть количество элементов в стеке
5. `IsEmpty() bool` - проверить, пуст ли стек

## Сигнатуры

```go
// Stack - универсальный стек
type Stack[T any] struct {
    items []T
}

// NewStack создает новый пустой стек
func NewStack[T any]() *Stack[T]

// Push добавляет элемент на вершину стека
func (s *Stack[T]) Push(value T)

// Pop извлекает элемент с вершины стека
// Возвращает (value, true) если стек не пустой
// Возвращает (zero, false) если стек пустой
func (s *Stack[T]) Pop() (T, bool)

// Peek возвращает элемент с вершины без извлечения
// Возвращает (value, true) если стек не пустой
// Возвращает (zero, false) если стек пустой
func (s *Stack[T]) Peek() (T, bool)

// Len возвращает количество элементов в стеке
func (s *Stack[T]) Len() int

// IsEmpty проверяет, пуст ли стек
func (s *Stack[T]) IsEmpty() bool
```

## Принцип работы стека

Стек работает по принципу LIFO (Last In, First Out - последний вошел, первый вышел):

```
Push(1): [1]
Push(2): [1, 2]
Push(3): [1, 2, 3]
Pop():   [1, 2] -> вернул 3
Peek():  [1, 2] -> вернул 2 (не удаляя)
Pop():   [1] -> вернул 2
```

## Пример использования

```go
func main() {
    // Стек целых чисел
    intStack := NewStack[int]()
    intStack.Push(1)
    intStack.Push(2)
    intStack.Push(3)

    fmt.Println(intStack.Len())      // 3
    fmt.Println(intStack.Peek())     // 3, true
    fmt.Println(intStack.Pop())      // 3, true
    fmt.Println(intStack.Pop())      // 2, true
    fmt.Println(intStack.Len())      // 1
    fmt.Println(intStack.IsEmpty())  // false

    // Стек строк
    strStack := NewStack[string]()
    strStack.Push("hello")
    strStack.Push("world")

    if val, ok := strStack.Pop(); ok {
        fmt.Println(val) // "world"
    }

    // Проверка пустого стека
    emptyStack := NewStack[float64]()
    if val, ok := emptyStack.Pop(); !ok {
        fmt.Println("Стек пуст!")
    }
}
```

## Подсказки

- Используй слайс для хранения элементов
- `Push`: добавление в конец слайса `s.items = append(s.items, value)`
- `Pop`: извлечение последнего элемента и уменьшение слайса
- Не забудь проверять границы слайса в Pop и Peek
- Используй `var zero T` для возврата нулевого значения типа

## Критерии выполнения

- [ ] Stack работает с любыми типами (int, string, структуры)
- [ ] Push добавляет элементы на вершину
- [ ] Pop извлекает элементы в порядке LIFO
- [ ] Peek не изменяет стек
- [ ] Pop и Peek корректно обрабатывают пустой стек
- [ ] Len и IsEmpty возвращают корректные значения
- [ ] Код компилируется без ошибок
