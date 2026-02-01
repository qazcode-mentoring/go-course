# Неделя 14: Паттерны конкурентности (Select, WaitGroup, Mutex)

## Теория

**Ссылка:** https://www.gocat.dev/tour/

Повтори уроки 90-96 и изучи дополнительные материалы по синхронизации.

## Паттерны синхронизации в Go

На прошлой неделе мы познакомились с горутинами и каналами. Эта неделя посвящена продвинутым паттернам конкурентности и примитивам синхронизации из пакета `sync`.

### Select: мультиплексирование каналов

`select` позволяет горутине ожидать несколько операций с каналами одновременно. Это мощный инструмент для:
- Таймаутов операций
- Отмены через context
- Неблокирующих операций (с `default`)
- Объединения нескольких источников данных (fan-in)

```go
// Таймаут операции
select {
case result := <-ch:
    fmt.Println("Получен результат:", result)
case <-time.After(5 * time.Second):
    fmt.Println("Таймаут! Операция заняла слишком много времени")
}

// Отмена через context
select {
case result := <-ch:
    processResult(result)
case <-ctx.Done():
    fmt.Println("Операция отменена:", ctx.Err())
    return
}

// Неблокирующая проверка канала
select {
case msg := <-ch:
    fmt.Println("Есть сообщение:", msg)
default:
    fmt.Println("Сообщений нет, продолжаем работу")
}
```

### sync.WaitGroup: ожидание группы горутин

`WaitGroup` используется для ожидания завершения набора горутин. Это основной паттерн для параллельной обработки данных.

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        processTask(id)
    }(i)
}

wg.Wait() // Ждём завершения всех горутин
fmt.Println("Все задачи выполнены!")
```

**Важно:**
- `Add()` вызывается ДО запуска горутины
- `Done()` вызывается в конце работы горутины (обычно через `defer`)
- `Wait()` блокирует до обнуления счётчика

### sync.Mutex: взаимное исключение

`Mutex` (mutual exclusion) защищает общие данные от одновременного доступа из нескольких горутин. Это предотвращает состояние гонки (race condition).

```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

**Правила работы с Mutex:**
- Всегда используй `defer mu.Unlock()` сразу после `Lock()`
- Не держи блокировку дольше, чем необходимо
- Не забывай разблокировать (иначе будет deadlock)
- Mutex нельзя копировать после первого использования

### sync.RWMutex: множественное чтение, единичная запись

`RWMutex` оптимизирован для случаев, когда чтений значительно больше, чем записей. Позволяет нескольким горутинам читать данные одновременно.

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

// Чтение - можно одновременно из многих горутин
func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

// Запись - эксклюзивный доступ
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}
```

**Когда использовать RWMutex:**
- Много операций чтения, мало записи
- Чтение занимает заметное время
- Для простых случаев обычный `Mutex` может быть быстрее

### Выбор инструмента синхронизации

| Задача | Инструмент |
|--------|-----------|
| Передача данных между горутинами | Каналы |
| Ожидание завершения горутин | WaitGroup |
| Защита общих данных | Mutex / RWMutex |
| Таймаут / отмена операции | select + time.After / context |
| Однократная инициализация | sync.Once |

### Race Detector

Go имеет встроенный детектор гонок. Используй его при разработке:

```bash
go run -race main.go
go test -race ./...
```

Детектор замедляет программу, но находит проблемы с конкурентным доступом.

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 14.1 | [select-timeout](./homework/14.1-select-timeout/) | Таймаут и отмена операций через select |
| 14.2 | [worker-pool](./homework/14.2-worker-pool/) | Пул воркеров с WaitGroup |
| 14.3 | [safe-counter](./homework/14.3-safe-counter/) | Потокобезопасный счётчик с Mutex |
| 14.4 | [concurrent-cache](./homework/14.4-concurrent-cache/) | Конкурентный кэш с RWMutex |

## Вопросы для самопроверки

Создай файл `answers-14.txt` и напиши ответы на вопросы:

1. В чём разница между `Mutex` и `RWMutex`? В каких случаях `RWMutex` даёт преимущество?
2. Что произойдёт, если вызвать `Lock()` дважды подряд в одной горутине? Как избежать этой проблемы?
3. Зачем в `select` используется `case <-ctx.Done()`? Как это связано с отменой операций?
4. Почему важно вызывать `wg.Add(1)` до запуска горутины, а не внутри неё?
5. Как с помощью `select` и `default` сделать неблокирующую отправку в канал?

## Дополнительные материалы

- [Go by Example: Mutexes](https://gobyexample.com/mutexes)
- [Go by Example: WaitGroups](https://gobyexample.com/waitgroups)
- [Go by Example: Select](https://gobyexample.com/select)
- [Go by Example: Timeouts](https://gobyexample.com/timeouts)
- [Go by Example: Non-Blocking Channel Operations](https://gobyexample.com/non-blocking-channel-operations)
- [Go Blog: Race Detector](https://go.dev/blog/race-detector)
- [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Wiki: MutexOrChannel](https://go.dev/wiki/MutexOrChannel)
