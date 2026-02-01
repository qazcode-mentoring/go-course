# ДЗ 13.4: Buffered Tasks

## Цель
Научиться использовать буферизованные каналы для организации очереди задач и понять разницу между буферизованными и небуферизованными каналами.

## Что нужно сделать

1. Создать структуру `Task` для представления задачи
2. Реализовать функцию `worker`, которая обрабатывает задачи из канала
3. Реализовать функцию `processTasksBuffered`, которая создаёт буферизованный канал и обрабатывает задачи
4. Реализовать функцию `processTasksUnbuffered` для сравнения с небуферизованным каналом

## Сигнатуры

```go
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
func worker(id int, tasks <-chan Task, results chan<- Result)

// processTasksBuffered обрабатывает задачи через буферизованный канал
// bufferSize - размер буфера канала задач
// workerCount - количество воркеров
func processTasksBuffered(tasks []Task, bufferSize, workerCount int) []Result

// processTasksUnbuffered обрабатывает задачи через небуферизованный канал
func processTasksUnbuffered(tasks []Task, workerCount int) []Result
```

## Теория: буферизованные каналы

```go
// Небуферизованный канал (по умолчанию)
ch := make(chan int)
// Отправка блокируется, пока кто-то не прочитает

// Буферизованный канал
ch := make(chan int, 10)
// Отправка не блокируется, пока буфер не заполнен
```

### Когда использовать буфер?

- **Небуферизованный**: когда нужна синхронизация "рука-в-руку"
- **Буферизованный**: когда производитель и потребитель работают с разной скоростью

## Пример использования

```go
func main() {
    tasks := []Task{
        {ID: 1, Name: "Задача 1", Duration: 100 * time.Millisecond},
        {ID: 2, Name: "Задача 2", Duration: 200 * time.Millisecond},
        {ID: 3, Name: "Задача 3", Duration: 150 * time.Millisecond},
        {ID: 4, Name: "Задача 4", Duration: 50 * time.Millisecond},
        {ID: 5, Name: "Задача 5", Duration: 100 * time.Millisecond},
    }

    // Обработка с буфером на 3 задачи и 2 воркерами
    results := processTasksBuffered(tasks, 3, 2)

    for _, r := range results {
        fmt.Printf("Задача %d: %s\n", r.TaskID, r.Output)
    }
}
```

## Ожидаемый вывод

```
=== Buffered Tasks ===

--- Буферизованный канал (буфер=5, воркеров=3) ---
[Worker 1] Начал задачу 1: "Отправить email"
[Worker 2] Начал задачу 2: "Обработать изображение"
[Worker 3] Начал задачу 3: "Сгенерировать отчёт"
[Worker 1] Завершил задачу 1
[Worker 1] Начал задачу 4: "Синхронизировать данные"
[Worker 3] Завершил задачу 3
[Worker 3] Начал задачу 5: "Отправить уведомление"
[Worker 2] Завершил задачу 2
[Worker 1] Завершил задачу 4
[Worker 3] Завершил задачу 5

Результаты:
- Задача 1: Выполнено
- Задача 2: Выполнено
- Задача 3: Выполнено
- Задача 4: Выполнено
- Задача 5: Выполнено

Общее время: 350ms (vs 700ms последовательно)
```

## Подсказки

### Паттерн worker pool

```go
func processWithWorkers(tasks []Task, workerCount int) {
    taskCh := make(chan Task, len(tasks))  // Буфер на все задачи
    resultCh := make(chan Result, len(tasks))

    // Запускаем воркеров
    var wg sync.WaitGroup
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            worker(workerID, taskCh, resultCh)
        }(i + 1)
    }

    // Отправляем задачи
    for _, task := range tasks {
        taskCh <- task
    }
    close(taskCh)  // Сигнализируем воркерам, что задач больше нет

    // Ждём завершения воркеров
    wg.Wait()
    close(resultCh)

    // Собираем результаты
    for result := range resultCh {
        // ...
    }
}
```

### Симуляция работы

```go
func worker(id int, tasks <-chan Task, results chan<- Result) {
    for task := range tasks {
        fmt.Printf("[Worker %d] Начал задачу %d\n", id, task.ID)

        // Симулируем работу
        time.Sleep(task.Duration)

        results <- Result{
            TaskID:    task.ID,
            Output:    "Выполнено",
            Completed: time.Now(),
        }

        fmt.Printf("[Worker %d] Завершил задачу %d\n", id, task.ID)
    }
}
```

## Критерии выполнения

- [ ] Структуры Task и Result определены корректно
- [ ] worker обрабатывает задачи из канала и отправляет результаты
- [ ] processTasksBuffered использует буферизованный канал
- [ ] processTasksUnbuffered использует небуферизованный канал
- [ ] Множественные воркеры работают параллельно
- [ ] Все задачи обрабатываются, все результаты собираются
- [ ] Код корректно использует WaitGroup для синхронизации
- [ ] Каналы корректно закрываются
