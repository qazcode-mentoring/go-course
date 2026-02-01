package main

import (
	"fmt"
	"sync"
	"time"
)

// Task представляет задачу для обработки
type Task struct {
	ID       int
	Name     string
	Duration time.Duration // Время выполнения задачи
}

// Result представляет результат выполнения задачи
type Result struct {
	TaskID    int
	Output    string
	Completed time.Time
}

// worker обрабатывает задачи из канала tasks и отправляет результаты в results
func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	// TODO: реализуй функцию
	// 1. Используй defer wg.Done()
	// 2. В цикле for-range читай задачи из канала tasks
	// 3. Для каждой задачи:
	//    - Выведи: "[Worker ID] Начал задачу TaskID: "TaskName""
	//    - Симулируй работу: time.Sleep(task.Duration)
	//    - Создай Result и отправь в канал results
	//    - Выведи: "[Worker ID] Завершил задачу TaskID"
}

// processTasksBuffered обрабатывает задачи через буферизованный канал
// bufferSize - размер буфера канала задач
// workerCount - количество воркеров
func processTasksBuffered(tasks []Task, bufferSize, workerCount int) []Result {
	// TODO: реализуй функцию
	// 1. Создай буферизованный канал задач: make(chan Task, bufferSize)
	// 2. Создай буферизованный канал результатов: make(chan Result, len(tasks))
	// 3. Создай WaitGroup для воркеров
	// 4. Запусти workerCount воркеров в горутинах
	// 5. Отправь все задачи в канал задач
	// 6. Закрой канал задач
	// 7. Дождись завершения всех воркеров (wg.Wait())
	// 8. Закрой канал результатов
	// 9. Собери результаты из канала в слайс и верни
	return nil
}

// processTasksUnbuffered обрабатывает задачи через небуферизованный канал
// workerCount - количество воркеров
func processTasksUnbuffered(tasks []Task, workerCount int) []Result {
	// TODO: реализуй функцию
	// Аналогично processTasksBuffered, но с небуферизованным каналом
	// make(chan Task) вместо make(chan Task, bufferSize)
	//
	// ВАЖНО: при небуферизованном канале отправка будет блокироваться,
	// пока воркер не прочитает задачу. Поэтому отправку задач нужно
	// делать в отдельной горутине!
	return nil
}

func main() {
	fmt.Println("=== Buffered Tasks ===")

	// Создаём список задач
	tasks := []Task{
		{ID: 1, Name: "Отправить email", Duration: 100 * time.Millisecond},
		{ID: 2, Name: "Обработать изображение", Duration: 200 * time.Millisecond},
		{ID: 3, Name: "Сгенерировать отчёт", Duration: 150 * time.Millisecond},
		{ID: 4, Name: "Синхронизировать данные", Duration: 80 * time.Millisecond},
		{ID: 5, Name: "Отправить уведомление", Duration: 50 * time.Millisecond},
		{ID: 6, Name: "Обновить кэш", Duration: 120 * time.Millisecond},
	}

	// Вычисляем общее время последовательного выполнения
	var totalSequential time.Duration
	for _, t := range tasks {
		totalSequential += t.Duration
	}

	// Тест с буферизованным каналом
	fmt.Println("\n--- Буферизованный канал (буфер=5, воркеров=3) ---")
	start := time.Now()
	results := processTasksBuffered(tasks, 5, 3)
	bufferedDuration := time.Since(start)

	fmt.Println("\nРезультаты:")
	for _, r := range results {
		fmt.Printf("- Задача %d: %s\n", r.TaskID, r.Output)
	}
	fmt.Printf("\nВремя выполнения: %v (vs %v последовательно)\n",
		bufferedDuration.Round(time.Millisecond),
		totalSequential)

	// Тест с небуферизованным каналом
	fmt.Println("\n--- Небуферизованный канал (воркеров=3) ---")
	start = time.Now()
	results = processTasksUnbuffered(tasks, 3)
	unbufferedDuration := time.Since(start)

	fmt.Println("\nРезультаты:")
	for _, r := range results {
		fmt.Printf("- Задача %d: %s\n", r.TaskID, r.Output)
	}
	fmt.Printf("\nВремя выполнения: %v\n", unbufferedDuration.Round(time.Millisecond))

	// Сравнение
	fmt.Println("\n--- Сравнение ---")
	fmt.Printf("Последовательно:      %v\n", totalSequential)
	fmt.Printf("Буферизованный:       %v\n", bufferedDuration.Round(time.Millisecond))
	fmt.Printf("Небуферизованный:     %v\n", unbufferedDuration.Round(time.Millisecond))
}
