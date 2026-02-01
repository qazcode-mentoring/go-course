package main

import (
	"fmt"
	"sync"
)

// Counter представляет потокобезопасный счётчик
type Counter struct {
	// TODO: добавь поля
	// - mu sync.Mutex для защиты данных
	// - value int для хранения значения
}

// NewCounter создаёт новый счётчик с начальным значением 0
func NewCounter() *Counter {
	// TODO: реализуй функцию
	return &Counter{}
}

// Increment увеличивает счётчик на 1 и возвращает новое значение
func (c *Counter) Increment() int {
	// TODO: реализуй метод
	// 1. Захвати мьютекс
	// 2. Увеличь значение
	// 3. Сохрани новое значение в локальную переменную
	// 4. Освободи мьютекс (используй defer)
	// 5. Верни новое значение

	return 0
}

// Decrement уменьшает счётчик на 1 и возвращает новое значение
func (c *Counter) Decrement() int {
	// TODO: реализуй метод
	// Аналогично Increment, но уменьшаем

	return 0
}

// Add добавляет delta к счётчику и возвращает новое значение
func (c *Counter) Add(delta int) int {
	// TODO: реализуй метод

	_ = delta // удали после реализации
	return 0
}

// Value возвращает текущее значение счётчика
func (c *Counter) Value() int {
	// TODO: реализуй метод
	// Не забудь: даже чтение должно быть защищено мьютексом!

	return 0
}

// Reset сбрасывает счётчик в 0 и возвращает предыдущее значение
func (c *Counter) Reset() int {
	// TODO: реализуй метод
	// 1. Захвати мьютекс
	// 2. Сохрани текущее значение
	// 3. Установи значение в 0
	// 4. Верни старое значение

	return 0
}

func main() {
	fmt.Println("=== ДЗ 14.3: Safe Counter ===")
	fmt.Println()

	// Тест 1: Базовые операции
	fmt.Println("--- Тест 1: Базовые операции ---")
	counter := NewCounter()

	for i := 0; i < 5; i++ {
		counter.Increment()
	}
	fmt.Printf("После 5 инкрементов: %d\n", counter.Value())

	counter.Decrement()
	fmt.Printf("После декремента: %d\n", counter.Value())

	counter.Add(10)
	fmt.Printf("После Add(10): %d\n", counter.Value())

	old := counter.Reset()
	fmt.Printf("После Reset: %d (было: %d)\n", counter.Value(), old)

	// Тест 2: Конкурентный доступ
	fmt.Println("\n--- Тест 2: Конкурентный доступ ---")

	counter = NewCounter()
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 1000
	expected := numGoroutines * incrementsPerGoroutine

	fmt.Printf("Запущено %d горутин, каждая делает %d инкрементов\n",
		numGoroutines, incrementsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	result := counter.Value()
	fmt.Printf("Итого: %d (ожидалось: %d)\n", result, expected)

	if result == expected {
		fmt.Println("Тест ПРОЙДЕН!")
	} else {
		fmt.Println("Тест ПРОВАЛЕН! Обнаружена потеря данных из-за race condition")
	}

	// Тест 3: Смешанные операции
	fmt.Println("\n--- Тест 3: Смешанные операции ---")

	counter = NewCounter()

	numIncrements := 50000
	numDecrements := 50000

	fmt.Printf("Инкременты: %d, Декременты: %d\n", numIncrements, numDecrements)

	// Половина горутин инкрементирует
	for i := 0; i < numIncrements; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	// Половина горутин декрементирует
	for i := 0; i < numDecrements; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Decrement()
		}()
	}

	wg.Wait()

	result = counter.Value()
	expectedMixed := numIncrements - numDecrements

	fmt.Printf("Итого: %d (ожидалось: %d)\n", result, expectedMixed)

	if result == expectedMixed {
		fmt.Println("Тест ПРОЙДЕН!")
	} else {
		fmt.Println("Тест ПРОВАЛЕН!")
	}

	// Напоминание о проверке
	fmt.Println("\n--- Проверка race detector ---")
	fmt.Println("Запусти программу с флагом -race:")
	fmt.Println("  go run -race main.go")
	fmt.Println("Если есть race condition, Go выведет WARNING: DATA RACE")
}
