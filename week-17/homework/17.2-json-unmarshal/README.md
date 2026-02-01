# ДЗ 17.2: Парсинг JSON

## Цель

Научиться парсить JSON в Go-структуры, работать с опциональными полями и динамическим JSON неизвестной структуры.

## Что нужно сделать

1. Распарсить JSON в структуру `User`
2. Обработать опциональные поля (могут отсутствовать в JSON)
3. Распарсить JSON с неизвестной структурой в `map[string]any`
4. Извлечь вложенные данные из динамического JSON

## Структуры

```go
type Address struct {
    City    string `json:"city"`
    Street  string `json:"street"`
    ZipCode string `json:"zip_code,omitempty"`
}

type User struct {
    ID        int      `json:"id"`
    Username  string   `json:"username"`
    Email     string   `json:"email"`
    Age       int      `json:"age,omitempty"`       // опционально
    IsActive  bool     `json:"is_active"`
    Tags      []string `json:"tags,omitempty"`      // опционально
    Address   *Address `json:"address,omitempty"`   // опционально, указатель!
}
```

## JSON для парсинга

### Полный пользователь
```json
{
  "id": 1,
  "username": "gopher",
  "email": "gopher@example.com",
  "age": 25,
  "is_active": true,
  "tags": ["developer", "golang"],
  "address": {
    "city": "Moscow",
    "street": "Tverskaya 1",
    "zip_code": "125009"
  }
}
```

### Минимальный пользователь (без опциональных полей)
```json
{
  "id": 2,
  "username": "newbie",
  "email": "newbie@example.com",
  "is_active": false
}
```

### Динамический JSON (структура заранее неизвестна)
```json
{
  "type": "notification",
  "payload": {
    "message": "Hello!",
    "priority": 5,
    "recipients": ["user1", "user2"]
  },
  "metadata": {
    "source": "api",
    "version": 2
  }
}
```

## Функции для реализации

```go
// ParseUser парсит JSON в структуру User
func ParseUser(jsonData string) (*User, error)

// ParseUsers парсит JSON-массив пользователей
func ParseUsers(jsonData string) ([]User, error)

// ParseDynamic парсит JSON неизвестной структуры
func ParseDynamic(jsonData string) (map[string]any, error)

// ExtractField извлекает значение по ключу из динамического JSON
func ExtractField(data map[string]any, key string) (any, bool)

// ExtractNestedField извлекает вложенное значение (например, "payload.message")
func ExtractNestedField(data map[string]any, path string) (any, bool)
```

## Пример использования

```go
func main() {
    jsonData := `{"id":1,"username":"test","email":"test@mail.com","is_active":true}`

    user, err := ParseUser(jsonData)
    if err != nil {
        panic(err)
    }

    fmt.Printf("User: %s (email: %s)\n", user.Username, user.Email)

    // Проверка опционального поля
    if user.Address != nil {
        fmt.Printf("City: %s\n", user.Address.City)
    } else {
        fmt.Println("Адрес не указан")
    }
}
```

## Подсказки

- `json.Unmarshal([]byte(jsonData), &result)` — передаём указатель!
- Для опциональных структур используй указатель (`*Address`), чтобы отличить отсутствие от пустого значения
- При работе с `map[string]any` числа всегда будут `float64`
- Для извлечения вложенных полей разбей путь по точке: `strings.Split(path, ".")`
- Type assertion: `value.(string)`, `value.(float64)`, `value.(map[string]any)`

## Важно: Type assertions

```go
data := map[string]any{"count": 42}

// Опасно — паника если тип неверный
count := data["count"].(int) // PANIC! число будет float64

// Безопасно — проверяем тип
if count, ok := data["count"].(float64); ok {
    fmt.Println(int(count)) // 42
}
```

## Критерии выполнения

- [ ] `ParseUser` корректно парсит полный JSON
- [ ] `ParseUser` корректно обрабатывает отсутствующие опциональные поля
- [ ] `ParseUsers` парсит массив пользователей
- [ ] `ParseDynamic` парсит JSON в `map[string]any`
- [ ] `ExtractField` безопасно извлекает значения с проверкой наличия
- [ ] `ExtractNestedField` работает с вложенными путями ("payload.message")
- [ ] Код обрабатывает ошибки парсинга
- [ ] Код компилируется без ошибок
