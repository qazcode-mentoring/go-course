# ДЗ 17.1: Сериализация структур в JSON

## Цель

Научиться преобразовывать Go-структуры в JSON с использованием struct tags для управления именами полей и форматом вывода.

## Что нужно сделать

1. Создать структуру `Product` и сериализовать её в JSON
2. Создать структуру `Order` со слайсом продуктов
3. Использовать struct tags для кастомных имён полей
4. Реализовать форматированный вывод с отступами

## Структуры

```go
type Product struct {
    ID          int     // json: "id"
    Name        string  // json: "name"
    Price       float64 // json: "price"
    InStock     bool    // json: "in_stock"
    Description string  // json: "description", пропустить если пустое
}

type Order struct {
    OrderID   string    // json: "order_id"
    Customer  string    // json: "customer"
    Products  []Product // json: "products"
    Total     float64   // json: "total"
    CreatedAt time.Time // json: "created_at"
}
```

## Требования к struct tags

- Все поля должны иметь snake_case имена в JSON
- Поле `Description` должно пропускаться, если оно пустое (`omitempty`)
- Используй `json.MarshalIndent` для красивого вывода

## Пример использования

```go
func main() {
    product := Product{
        ID:      1,
        Name:    "Ноутбук",
        Price:   89999.99,
        InStock: true,
        // Description пустое - не должно быть в JSON
    }

    data, err := json.MarshalIndent(product, "", "  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}
```

## Ожидаемый вывод

```json
{
  "id": 1,
  "name": "Ноутбук",
  "price": 89999.99,
  "in_stock": true
}
```

Для заказа:
```json
{
  "order_id": "ORD-001",
  "customer": "Иван Петров",
  "products": [
    {
      "id": 1,
      "name": "Ноутбук",
      "price": 89999.99,
      "in_stock": true
    },
    {
      "id": 2,
      "name": "Мышь",
      "price": 1299.50,
      "in_stock": true,
      "description": "Беспроводная мышь"
    }
  ],
  "total": 91299.49,
  "created_at": "2024-01-15T10:30:00Z"
}
```

## Подсказки

- Struct tags пишутся после типа в обратных кавычках: `` `json:"name"` ``
- Для omitempty: `` `json:"field,omitempty"` ``
- `json.MarshalIndent(v, prefix, indent)` — prefix обычно пустой, indent обычно 2 или 4 пробела
- Для time.Time Go автоматически использует формат RFC3339

## Критерии выполнения

- [ ] Структура `Product` имеет правильные json-теги
- [ ] Структура `Order` содержит слайс продуктов
- [ ] Поле `Description` пропускается когда пустое (omitempty)
- [ ] JSON выводится с отступами (MarshalIndent)
- [ ] Имена полей в JSON используют snake_case
- [ ] Код компилируется и выполняется без ошибок
