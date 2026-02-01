# ДЗ 19.3: Query Parameters - фильтрация, сортировка, пагинация

## Цель

Научиться работать с query parameters для реализации фильтрации, сортировки и пагинации в REST API. Освоить типичные паттерны обработки списковых endpoints.

## Что нужно сделать

Реализовать API для списка продуктов с поддержкой:

1. **Фильтрация** по категории, минимальной/максимальной цене, статусу наличия
2. **Сортировка** по имени, цене, дате создания
3. **Пагинация** с параметрами page и per_page
4. **Поиск** по названию продукта

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | /api/products | Список продуктов с фильтрами |
| GET | /api/products/{id} | Получить продукт по ID |

## Query Parameters

### Фильтрация

| Параметр | Тип | Описание | Пример |
|----------|-----|----------|--------|
| category | string | Фильтр по категории | `?category=electronics` |
| min_price | float | Минимальная цена | `?min_price=100` |
| max_price | float | Максимальная цена | `?max_price=500` |
| in_stock | bool | Только в наличии | `?in_stock=true` |
| search | string | Поиск по названию | `?search=phone` |

### Сортировка

| Параметр | Значения | Описание |
|----------|----------|----------|
| sort | name, price, created_at | Поле сортировки |
| order | asc, desc | Направление (по умолчанию asc) |

### Пагинация

| Параметр | По умолчанию | Описание |
|----------|--------------|----------|
| page | 1 | Номер страницы |
| per_page | 10 | Количество на странице (максимум 100) |

## Структуры данных

```go
// Product представляет продукт
type Product struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Category  string    `json:"category"`
    Price     float64   `json:"price"`
    InStock   bool      `json:"in_stock"`
    CreatedAt time.Time `json:"created_at"`
}

// ProductsResponse — ответ со списком продуктов и метаданными пагинации
type ProductsResponse struct {
    Products   []Product `json:"products"`
    Total      int       `json:"total"`       // общее количество (с учётом фильтров)
    Page       int       `json:"page"`        // текущая страница
    PerPage    int       `json:"per_page"`    // элементов на странице
    TotalPages int       `json:"total_pages"` // всего страниц
}

// ProductFilter содержит параметры фильтрации
type ProductFilter struct {
    Category string
    MinPrice float64
    MaxPrice float64
    InStock  *bool  // nil = не фильтровать, true/false = фильтровать
    Search   string
}

// SortParams содержит параметры сортировки
type SortParams struct {
    Field string // "name", "price", "created_at"
    Order string // "asc", "desc"
}

// PaginationParams содержит параметры пагинации
type PaginationParams struct {
    Page    int
    PerPage int
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
    Error string `json:"error"`
}
```

## Требования к GET /api/products

### Формат ответа

Всегда возвращает объект ProductsResponse:

```json
{
    "products": [
        {"id": 1, "name": "iPhone 15", "category": "electronics", "price": 999.99, "in_stock": true, "created_at": "2024-01-01T10:00:00Z"},
        {"id": 2, "name": "MacBook Pro", "category": "electronics", "price": 1999.99, "in_stock": true, "created_at": "2024-01-02T10:00:00Z"}
    ],
    "total": 15,
    "page": 1,
    "per_page": 10,
    "total_pages": 2
}
```

### Логика работы

1. **Сначала фильтрация** — отбираем продукты по заданным критериям
2. **Затем сортировка** — сортируем отфильтрованные результаты
3. **В конце пагинация** — выбираем нужную страницу

### Правила пагинации

- `page` < 1 → page = 1
- `per_page` < 1 → per_page = 10
- `per_page` > 100 → per_page = 100
- `total_pages` = ceil(total / per_page)
- Если запрошена страница за пределами → возвращается пустой массив products

## Примеры использования (curl)

```bash
# Все продукты (с пагинацией по умолчанию)
$ curl "http://localhost:8080/api/products"
{"products":[...],"total":20,"page":1,"per_page":10,"total_pages":2}

# Вторая страница
$ curl "http://localhost:8080/api/products?page=2"
{"products":[...],"total":20,"page":2,"per_page":10,"total_pages":2}

# 5 элементов на странице
$ curl "http://localhost:8080/api/products?per_page=5"
{"products":[...],"total":20,"page":1,"per_page":5,"total_pages":4}

# Фильтр по категории
$ curl "http://localhost:8080/api/products?category=electronics"
{"products":[...],"total":8,"page":1,"per_page":10,"total_pages":1}

# Фильтр по цене (от 100 до 500)
$ curl "http://localhost:8080/api/products?min_price=100&max_price=500"
{"products":[...],"total":5,"page":1,"per_page":10,"total_pages":1}

# Только в наличии
$ curl "http://localhost:8080/api/products?in_stock=true"
{"products":[...],"total":15,"page":1,"per_page":10,"total_pages":2}

# Поиск по названию (регистронезависимый)
$ curl "http://localhost:8080/api/products?search=phone"
{"products":[...],"total":3,"page":1,"per_page":10,"total_pages":1}

# Сортировка по цене (по убыванию)
$ curl "http://localhost:8080/api/products?sort=price&order=desc"
{"products":[...],"total":20,"page":1,"per_page":10,"total_pages":2}

# Сортировка по имени
$ curl "http://localhost:8080/api/products?sort=name"
{"products":[...],"total":20,"page":1,"per_page":10,"total_pages":2}

# Комбинация параметров
$ curl "http://localhost:8080/api/products?category=electronics&in_stock=true&sort=price&order=asc&page=1&per_page=5"
{"products":[...],"total":6,"page":1,"per_page":5,"total_pages":2}

# Получить конкретный продукт
$ curl "http://localhost:8080/api/products/1"
{"id":1,"name":"iPhone 15","category":"electronics","price":999.99,"in_stock":true,"created_at":"..."}

# Продукт не найден
$ curl -i "http://localhost:8080/api/products/999"
HTTP/1.1 404 Not Found
{"error":"product not found"}
```

## Вспомогательные функции

```go
// ParseFilter извлекает параметры фильтрации из запроса
func ParseFilter(r *http.Request) ProductFilter

// ParseSort извлекает параметры сортировки из запроса
func ParseSort(r *http.Request) SortParams

// ParsePagination извлекает параметры пагинации из запроса
func ParsePagination(r *http.Request) PaginationParams

// FilterProducts применяет фильтры к списку продуктов
func FilterProducts(products []Product, filter ProductFilter) []Product

// SortProducts сортирует продукты
func SortProducts(products []Product, params SortParams)

// Paginate возвращает срез для указанной страницы
func Paginate(products []Product, params PaginationParams) []Product

// CalculateTotalPages вычисляет количество страниц
func CalculateTotalPages(total, perPage int) int
```

## Подсказки

### Получение query параметров

```go
query := r.URL.Query()

// Одиночное значение
category := query.Get("category")    // "" если нет
page := query.Get("page")            // "" если нет

// Проверка наличия параметра
if _, ok := query["in_stock"]; ok {
    // параметр передан
}

// Конвертация в число
pageNum, _ := strconv.Atoi(query.Get("page"))

// Конвертация в float
price, _ := strconv.ParseFloat(query.Get("min_price"), 64)

// Конвертация в bool
inStock, _ := strconv.ParseBool(query.Get("in_stock"))
```

### Сортировка slice

```go
import "sort"

// Сортировка по полю
sort.Slice(products, func(i, j int) bool {
    if params.Order == "desc" {
        return products[i].Price > products[j].Price
    }
    return products[i].Price < products[j].Price
})
```

### Пагинация

```go
func Paginate(items []Product, page, perPage int) []Product {
    start := (page - 1) * perPage
    if start >= len(items) {
        return []Product{}
    }

    end := start + perPage
    if end > len(items) {
        end = len(items)
    }

    return items[start:end]
}
```

### Вычисление количества страниц

```go
func CalculateTotalPages(total, perPage int) int {
    if perPage <= 0 {
        return 0
    }
    return (total + perPage - 1) / perPage
}
```

## Критерии выполнения

- [ ] GET /api/products возвращает ProductsResponse с метаданными пагинации
- [ ] Работает фильтрация по category
- [ ] Работает фильтрация по min_price и max_price
- [ ] Работает фильтрация по in_stock
- [ ] Работает поиск по search (регистронезависимый)
- [ ] Работает сортировка по name, price, created_at
- [ ] Работает сортировка в обоих направлениях (asc/desc)
- [ ] Работает пагинация с page и per_page
- [ ] per_page ограничен максимумом 100
- [ ] total_pages вычисляется корректно
- [ ] Фильтры, сортировка и пагинация работают вместе
- [ ] GET /api/products/{id} возвращает продукт или 404
- [ ] Все JSON ответы имеют Content-Type: application/json
- [ ] Код компилируется без ошибок
