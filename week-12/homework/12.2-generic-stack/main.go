package main

import "fmt"

// Stack - универсальный стек для любых типов данных
type Stack[T any] struct {
	items []T
}

// NewStack создает новый пустой стек
func NewStack[T any]() *Stack[T] {
	// TODO: реализуй функцию
	// Создай и верни указатель на новый Stack с пустым слайсом
	return &Stack[T]{items: nil}
}

// Push добавляет элемент на вершину стека
func (s *Stack[T]) Push(value T) {
	// TODO: реализуй метод
	// Добавь value в конец слайса items
	// Используй: s.items = append(s.items, value)
	s.items = append(s.items, value)
}

// Pop извлекает элемент с вершины стека
// Возвращает (value, true) если стек не пустой
// Возвращает (zero, false) если стек пустой
func (s *Stack[T]) Pop() (T, bool) {
	// TODO: реализуй метод
	// 1. Проверь, что стек не пустой
	// 2. Получи последний элемент
	// 3. Удали последний элемент из слайса: s.items = s.items[:len(s.items)-1]
	// 4. Верни элемент и true
	var zero T

	if len(s.items) != 0 {
		last := s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
		return last, true
	}

	return zero, false
}

// Peek возвращает элемент с вершины без извлечения
// Возвращает (value, true) если стек не пустой
// Возвращает (zero, false) если стек пустой
func (s *Stack[T]) Peek() (T, bool) {
	// TODO: реализуй метод
	// 1. Проверь, что стек не пустой
	// 2. Верни последний элемент без удаления
	var zero T

	if len(s.items) != 0 {
		return s.items[len(s.items)-1], true
	}

	return zero, false
}

// Len возвращает количество элементов в стеке
func (s *Stack[T]) Len() int {
	// TODO: реализуй метод
	return len(s.items)
}

// IsEmpty проверяет, пуст ли стек
func (s *Stack[T]) IsEmpty() bool {
	// TODO: реализуй метод
	if len(s.items) != 0 {
		return false
	}
	return true
}

// Пример пользовательской структуры для теста
type Task struct {
	ID   int
	Name string
}

func main() {
	fmt.Println("=== Generic Stack ===")

	// Тест со стеком целых чисел
	fmt.Println("\n--- Stack[int] ---")
	intStack := NewStack[int]()

	fmt.Printf("Пустой? %v (ожидается: true)\n", intStack.IsEmpty())

	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)

	fmt.Printf("Размер после Push(10,20,30): %d (ожидается: 3)\n", intStack.Len())
	fmt.Printf("Пустой? %v (ожидается: false)\n", intStack.IsEmpty())

	if val, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d (ожидается: 30)\n", val)
	}
	fmt.Printf("Размер после Peek: %d (ожидается: 3)\n", intStack.Len())

	if val, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d (ожидается: 30)\n", val)
	}
	if val, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d (ожидается: 20)\n", val)
	}
	fmt.Printf("Размер после двух Pop: %d (ожидается: 1)\n", intStack.Len())

	// Тест со стеком строк
	fmt.Println("\n--- Stack[string] ---")
	strStack := NewStack[string]()

	strStack.Push("first")
	strStack.Push("second")
	strStack.Push("third")

	fmt.Printf("Размер: %d (ожидается: 3)\n", strStack.Len())

	for !strStack.IsEmpty() {
		if val, ok := strStack.Pop(); ok {
			fmt.Printf("Pop: %s\n", val)
		}
	}
	// Ожидается: third, second, first (LIFO порядок)

	// Тест с пустым стеком
	fmt.Println("\n--- Пустой стек ---")
	emptyStack := NewStack[float64]()

	if _, ok := emptyStack.Pop(); !ok {
		fmt.Println("Pop на пустом стеке: корректно вернул false")
	}
	if _, ok := emptyStack.Peek(); !ok {
		fmt.Println("Peek на пустом стеке: корректно вернул false")
	}

	// Тест с пользовательской структурой
	fmt.Println("\n--- Stack[Task] ---")
	taskStack := NewStack[Task]()

	taskStack.Push(Task{ID: 1, Name: "Написать код"})
	taskStack.Push(Task{ID: 2, Name: "Написать тесты"})
	taskStack.Push(Task{ID: 3, Name: "Сделать ревью"})

	fmt.Printf("Количество задач: %d\n", taskStack.Len())

	if task, ok := taskStack.Pop(); ok {
		fmt.Printf("Следующая задача: %d - %s\n", task.ID, task.Name)
		// Ожидается: 3 - Сделать ревью
	}
	if task, ok := taskStack.Peek(); ok {
		fmt.Printf("Задача на очереди: %d - %s\n", task.ID, task.Name)
		// Ожидается: 2 - Написать тесты
	}
}
