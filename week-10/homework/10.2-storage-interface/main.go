package main

import (
	"fmt"
)

// Storage — интерфейс для хранилища ключ-значение
type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
	Keys() []string
}

// MemoryStorage — реализация Storage на основе map
type MemoryStorage struct {
	data map[string]string
}

// NewMemoryStorage создаёт новое хранилище на основе map
func NewMemoryStorage() *MemoryStorage {
	// TODO: реализуй функцию
	return &MemoryStorage{data: make(map[string]string)}
}

// Set сохраняет значение по ключу
func (m *MemoryStorage) Set(key, value string) {
	// TODO: реализуй метод
	m.data[key] = value
}

// Get возвращает значение по ключу
func (m *MemoryStorage) Get(key string) (string, bool) {
	// TODO: реализуй метод
	val, ok := m.data[key]
	return val, ok
}

// Delete удаляет значение по ключу
func (m *MemoryStorage) Delete(key string) {
	// TODO: реализуй метод
	delete(m.data, key)
}

// Keys возвращает все ключи
func (m *MemoryStorage) Keys() []string {
	// TODO: реализуй метод
	keys := make([]string, 0, len(m.data))

	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// KeyValue — пара ключ-значение для SliceStorage
type KeyValue struct {
	Key, Value string
}

// SliceStorage — реализация Storage на основе слайса
type SliceStorage struct {
	items []KeyValue
}

// NewSliceStorage создаёт новое хранилище на основе слайса
func NewSliceStorage() *SliceStorage {
	// TODO: реализуй функцию
	return &SliceStorage{items: make([]KeyValue, 0)}
}

// Set сохраняет значение (обновляет если ключ существует)
func (s *SliceStorage) Set(key, value string) {
	// TODO: реализуй метод
	// Если ключ существует — обнови значение
	// Иначе — добавь новую пару
	existKey := false
	for i := range s.items {
		if s.items[i].Key == key {
			existKey = true
			s.items[i].Value = value
		}
	}
	if !existKey {
		s.items = append(s.items, KeyValue{key, value})
	}
}

// Get возвращает значение по ключу
func (s *SliceStorage) Get(key string) (string, bool) {
	// TODO: реализуй метод
	for i := range s.items {
		if s.items[i].Key == key {
			return s.items[i].Value, true
		}
	}
	return "", false
}

// Delete удаляет значение по ключу
func (s *SliceStorage) Delete(key string) {
	// TODO: реализуй метод
	for i := range s.items {
		if s.items[i].Key == key {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return
		}
	}
}

// Keys возвращает все ключи
func (s *SliceStorage) Keys() []string {
	// TODO: реализуй метод
	keys := make([]string, 0)
	for i := range s.items {
		keys = append(keys, s.items[i].Key)
	}
	return keys
}

// CopyStorage копирует все данные из одного Storage в другой
func CopyStorage(from, to Storage) {
	// TODO: реализуй функцию
	// Получи все ключи из from
	// Для каждого ключа скопируй значение в to
	keys := from.Keys()

	for _, key := range keys {
		if val, ok := from.Get(key); ok {
			to.Set(key, val)
		}
	}
}

func main() {
	fmt.Println("=== Интерфейс Storage ===")

	// Работа с MemoryStorage
	fmt.Println("\n--- MemoryStorage ---")
	mem := NewMemoryStorage()
	mem.Set("name", "Алексей")
	mem.Set("city", "Москва")
	mem.Set("lang", "Go")

	if val, ok := mem.Get("name"); ok {
		fmt.Println("name:", val)
	}
	fmt.Println("Ключи:", mem.Keys())

	// Работа с SliceStorage
	fmt.Println("\n--- SliceStorage ---")
	slice := NewSliceStorage()
	slice.Set("a", "1")
	slice.Set("b", "2")
	fmt.Println("Ключи:", slice.Keys())

	// Копирование между хранилищами
	fmt.Println("\n--- CopyStorage ---")
	fmt.Println("Копируем из MemoryStorage в SliceStorage")
	CopyStorage(mem, slice)
	fmt.Println("SliceStorage после копирования:")
	fmt.Println("Ключи:", slice.Keys())

	if val, ok := slice.Get("city"); ok {
		fmt.Println("city:", val)
	}
}
