# ДЗ 14.4: Concurrent Cache

## Цель
Научиться использовать `sync.RWMutex` для создания эффективного конкурентного кэша с множественным чтением.

## Что нужно сделать

Реализовать потокобезопасный кэш с поддержкой TTL (time-to-live) для записей. Кэш должен эффективно обрабатывать много одновременных операций чтения.

## Сигнатуры

```go
// CacheItem представляет элемент кэша
type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

// Cache представляет потокобезопасный кэш с TTL
type Cache struct {
    // TODO: добавь поля
}

// NewCache создаёт новый кэш
func NewCache() *Cache

// Set добавляет или обновляет значение в кэше с указанным TTL.
// Если ttl <= 0, запись не истекает.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration)

// Get возвращает значение из кэша.
// Возвращает nil и false, если ключ не найден или истёк.
func (c *Cache) Get(key string) (interface{}, bool)

// Delete удаляет значение из кэша.
// Возвращает true, если значение было удалено.
func (c *Cache) Delete(key string) bool

// Count возвращает количество элементов в кэше (включая истёкшие).
func (c *Cache) Count() int

// Cleanup удаляет все истёкшие записи.
// Возвращает количество удалённых записей.
func (c *Cache) Cleanup() int
```

## Теория: RWMutex

`sync.RWMutex` (Read-Write Mutex) - это мьютекс, который различает операции чтения и записи:

- **RLock/RUnlock** - блокировка на чтение. Несколько горутин могут держать RLock одновременно.
- **Lock/Unlock** - блокировка на запись. Только одна горутина может держать Lock, при этом никто не может держать RLock.

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

// Чтение - можно одновременно из многих горутин
func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

// Запись - эксклюзивный доступ
func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}
```

### Когда RWMutex лучше Mutex?

RWMutex даёт преимущество, когда:
- Операций чтения значительно больше, чем записи (>10:1)
- Операции чтения занимают заметное время
- Много горутин одновременно читают данные

Для простых случаев обычный `Mutex` может быть даже быстрее из-за меньших накладных расходов.

## Пример использования

```go
func main() {
    cache := NewCache()

    // Записываем данные
    cache.Set("user:1", "Alice", 5*time.Second)
    cache.Set("user:2", "Bob", 10*time.Second)
    cache.Set("config", map[string]int{"timeout": 30}, 0) // Без TTL

    // Читаем данные
    if val, ok := cache.Get("user:1"); ok {
        fmt.Println("user:1 =", val)
    }

    // Параллельное чтение
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cache.Get("config")
        }()
    }
    wg.Wait()

    // Ждём истечения TTL
    time.Sleep(6 * time.Second)

    if _, ok := cache.Get("user:1"); !ok {
        fmt.Println("user:1 истёк")
    }

    // Очистка истёкших записей
    deleted := cache.Cleanup()
    fmt.Printf("Удалено истёкших записей: %d\n", deleted)
}
```

## Ожидаемый вывод

```
=== Тест 1: Базовые операции ===
Set user:1 = Alice (TTL: 2s)
Set user:2 = Bob (TTL: 5s)
Set config = {...} (без TTL)
Get user:1: Alice (ok: true)
Get user:2: Bob (ok: true)
Get config: map[timeout:30] (ok: true)
Get unknown: <nil> (ok: false)
Count: 3

=== Тест 2: TTL ===
Ждём 3 секунды...
Get user:1: <nil> (ok: false) - истёк!
Get user:2: Bob (ok: true) - ещё живой
Get config: map[timeout:30] (ok: true) - без TTL

=== Тест 3: Cleanup ===
Удалено истёкших записей: 1
Count после cleanup: 2

=== Тест 4: Конкурентный доступ ===
Запуск 50 писателей и 200 читателей...
Все операции завершены!
Тест ПРОЙДЕН!
```

## Подсказки

- Для Get используй `RLock/RUnlock`
- Для Set, Delete, Cleanup используй `Lock/Unlock`
- Проверяй TTL при каждом Get, а не только в Cleanup
- При проверке TTL используй `time.Now().After(item.ExpiresAt)`
- TTL = 0 означает "никогда не истекает" (ExpiresAt будет нулевым значением)
- Используй `time.Time.IsZero()` для проверки нулевого времени

## Бонусное задание

Добавь автоматическую очистку истёкших записей в фоновом режиме:

```go
// NewCacheWithCleanup создаёт кэш с автоматической очисткой каждые interval.
// Возвращает кэш и функцию для остановки очистки.
func NewCacheWithCleanup(interval time.Duration) (*Cache, func())
```

## Критерии выполнения

- [ ] Структура `Cache` использует `sync.RWMutex`
- [ ] Метод `Get` использует `RLock` для множественного чтения
- [ ] Методы `Set`, `Delete`, `Cleanup` используют `Lock`
- [ ] TTL корректно работает (истёкшие записи не возвращаются)
- [ ] Метод `Cleanup` удаляет только истёкшие записи
- [ ] Тест с конкурентным доступом проходит успешно
- [ ] Программа проходит проверку `go run -race` без предупреждений
