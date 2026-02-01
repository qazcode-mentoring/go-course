package main

import (
	"fmt"
	"time"
)

// cacheEntry хранит значение и время истечения
type cacheEntry[V any] struct {
	value      V
	expiration time.Time
}

// isExpired проверяет, истекла ли запись
func (e cacheEntry[V]) isExpired() bool {
	return time.Now().After(e.expiration)
}

// Cache - универсальный кэш с TTL
// K - тип ключа (должен быть comparable для использования в map)
// V - тип значения (любой)
type Cache[K comparable, V any] struct {
	data       map[K]cacheEntry[V]
	defaultTTL time.Duration
}

// NewCache создает новый кэш с указанным TTL по умолчанию
func NewCache[K comparable, V any](defaultTTL time.Duration) *Cache[K, V] {
	// TODO: реализуй функцию
	// Создай и верни указатель на новый Cache с инициализированной map
	return nil
}

// Set добавляет значение с дефолтным TTL
func (c *Cache[K, V]) Set(key K, value V) {
	// TODO: реализуй метод
	// Используй c.defaultTTL для установки времени истечения
	// Время истечения = time.Now().Add(c.defaultTTL)
}

// SetWithTTL добавляет значение с указанным TTL
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	// TODO: реализуй метод
	// Создай cacheEntry с value и expiration = time.Now().Add(ttl)
	// Сохрани в c.data[key]
}

// Get возвращает значение, если оно существует и не истекло
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: реализуй метод
	// 1. Проверь, существует ли ключ в c.data
	// 2. Если существует, проверь не истекла ли запись
	// 3. Если истекла, можно удалить и вернуть (zero, false)
	// 4. Если валидна, верни (value, true)
	var zero V
	return zero, false
}

// Delete удаляет значение из кэша
func (c *Cache[K, V]) Delete(key K) {
	// TODO: реализуй метод
	// Используй delete(c.data, key)
}

// Len возвращает количество валидных (не истекших) записей
func (c *Cache[K, V]) Len() int {
	// TODO: реализуй метод
	// Пройди по всем записям и посчитай только те, которые не истекли
	return 0
}

// Clear очищает весь кэш
func (c *Cache[K, V]) Clear() {
	// TODO: реализуй метод
	// Можно создать новую пустую map: c.data = make(map[K]cacheEntry[V])
}

// Cleanup удаляет все истекшие записи
func (c *Cache[K, V]) Cleanup() {
	// TODO: реализуй метод
	// Пройди по всем записям и удали те, которые истекли
}

// Пример структуры для тестирования
type UserSession struct {
	UserID    int
	Token     string
	CreatedAt time.Time
}

func main() {
	fmt.Println("=== Generic Cache с TTL ===")

	// Базовый тест с string -> string
	fmt.Println("\n--- Базовый тест ---")
	cache := NewCache[string, string](5 * time.Second)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, ok := cache.Get("key1"); ok {
		fmt.Printf("key1 = %s (ожидается: value1)\n", val)
	}

	fmt.Printf("Размер кэша: %d (ожидается: 2)\n", cache.Len())

	// Тест Delete
	fmt.Println("\n--- Тест Delete ---")
	cache.Delete("key1")
	if _, ok := cache.Get("key1"); !ok {
		fmt.Println("key1 удален (ожидается)")
	}
	fmt.Printf("Размер после Delete: %d (ожидается: 1)\n", cache.Len())

	// Тест с коротким TTL
	fmt.Println("\n--- Тест TTL ---")
	shortCache := NewCache[string, string](100 * time.Millisecond)
	shortCache.Set("temp", "data")

	if val, ok := shortCache.Get("temp"); ok {
		fmt.Printf("До истечения: %s (ожидается: data)\n", val)
	}

	fmt.Println("Ждем 150ms...")
	time.Sleep(150 * time.Millisecond)

	if _, ok := shortCache.Get("temp"); !ok {
		fmt.Println("После истечения: ключ недоступен (ожидается)")
	}

	// Тест SetWithTTL
	fmt.Println("\n--- Тест SetWithTTL ---")
	mixedCache := NewCache[string, string](5 * time.Second)
	mixedCache.Set("long", "5 seconds")
	mixedCache.SetWithTTL("short", "100ms", 100*time.Millisecond)

	fmt.Printf("Размер до истечения: %d (ожидается: 2)\n", mixedCache.Len())

	time.Sleep(150 * time.Millisecond)

	fmt.Printf("Размер после 150ms: %d (ожидается: 1)\n", mixedCache.Len())

	if _, ok := mixedCache.Get("short"); !ok {
		fmt.Println("short истек (ожидается)")
	}
	if val, ok := mixedCache.Get("long"); ok {
		fmt.Printf("long еще валиден: %s (ожидается: 5 seconds)\n", val)
	}

	// Тест Cleanup
	fmt.Println("\n--- Тест Cleanup ---")
	cleanupCache := NewCache[int, string](50 * time.Millisecond)
	cleanupCache.Set(1, "one")
	cleanupCache.Set(2, "two")
	cleanupCache.Set(3, "three")

	time.Sleep(100 * time.Millisecond)

	// Добавим свежие данные
	cleanupCache.Set(4, "four")
	cleanupCache.Set(5, "five")

	fmt.Println("Перед Cleanup (3 истекших + 2 свежих)")
	cleanupCache.Cleanup()
	fmt.Printf("После Cleanup: %d записей (ожидается: 2)\n", cleanupCache.Len())

	// Тест Clear
	fmt.Println("\n--- Тест Clear ---")
	clearCache := NewCache[string, int](time.Minute)
	clearCache.Set("a", 1)
	clearCache.Set("b", 2)
	clearCache.Set("c", 3)

	fmt.Printf("До Clear: %d (ожидается: 3)\n", clearCache.Len())
	clearCache.Clear()
	fmt.Printf("После Clear: %d (ожидается: 0)\n", clearCache.Len())

	// Тест с пользовательской структурой
	fmt.Println("\n--- Тест с UserSession ---")
	sessionCache := NewCache[string, UserSession](time.Hour)

	session := UserSession{
		UserID:    42,
		Token:     "abc123xyz",
		CreatedAt: time.Now(),
	}
	sessionCache.Set("session:user42", session)

	if s, ok := sessionCache.Get("session:user42"); ok {
		fmt.Printf("UserID: %d, Token: %s\n", s.UserID, s.Token)
	}

	// Тест с числовым ключом
	fmt.Println("\n--- Тест с int ключом ---")
	intKeyCache := NewCache[int, []string](time.Minute)
	intKeyCache.Set(1, []string{"apple", "banana"})
	intKeyCache.Set(2, []string{"cat", "dog"})

	if val, ok := intKeyCache.Get(1); ok {
		fmt.Printf("Ключ 1: %v\n", val)
	}

	fmt.Println("\n=== Все тесты выполнены ===")
}
