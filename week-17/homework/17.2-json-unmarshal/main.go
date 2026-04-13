package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Address представляет адрес пользователя
type Address struct {
	City    string `json:"city"`
	Street  string `json:"street"`
	ZipCode string `json:"zip_code,omitempty"`
}

// User представляет пользователя системы
type User struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Age      int      `json:"age,omitempty"`
	IsActive bool     `json:"is_active"`
	Tags     []string `json:"tags,omitempty"`
	Address  *Address `json:"address,omitempty"` // указатель для опциональности
}

// ParseUser парсит JSON-строку в структуру User
func ParseUser(jsonData string) (*User, error) {
	// TODO: реализуй парсинг JSON в структуру User
	// 1. Создай переменную типа User
	// 2. Используй json.Unmarshal([]byte(jsonData), &user)
	// 3. Верни указатель на user и ошибку
	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ParseUsers парсит JSON-массив пользователей
func ParseUsers(jsonData string) ([]User, error) {
	// TODO: реализуй парсинг массива пользователей
	// 1. Создай слайс []User
	// 2. Используй json.Unmarshal
	// 3. Верни слайс и ошибку
	var users []User
	for _, user := range users {
		err := json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

// ParseDynamic парсит JSON неизвестной структуры в map
func ParseDynamic(jsonData string) (map[string]any, error) {
	// TODO: реализуй парсинг в map[string]any
	// Это полезно когда структура JSON заранее неизвестна
	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ExtractField безопасно извлекает значение по ключу
func ExtractField(data map[string]any, key string) (any, bool) {
	// TODO: извлеки значение из map по ключу
	// Верни значение и true если ключ существует
	// Верни nil и false если ключа нет
	if _, ok := data[key]; ok {
		value := data[key]
		return value, true
	}

	return nil, false
}

// ExtractNestedField извлекает вложенное значение по пути "key1.key2.key3"
func ExtractNestedField(data map[string]any, path string) (any, bool) {
	// TODO: реализуй извлечение вложенных полей
	// 1. Разбей path по точке: strings.Split(path, ".")
	// 2. Пройди по частям пути, на каждом шаге проверяя тип
	// 3. Для промежуточных ключей ожидай map[string]any
	// 4. Верни финальное значение или nil, false если путь не найден
	keys := strings.Split(path, ".")
	var current any = data

	for i, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		val, found := m[key]
		if !found {
			return nil, false
		}

		if i == len(keys)-1 {
			return val, true
		}

		current = val
	}

	return nil, false
}

// PrintUser выводит информацию о пользователе
func PrintUser(u *User) {
	fmt.Printf("  ID: %d\n", u.ID)
	fmt.Printf("  Username: %s\n", u.Username)
	fmt.Printf("  Email: %s\n", u.Email)
	fmt.Printf("  Active: %v\n", u.IsActive)

	if u.Age > 0 {
		fmt.Printf("  Age: %d\n", u.Age)
	}

	if len(u.Tags) > 0 {
		fmt.Printf("  Tags: %v\n", u.Tags)
	}

	if u.Address != nil {
		fmt.Printf("  Address: %s, %s", u.Address.City, u.Address.Street)
		if u.Address.ZipCode != "" {
			fmt.Printf(" (%s)", u.Address.ZipCode)
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("=== Парсинг JSON ===")

	// JSON с полными данными
	fullUserJSON := `{
		"id": 1,
		"username": "gopher",
		"email": "gopher@example.com",
		"age": 25,
		"is_active": true,
		"tags": ["developer", "golang", "backend"],
		"address": {
			"city": "Moscow",
			"street": "Tverskaya 1",
			"zip_code": "125009"
		}
	}`

	// JSON с минимальными данными
	minUserJSON := `{
		"id": 2,
		"username": "newbie",
		"email": "newbie@example.com",
		"is_active": false
	}`

	// Парсим полного пользователя
	fmt.Println("\n--- Полный пользователь ---")
	fullUser, err := ParseUser(fullUserJSON)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else if fullUser != nil {
		PrintUser(fullUser)
	}

	// Парсим минимального пользователя
	fmt.Println("\n--- Минимальный пользователь ---")
	minUser, err := ParseUser(minUserJSON)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else if minUser != nil {
		PrintUser(minUser)
	}

	// Парсим массив пользователей
	usersJSON := `[
		{"id": 1, "username": "alice", "email": "alice@mail.com", "is_active": true},
		{"id": 2, "username": "bob", "email": "bob@mail.com", "is_active": true},
		{"id": 3, "username": "charlie", "email": "charlie@mail.com", "is_active": false}
	]`

	fmt.Println("\n--- Массив пользователей ---")
	users, err := ParseUsers(usersJSON)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		for i, u := range users {
			fmt.Printf("%d. %s (%s) - active: %v\n", i+1, u.Username, u.Email, u.IsActive)
		}
	}

	// Парсим динамический JSON
	dynamicJSON := `{
		"type": "notification",
		"payload": {
			"message": "Привет, мир!",
			"priority": 5,
			"recipients": ["user1", "user2", "user3"]
		},
		"metadata": {
			"source": "api",
			"version": 2
		}
	}`

	fmt.Println("\n--- Динамический JSON ---")
	data, err := ParseDynamic(dynamicJSON)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else if data != nil {
		// Извлекаем простое поле
		if msgType, ok := ExtractField(data, "type"); ok {
			fmt.Printf("Type: %v\n", msgType)
		}

		// Извлекаем вложенное поле
		if message, ok := ExtractNestedField(data, "payload.message"); ok {
			fmt.Printf("Message: %v\n", message)
		}

		if priority, ok := ExtractNestedField(data, "payload.priority"); ok {
			// Внимание: числа из JSON всегда float64!
			fmt.Printf("Priority: %v (тип: %T)\n", priority, priority)
		}

		if source, ok := ExtractNestedField(data, "metadata.source"); ok {
			fmt.Printf("Source: %v\n", source)
		}

		// Несуществующий путь
		if _, ok := ExtractNestedField(data, "payload.nonexistent"); !ok {
			fmt.Println("Поле 'payload.nonexistent' не найдено (ожидаемо)")
		}
	}

	// Обработка ошибок парсинга
	fmt.Println("\n--- Обработка ошибок ---")
	invalidJSON := `{"id": "not a number", "username": 123}`
	_, err = ParseUser(invalidJSON)
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}
}
