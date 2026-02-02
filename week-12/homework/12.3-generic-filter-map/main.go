package main

import "fmt"

// Filter возвращает новый слайс с элементами, для которых predicate вернул true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: реализуй функцию
	// 1. Создай пустой слайс результатов
	// 2. Пройди по всем элементам
	// 3. Если predicate(element) == true, добавь в результат
	// 4. Верни результат
	var res []T

	for _, v := range slice {
		if predicate(v) {
			res = append(res, v)
		}
	}

	return res
}

// Map преобразует каждый элемент слайса с помощью функции transform
func Map[T, R any](slice []T, transform func(T) R) []R {
	// TODO: реализуй функцию
	// 1. Создай слайс результатов с той же длиной
	// 2. Пройди по всем элементам
	// 3. Преобразуй каждый элемент: result[i] = transform(slice[i])
	// 4. Верни результат

	res := make([]R, len(slice))

	for i := range slice {
		res[i] = transform(slice[i])
	}

	return res
}

// Reduce сворачивает слайс в одно значение
func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R {
	// TODO: реализуй функцию
	// 1. Инициализируй аккумулятор значением initial
	// 2. Пройди по всем элементам
	// 3. Обнови аккумулятор: acc = reducer(acc, element)
	// 4. Верни аккумулятор

	acc := initial

	for _, v := range slice {
		acc = reducer(acc, v)
	}

	return acc
}

// Find ищет первый элемент, для которого predicate вернул true
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	// TODO: реализуй функцию
	// 1. Пройди по всем элементам
	// 2. Если predicate(element) == true, верни (element, true)
	// 3. Если ничего не найдено, верни (zero, false)
	var zero T

	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}

	return zero, false
}

// Any возвращает true, если хотя бы один элемент удовлетворяет predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	// TODO: реализуй функцию
	// Верни true при первом совпадении

	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}

	return false
}

// All возвращает true, если все элементы удовлетворяют predicate
// Для пустого слайса возвращает true
func All[T any](slice []T, predicate func(T) bool) bool {
	// TODO: реализуй функцию
	// Верни false при первом несовпадении

	if len(slice) == 0 {
		return true
	}

	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}

	return true
}

// User - пример пользовательской структуры
type User struct {
	Name   string
	Age    int
	Active bool
}

func main() {
	fmt.Println("=== Generic Filter и Map ===")

	// Тест Filter
	fmt.Println("\n--- Filter ---")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("Исходный слайс: %v\n", nums)

	evens := Filter(nums, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Четные числа: %v (ожидается: [2 4 6 8 10])\n", evens)

	greaterThan5 := Filter(nums, func(n int) bool {
		return n > 5
	})
	fmt.Printf("Больше 5: %v (ожидается: [6 7 8 9 10])\n", greaterThan5)

	// Тест Map
	fmt.Println("\n--- Map ---")
	squares := Map(nums, func(n int) int {
		return n * n
	})
	fmt.Printf("Квадраты: %v\n", squares)
	fmt.Println("(ожидается: [1 4 9 16 25 36 49 64 81 100])")

	// Map с изменением типа
	strs := Map(nums, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Printf("Строки: %v\n", strs)
	fmt.Println("(ожидается: [#1 #2 #3 #4 #5 #6 #7 #8 #9 #10])")

	// Тест Reduce
	fmt.Println("\n--- Reduce ---")
	sum := Reduce(nums, 0, func(acc, n int) int {
		return acc + n
	})
	fmt.Printf("Сумма: %d (ожидается: 55)\n", sum)

	product := Reduce([]int{1, 2, 3, 4, 5}, 1, func(acc, n int) int {
		return acc * n
	})
	fmt.Printf("Произведение 1-5: %d (ожидается: 120)\n", product)

	// Конкатенация строк через Reduce
	words := []string{"Go", "is", "awesome"}
	sentence := Reduce(words, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + " " + s
	})
	fmt.Printf("Предложение: %q (ожидается: \"Go is awesome\")\n", sentence)

	// Тест Find
	fmt.Println("\n--- Find ---")
	if val, ok := Find(nums, func(n int) bool { return n > 5 }); ok {
		fmt.Printf("Первое число > 5: %d (ожидается: 6)\n", val)
	}

	if _, ok := Find(nums, func(n int) bool { return n > 100 }); !ok {
		fmt.Println("Число > 100 не найдено (ожидается)")
	}

	// Тест Any и All
	fmt.Println("\n--- Any и All ---")
	fmt.Printf("Any > 5: %v (ожидается: true)\n", Any(nums, func(n int) bool { return n > 5 }))
	fmt.Printf("Any > 100: %v (ожидается: false)\n", Any(nums, func(n int) bool { return n > 100 }))
	fmt.Printf("All > 0: %v (ожидается: true)\n", All(nums, func(n int) bool { return n > 0 }))
	fmt.Printf("All > 5: %v (ожидается: false)\n", All(nums, func(n int) bool { return n > 5 }))

	// Пустой слайс
	fmt.Println("\n--- Пустой слайс ---")
	empty := []int{}
	fmt.Printf("Any на пустом: %v (ожидается: false)\n", Any(empty, func(n int) bool { return true }))
	fmt.Printf("All на пустом: %v (ожидается: true)\n", All(empty, func(n int) bool { return false }))

	// Комплексный пример с пользователями
	fmt.Println("\n--- Работа с User ---")
	users := []User{
		{Name: "Алексей", Age: 25, Active: true},
		{Name: "Мария", Age: 17, Active: true},
		{Name: "Иван", Age: 30, Active: false},
		{Name: "Елена", Age: 22, Active: true},
	}

	// Найти активных взрослых
	activeAdults := Filter(users, func(u User) bool {
		return u.Age >= 18 && u.Active
	})
	fmt.Printf("Активные взрослые: %d человек\n", len(activeAdults))
	// Ожидается: 2 (Алексей и Елена)

	// Получить имена
	names := Map(activeAdults, func(u User) string {
		return u.Name
	})
	fmt.Printf("Имена: %v\n", names)
	// Ожидается: [Алексей Елена]

	// Средний возраст активных взрослых
	if len(activeAdults) > 0 {
		totalAge := Reduce(activeAdults, 0, func(sum int, u User) int {
			return sum + u.Age
		})
		avgAge := float64(totalAge) / float64(len(activeAdults))
		fmt.Printf("Средний возраст: %.1f\n", avgAge)
		// Ожидается: 23.5
	}
}
