# ДЗ 12.4: Generic Cache с TTL

## Цель
Научиться создавать сложные обобщенные структуры данных с несколькими параметрами типа.

## Что нужно сделать

Реализовать универсальный кэш `Cache[K, V]` с поддержкой времени жизни записей (TTL - Time To Live):

1. `Set(key K, value V)` - добавить значение с дефолтным TTL
2. `SetWithTTL(key K, value V, ttl time.Duration)` - добавить значение с указанным TTL
3. `Get(key K) (V, bool)` - получить значение (проверяя срок годности)
4. `Delete(key K)` - удалить значение
5. `Len() int` - количество валидных записей
6. `Clear()` - очистить весь кэш
7. `Cleanup()` - удалить просроченные записи

## Сигнатуры

```go
// cacheEntry хранит значение и время истечения
type cacheEntry[V any] struct {
    value      V
    expiration time.Time
}

// Cache - универсальный кэш с TTL
// K - тип ключа (должен быть comparable для использования в map)
// V - тип значения (любой)
type Cache[K comparable, V any] struct {
    data       map[K]cacheEntry[V]
    defaultTTL time.Duration
}

// NewCache создает новый кэш с указанным TTL по умолчанию
func NewCache[K comparable, V any](defaultTTL time.Duration) *Cache[K, V]

// Set добавляет значение с дефолтным TTL
func (c *Cache[K, V]) Set(key K, value V)

// SetWithTTL добавляет значение с указанным TTL
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration)

// Get возвращает значение, если оно существует и не истекло
// Возвращает (value, true) если найдено и валидно
// Возвращает (zero, false) если не найдено или истекло
func (c *Cache[K, V]) Get(key K) (V, bool)

// Delete удаляет значение из кэша
func (c *Cache[K, V]) Delete(key K)

// Len возвращает количество валидных (не истекших) записей
func (c *Cache[K, V]) Len() int

// Clear очищает весь кэш
func (c *Cache[K, V]) Clear()

// Cleanup удаляет все истекшие записи
func (c *Cache[K, V]) Cleanup()
```

## Пример использования

```go
func main() {
    // Кэш строк с ключом-строкой, TTL по умолчанию 5 секунд
    cache := NewCache[string, string](5 * time.Second)

    // Добавляем данные
    cache.Set("user:1", "Alice")
    cache.Set("user:2", "Bob")

    // Получаем данные
    if val, ok := cache.Get("user:1"); ok {
        fmt.Println(val) // "Alice"
    }

    // Добавляем с кастомным TTL
    cache.SetWithTTL("session:abc", "token123", 1*time.Second)

    // Ждем истечения
    time.Sleep(2 * time.Second)

    // session:abc истек
    if _, ok := cache.Get("session:abc"); !ok {
        fmt.Println("Сессия истекла")
    }

    // user:1 еще валиден (TTL 5 секунд)
    if val, ok := cache.Get("user:1"); ok {
        fmt.Println(val) // "Alice"
    }

    // Кэш с числовым ключом
    intCache := NewCache[int, []byte](time.Minute)
    intCache.Set(42, []byte("binary data"))
}
```

## Подсказки

- Используй `time.Now().Add(ttl)` для вычисления времени истечения
- Для проверки истечения: `time.Now().After(entry.expiration)`
- В `Get` проверяй истечение и возвращай false если запись просрочена
- В `Len` считай только валидные записи
- `Cleanup` полезен для периодической очистки памяти

## Работа с временем в Go

```go
import "time"

// Текущее время
now := time.Now()

// Добавить duration
future := now.Add(5 * time.Second)

// Сравнение времени
if time.Now().After(future) {
    // future в прошлом
}

// Пауза
time.Sleep(time.Second)
```

## Бонусное задание (необязательно)

Добавь метод `GetOrSet`, который возвращает существующее значение или устанавливает новое:

```go
// GetOrSet возвращает существующее значение или вычисляет и сохраняет новое
func (c *Cache[K, V]) GetOrSet(key K, compute func() V) V
```

## Критерии выполнения

- [ ] Cache работает с разными комбинациями типов K и V
- [ ] Set корректно добавляет записи с дефолтным TTL
- [ ] SetWithTTL корректно добавляет записи с кастомным TTL
- [ ] Get не возвращает истекшие записи
- [ ] Delete удаляет записи
- [ ] Len возвращает только количество валидных записей
- [ ] Clear очищает весь кэш
- [ ] Cleanup удаляет просроченные записи
- [ ] Код компилируется без ошибок
