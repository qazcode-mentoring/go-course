package main

import (
	"fmt"
	"sync"
)

// sayHello выводит приветствие с номером и именем
// Формат: "Привет #N от NAME!"
func sayHello(id int, name string, wg *sync.WaitGroup) {
	// TODO: реализуй функцию
	// 1. Используй defer wg.Done() для уменьшения счётчика WaitGroup
	// 2. Выведи приветствие в формате "Привет #N от NAME!"
}

// runConcurrent запускает горутины для каждого имени из списка
// и ждёт их завершения
func runConcurrent(names []string) {
	// TODO: реализуй функцию
	// 1. Создай sync.WaitGroup
	// 2. Для каждого имени:
	//    - Вызови wg.Add(1)
	//    - Запусти горутину с sayHello
	// 3. Дождись завершения всех горутин с wg.Wait()
}

func main() {
	fmt.Println("=== Concurrent Hello ===")
	fmt.Println()

	names := []string{"Алиса", "Боб", "Чарли", "Дэйв", "Ева"}

	// Запустим несколько раз, чтобы увидеть разный порядок
	for i := 1; i <= 3; i++ {
		fmt.Printf("--- Запуск %d ---\n", i)
		runConcurrent(names)
		fmt.Println()
	}

	// TODO: напиши здесь в комментарии, почему порядок вывода
	// может отличаться при каждом запуске
	// Ответ: ...
}
