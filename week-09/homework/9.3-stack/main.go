package main

import "fmt"

// Stack представляет стек целых чисел (LIFO)
type Stack struct {
	items []int
}

// NewStack создаёт и возвращает новый пустой стек
func NewStack() *Stack {
	// TODO: реализуй функцию
	return &Stack{make([]int, 0)}
}

// Push добавляет элемент на вершину стека
func (s *Stack) Push(value int) {
	// TODO: реализуй метод
	// Используй append
	s.items = append(s.items, value)
}

// Pop удаляет и возвращает верхний элемент
// Возвращает (value, true) или (0, false) если стек пуст
func (s *Stack) Pop() (int, bool) {
	// TODO: реализуй метод
	// Проверь, не пуст ли стек
	// Верни последний элемент и укороти слайс
	if len(s.items) == 0 {
		return 0, false
	}

	lastIndex := len(s.items) - 1
	lastValue := s.items[lastIndex]
	s.items = s.items[:lastIndex]

	return lastValue, true

}

// Peek возвращает верхний элемент без удаления
// Возвращает (value, true) или (0, false) если стек пуст
func (s *Stack) Peek() (int, bool) {
	// TODO: реализуй метод
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty возвращает true если стек пуст
func (s *Stack) IsEmpty() bool {
	// TODO: реализуй метод
	if len(s.items) != 0 {
		return false
	}
	return true
}

// Size возвращает количество элементов в стеке
func (s *Stack) Size() int {
	// TODO: реализуй метод
	return len(s.items)
}

func main() {
	stack := NewStack()

	fmt.Println("=== Stack (LIFO) ===")
	fmt.Println("Пустой стек?", stack.IsEmpty())

	// Добавляем элементы
	fmt.Println("\nPush: 1, 2, 3")
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	fmt.Println("Размер:", stack.Size())
	fmt.Println("Пустой стек?", stack.IsEmpty())

	// Peek
	if val, ok := stack.Peek(); ok {
		fmt.Println("Peek:", val)
	}

	// Pop
	fmt.Println("\nPop элементы:")
	for !stack.IsEmpty() {
		if val, ok := stack.Pop(); ok {
			fmt.Println("  Pop:", val)
		}
	}

	fmt.Println("\nПосле всех Pop:")
	fmt.Println("Размер:", stack.Size())
	fmt.Println("Пустой стек?", stack.IsEmpty())

	// Pop из пустого стека
	if _, ok := stack.Pop(); !ok {
		fmt.Println("Pop из пустого стека вернул ok=false")
	}
}
