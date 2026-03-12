package main

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem представляет элемент кэша
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired проверяет, истёк ли срок жизни элемента
func (ci *CacheItem) IsExpired() bool {
	// Если ExpiresAt нулевой, элемент никогда не истекает
	if ci.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(ci.ExpiresAt)
}

// Cache представляет потокобезопасный кэш с TTL
type Cache struct {
	// TODO: добавь поля
	// - mu sync.RWMutex для защиты данных
	// - items map[string]CacheItem для хранения элементов
	mu    sync.RWMutex
	items map[string]CacheItem
}

// NewCache создаёт новый кэш
func NewCache() *Cache {
	// TODO: реализуй функцию
	// Не забудь инициализировать map!
	return &Cache{mu: sync.RWMutex{}, items: map[string]CacheItem{}}
}

// Set добавляет или обновляет значение в кэше с указанным TTL.
// Если ttl <= 0, запись не истекает.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	// TODO: реализуй метод
	// 1. Захвати мьютекс на запись (Lock)
	// 2. Вычисли ExpiresAt:
	//    - если ttl > 0: time.Now().Add(ttl)
	//    - иначе: time.Time{} (нулевое значение)
	// 3. Создай CacheItem и сохрани в map
	// 4. Освободи мьютекс

	c.mu.Lock()
	defer c.mu.Unlock()

	var expires time.Time
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	} else {
		expires = time.Time{}
	}

	item := CacheItem{
		Value:     value,
		ExpiresAt: expires,
	}

	c.items[key] = item
}

// Get возвращает значение из кэша.
// Возвращает nil и false, если ключ не найден или истёк.
func (c *Cache) Get(key string) (interface{}, bool) {
	// TODO: реализуй метод
	// 1. Захвати мьютекс на чтение (RLock)
	// 2. Найди элемент в map
	// 3. Если не найден - верни nil, false
	// 4. Если найден, но истёк - верни nil, false
	// 5. Верни значение и true
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.items[key]; ok {
		if v.IsExpired() {
			return nil, false
		}
		return v.Value, true
	}

	return nil, false
}

// Delete удаляет значение из кэша.
// Возвращает true, если значение было удалено.
func (c *Cache) Delete(key string) bool {
	// TODO: реализуй метод
	// 1. Захвати мьютекс на запись (Lock)
	// 2. Проверь, есть ли ключ
	// 3. Удали ключ из map
	// 4. Верни true, если ключ был
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		delete(c.items, key)
		return true
	}
	return false
}

// Count возвращает количество элементов в кэше (включая истёкшие).
func (c *Cache) Count() int {
	// TODO: реализуй метод
	// Используй RLock для чтения
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Cleanup удаляет все истёкшие записи.
// Возвращает количество удалённых записей.
func (c *Cache) Cleanup() int {
	// TODO: реализуй метод
	// 1. Захвати мьютекс на запись (Lock)
	// 2. Пройди по всем элементам
	// 3. Если элемент истёк - удали его
	// 4. Посчитай количество удалённых
	// 5. Верни счётчик
	c.mu.Lock()
	defer c.mu.Unlock()

	countDeleted := 0
	for key, item := range c.items {
		if item.IsExpired() {
			delete(c.items, key)
			countDeleted++
		}
	}
	return countDeleted
}

func main() {
	fmt.Println("=== ДЗ 14.4: Concurrent Cache ===")
	fmt.Println()

	// Тест 1: Базовые операции
	fmt.Println("--- Тест 1: Базовые операции ---")
	cache := NewCache()

	cache.Set("user:1", "Alice", 2*time.Second)
	fmt.Println("Set user:1 = Alice (TTL: 2s)")

	cache.Set("user:2", "Bob", 5*time.Second)
	fmt.Println("Set user:2 = Bob (TTL: 5s)")

	cache.Set("config", map[string]int{"timeout": 30}, 0)
	fmt.Println("Set config = {...} (без TTL)")

	val, ok := cache.Get("user:1")
	fmt.Printf("Get user:1: %v (ok: %v)\n", val, ok)

	val, ok = cache.Get("user:2")
	fmt.Printf("Get user:2: %v (ok: %v)\n", val, ok)

	val, ok = cache.Get("config")
	fmt.Printf("Get config: %v (ok: %v)\n", val, ok)

	val, ok = cache.Get("unknown")
	fmt.Printf("Get unknown: %v (ok: %v)\n", val, ok)

	fmt.Printf("Count: %d\n", cache.Count())

	// Тест 2: TTL
	fmt.Println("\n--- Тест 2: TTL ---")
	fmt.Println("Ждём 3 секунды...")
	time.Sleep(3 * time.Second)

	val, ok = cache.Get("user:1")
	if !ok {
		fmt.Println("Get user:1: <nil> (ok: false) - истёк!")
	} else {
		fmt.Printf("Get user:1: %v (ok: %v)\n", val, ok)
	}

	val, ok = cache.Get("user:2")
	fmt.Printf("Get user:2: %v (ok: %v) - ещё живой\n", val, ok)

	val, ok = cache.Get("config")
	fmt.Printf("Get config: %v (ok: %v) - без TTL\n", val, ok)

	// Тест 3: Cleanup
	fmt.Println("\n--- Тест 3: Cleanup ---")
	deleted := cache.Cleanup()
	fmt.Printf("Удалено истёкших записей: %d\n", deleted)
	fmt.Printf("Count после cleanup: %d\n", cache.Count())

	// Тест 4: Конкурентный доступ
	fmt.Println("\n--- Тест 4: Конкурентный доступ ---")

	cache = NewCache() // Новый кэш

	var wg sync.WaitGroup

	// Наполняем кэш данными
	for i := 0; i < 100; i++ {
		cache.Set(fmt.Sprintf("key:%d", i), i, time.Minute)
	}

	numWriters := 50
	numReaders := 200

	fmt.Printf("Запуск %d писателей и %d читателей...\n", numWriters, numReaders)

	// Писатели
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Set(fmt.Sprintf("writer:%d:key:%d", id, j), id*1000+j, time.Second)
			}
		}(i)
	}

	// Читатели
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Get(fmt.Sprintf("key:%d", j%100))
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Все операции завершены!")
	fmt.Println("Тест ПРОЙДЕН!")

	// Тест 5: Delete
	fmt.Println("\n--- Тест 5: Delete ---")
	cache = NewCache()
	cache.Set("to-delete", "value", time.Minute)

	val, ok = cache.Get("to-delete")
	fmt.Printf("До Delete: %v (ok: %v)\n", val, ok)

	deleted2 := cache.Delete("to-delete")
	fmt.Printf("Delete вернул: %v\n", deleted2)

	val, ok = cache.Get("to-delete")
	fmt.Printf("После Delete: %v (ok: %v)\n", val, ok)

	deleted3 := cache.Delete("non-existent")
	fmt.Printf("Delete несуществующего: %v\n", deleted3)

	// Напоминание
	fmt.Println("\n--- Проверка race detector ---")
	fmt.Println("Запусти программу с флагом -race:")
	fmt.Println("  go run -race main.go")
}
