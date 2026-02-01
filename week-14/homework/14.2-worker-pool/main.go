package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job представляет задачу для обработки
type Job struct {
	ID   int
	Data string
}

// Result представляет результат обработки задачи
type Result struct {
	JobID    int
	Output   string
	WorkerID int
	Duration time.Duration
}

// processJob обрабатывает одну задачу (симуляция работы).
// Эмулирует работу случайной продолжительности от 100ms до 500ms.
func processJob(job Job) string {
	// Эмуляция работы
	delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(delay)

	return fmt.Sprintf("processed: %s", job.Data)
}

// worker обрабатывает задачи из канала jobs и отправляет результаты в канал results.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	// TODO: реализуй функцию
	// 1. Используй defer wg.Done() для сигнализации о завершении
	// 2. Читай задачи из канала jobs в цикле (for job := range jobs)
	// 3. Для каждой задачи:
	//    a) Выведи сообщение о начале обработки
	//    b) Засеки время начала
	//    c) Вызови processJob для обработки
	//    d) Засеки время окончания
	//    e) Выведи сообщение о завершении
	//    f) Отправь Result в канал results

	_ = id      // удали после реализации
	_ = jobs    // удали после реализации
	_ = results // удали после реализации
	_ = wg      // удали после реализации
}

// RunWorkerPool запускает пул воркеров и обрабатывает задачи.
func RunWorkerPool(numWorkers int, jobs []Job) []Result {
	// TODO: реализуй функцию
	// 1. Создай канал для задач (буферизованный, размер = len(jobs))
	// 2. Создай канал для результатов (буферизованный, размер = len(jobs))
	// 3. Создай WaitGroup для воркеров
	//
	// 4. Запусти numWorkers воркеров:
	//    - wg.Add(1)
	//    - go worker(id, jobsChan, resultsChan, &wg)
	//
	// 5. Отправь все задачи в канал jobs
	// 6. Закрой канал jobs (воркеры завершатся, когда обработают все задачи)
	//
	// 7. Дождись завершения всех воркеров в отдельной горутине
	//    и закрой канал результатов:
	//    go func() {
	//        wg.Wait()
	//        close(resultsChan)
	//    }()
	//
	// 8. Собери все результаты из канала в слайс
	// 9. Верни слайс результатов

	_ = numWorkers // удали после реализации
	_ = jobs       // удали после реализации

	return nil
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== ДЗ 14.2: Worker Pool ===")
	fmt.Println()

	// Создаём задачи
	numJobs := 10
	jobs := make([]Job, numJobs)
	for i := 0; i < numJobs; i++ {
		jobs[i] = Job{
			ID:   i + 1,
			Data: fmt.Sprintf("task-%d", i+1),
		}
	}

	// Тест 1: Один воркер (последовательная обработка)
	fmt.Println("--- Тест 1: 1 воркер, 5 задач ---")
	start := time.Now()
	results := RunWorkerPool(1, jobs[:5])
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %v\n", elapsed)
	fmt.Printf("Результатов: %d\n", len(results))

	// Тест 2: Три воркера (параллельная обработка)
	fmt.Println("\n--- Тест 2: 3 воркера, 10 задач ---")
	start = time.Now()
	results = RunWorkerPool(3, jobs)
	elapsed = time.Since(start)
	fmt.Printf("\nВремя выполнения: %v\n", elapsed)
	fmt.Printf("Результатов: %d\n", len(results))

	// Выводим детали результатов
	fmt.Println("\nДетали результатов:")
	for _, r := range results {
		fmt.Printf("  Job %2d: %s (worker %d, %v)\n",
			r.JobID, r.Output, r.WorkerID, r.Duration)
	}

	// Тест 3: Много воркеров
	fmt.Println("\n--- Тест 3: 5 воркеров, 10 задач ---")
	start = time.Now()
	results = RunWorkerPool(5, jobs)
	elapsed = time.Since(start)
	fmt.Printf("Время выполнения: %v\n", elapsed)
	fmt.Printf("Результатов: %d\n", len(results))

	// Проверка: с большим числом воркеров время должно быть меньше
	fmt.Println("\n--- Сравнение производительности ---")
	fmt.Println("Больше воркеров = меньше общее время (до определённого предела)")
}
