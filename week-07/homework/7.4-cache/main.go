package main

import "fmt"

// Cache представляет простой кэш ключ-значение
type Cache struct {
	data map[string]string
}

// NewCache создаёт и возвращает новый кэш
func NewCache() *Cache {
	// TODO: реализуй функцию
	// Не забудь инициализировать map через make
	return &Cache{data: make(map[string]string)}
}

// Set сохраняет значение по ключу
func (c *Cache) Set(key, value string) {
	// TODO: реализуй метод
	c.data[key] = value
}

// Get возвращает значение по ключу и флаг существования
func (c *Cache) Get(key string) (string, bool) {
	// TODO: реализуй метод
	value, ok := c.data[key]
	return value, ok
}

// Delete удаляет значение по ключу
func (c *Cache) Delete(key string) {
	// TODO: реализуй метод
	delete(c.data, key)
}

// Clear очищает весь кэш
func (c *Cache) Clear() {
	// TODO: реализуй метод
	// Можно создать новый пустой map
	c.data = make(map[string]string)
}

// Size возвращает количество элементов в кэше
func (c *Cache) Size() int {
	// TODO: реализуй метод
	return len(c.data)
}

func main() {
	cache := NewCache()

	fmt.Println("=== Простой кэш ===")

	// Добавление
	cache.Set("user:1", "Алексей")
	cache.Set("user:2", "Мария")
	cache.Set("user:3", "Иван")
	fmt.Println("Добавлено 3 элемента, размер:", cache.Size())

	// Получение
	if name, ok := cache.Get("user:1"); ok {
		fmt.Println("user:1 =", name)
	}

	// Проверка несуществующего
	if _, ok := cache.Get("user:999"); !ok {
		fmt.Println("user:999 не найден")
	}

	// Удаление
	cache.Delete("user:2")
	fmt.Println("После удаления user:2, размер:", cache.Size())

	// Очистка
	cache.Clear()
	fmt.Println("После очистки, размер:", cache.Size())
}
