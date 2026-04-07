package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Task представляет задачу для обработки
type Task struct {
	ID       int
	Data     string
	Priority int
}

// Result представляет результат обработки задачи
type Result struct {
	TaskID   int
	WorkerID int
	Output   string
	Duration time.Duration
	Error    error
}

// Worker обрабатывает задачи из канала
type Worker struct {
	ID int
	// TODO: добавь нужные поля
}

// NewWorker создаёт нового воркера
func NewWorker(id int) *Worker {
	return &Worker{ID: id}
}

// processTask эмулирует обработку задачи
func (w *Worker) processTask(ctx context.Context, task Task) Result {
	start := time.Now()

	// Случайное время обработки от 100ms до 500ms
	processingTime := time.Duration(100+rand.IntN(400)) * time.Millisecond

	select {
	case <-time.After(processingTime):
		return Result{
			TaskID:   task.ID,
			WorkerID: w.ID,
			Output:   fmt.Sprintf("Processed: %s", task.Data),
			Duration: time.Since(start),
		}
	case <-ctx.Done():
		return Result{
			TaskID:   task.ID,
			WorkerID: w.ID,
			Duration: time.Since(start),
			Error:    ctx.Err(),
		}
	}
}

// Start запускает воркера. Воркер читает задачи из jobs и отправляет результаты в results.
// При отмене контекста воркер завершает текущую задачу и выходит.
func (w *Worker) Start(ctx context.Context, jobs <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: завершение по контексту\n", w.ID)
			return
		case task, ok := <-jobs:
			if !ok {
				fmt.Printf("Worker %d: канал задач закрыт\n", w.ID)
				return
			}

			fmt.Printf("Worker %d: обработка задачи %d\n", w.ID, task.ID)
			result := w.processTask(ctx, task)

			// Отправляем результат, учитывая возможную отмену контекста
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
	workers    []*Worker
	jobs       chan Task
	results    chan Result
	workerWg   sync.WaitGroup
	started    bool
	stopped    bool
	mu         sync.Mutex
	cancelFunc context.CancelFunc
}

// NewWorkerPool создаёт пул из n воркеров
func NewWorkerPool(n int) *WorkerPool {
	workers := make([]*Worker, n)
	for i := range n {
		workers[i] = NewWorker(i + 1)
	}

	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Task, 100), // буферизированный канал
		results: make(chan Result, 100),
	}
}

// Start запускает пул воркеров
func (p *WorkerPool) Start(ctx context.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return
	}

	// Создаем дочерний контекст для управления воркерами
	var poolCtx context.Context
	poolCtx, p.cancelFunc = context.WithCancel(ctx)

	for _, worker := range p.workers {
		p.workerWg.Add(1)
		go worker.Start(poolCtx, p.jobs, p.results, &p.workerWg)
	}

	p.started = true
}

// Submit добавляет задачу в очередь. Возвращает false если пул остановлен.
func (p *WorkerPool) Submit(task Task) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return false
	}

	p.jobs <- task
	return true
}

// Results возвращает канал с результатами
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Shutdown корректно останавливает пул
// Ждёт завершения всех задач или до истечения таймаута
func (p *WorkerPool) Shutdown(timeout time.Duration) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return nil
	}
	p.stopped = true
	p.mu.Unlock()

	close(p.jobs)  // Больше задачи не принимаем
	p.cancelFunc() // Сигнализируем воркерам об отмене

	done := make(chan struct{})
	go func() {
		p.workerWg.Wait()
		close(p.results) // Закрываем результаты только когда все воркеры вышли
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout")
	}
}

// Service представляет сервис с фоновыми задачами
type Service struct {
	pool      *WorkerPool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	startTime time.Time
}

// NewService создаёт новый сервис
func NewService() *Service {
	return &Service{
		pool: NewWorkerPool(4),
	}
}

// Start запускает сервис и все фоновые задачи
func (s *Service) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.startTime = time.Now()

	// Запускаем пул
	s.pool.Start(s.ctx)

	// Фоновые задачи
	s.wg.Add(1)
	go s.runHealthCheck()

	s.wg.Add(1)
	go s.processResults()

	return nil
}

// runHealthCheck выполняет периодическую проверку (каждые 2 секунды)
func (s *Service) runHealthCheck() {
	defer s.wg.Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("Health check: остановлен")
			return
		case <-ticker.C:
			fmt.Printf("Health check: OK (uptime: %v)\n",
				time.Since(s.startTime).Round(time.Second))
		}
	}
}

// processResults обрабатывает результаты из пула
func (s *Service) processResults() {
	defer s.wg.Done()

	for result := range s.pool.Results() {
		if result.Error != nil {
			fmt.Printf("Результат задачи %d: ошибка - %v\n",
				result.TaskID, result.Error)
		} else {
			fmt.Printf("Результат задачи %d (worker %d): %s (%v)\n",
				result.TaskID, result.WorkerID, result.Output,
				result.Duration.Round(time.Millisecond))
		}
	}

	fmt.Println("Обработчик результатов: остановлен")
}

// Shutdown корректно останавливает сервис
func (s *Service) Shutdown(ctx context.Context) error {
	fmt.Println("Сервис: начало остановки...")

	s.cancel() // 1. Отменяем контекст сервиса

	// 2. Останавливаем пул (даем ему 5 секунд на завершение задач)
	err := s.pool.Shutdown(5 * time.Second)

	// 3. Ждем завершения фоновых задач (HealthCheck и processResults)
	waitDone := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		fmt.Println("Сервис остановлен")
		return err
	case <-ctx.Done():
		return fmt.Errorf("service shutdown interrupted: %w", ctx.Err())
	}
}

// SubmitTask добавляет задачу в сервис
func (s *Service) SubmitTask(task Task) bool {
	return s.pool.Submit(task)
}

func main() {
	fmt.Println("=== ДЗ 15.4: Graceful Shutdown ===")
	fmt.Println()

	// Создаём контекст с обработкой сигналов
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Создаём и запускаем сервис
	svc := NewService()

	fmt.Println("Запуск сервиса...")
	if err := svc.Start(ctx); err != nil {
		fmt.Printf("Ошибка запуска: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Сервис запущен. Нажми Ctrl+C для остановки.")
	fmt.Println()

	// Добавляем задачи
	go func() {
		for i := range 10 {
			task := Task{
				ID:       i + 1,
				Data:     fmt.Sprintf("data-%d", i+1),
				Priority: rand.IntN(3) + 1,
			}

			if svc.SubmitTask(task) {
				fmt.Printf("Добавлена задача %d\n", task.ID)
			} else {
				fmt.Printf("Не удалось добавить задачу %d (сервис остановлен)\n", task.ID)
				return
			}

			// Небольшая задержка между задачами
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Ждём сигнала завершения
	<-ctx.Done()
	fmt.Println()
	fmt.Println("=== Получен сигнал завершения ===")

	// Корректно останавливаем сервис
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Остановка сервиса...")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Ошибка при остановке: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== Сервис успешно остановлен ===")
}
