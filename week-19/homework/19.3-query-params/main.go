package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}

// ProductFilter содержит параметры фильтрации
type ProductFilter struct {
	Category string
	MinPrice float64
	MaxPrice float64
	InStock  *bool // nil = не фильтровать
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

// ProductStorage — хранилище продуктов
type ProductStorage struct {
	mu       sync.RWMutex
	products map[int]Product
}

// NewProductStorage создаёт хранилище с тестовыми данными
func NewProductStorage() *ProductStorage {
	baseTime := time.Now().Add(-24 * time.Hour)

	products := map[int]Product{
		1:  {ID: 1, Name: "iPhone 15", Category: "electronics", Price: 999.99, InStock: true, CreatedAt: baseTime.Add(1 * time.Hour)},
		2:  {ID: 2, Name: "MacBook Pro", Category: "electronics", Price: 1999.99, InStock: true, CreatedAt: baseTime.Add(2 * time.Hour)},
		3:  {ID: 3, Name: "Samsung Galaxy S24", Category: "electronics", Price: 899.99, InStock: true, CreatedAt: baseTime.Add(3 * time.Hour)},
		4:  {ID: 4, Name: "Sony Headphones", Category: "electronics", Price: 299.99, InStock: false, CreatedAt: baseTime.Add(4 * time.Hour)},
		5:  {ID: 5, Name: "Dell Monitor", Category: "electronics", Price: 449.99, InStock: true, CreatedAt: baseTime.Add(5 * time.Hour)},
		6:  {ID: 6, Name: "Office Chair", Category: "furniture", Price: 299.99, InStock: true, CreatedAt: baseTime.Add(6 * time.Hour)},
		7:  {ID: 7, Name: "Standing Desk", Category: "furniture", Price: 599.99, InStock: true, CreatedAt: baseTime.Add(7 * time.Hour)},
		8:  {ID: 8, Name: "Bookshelf", Category: "furniture", Price: 149.99, InStock: false, CreatedAt: baseTime.Add(8 * time.Hour)},
		9:  {ID: 9, Name: "Coffee Table", Category: "furniture", Price: 199.99, InStock: true, CreatedAt: baseTime.Add(9 * time.Hour)},
		10: {ID: 10, Name: "Desk Lamp", Category: "furniture", Price: 49.99, InStock: true, CreatedAt: baseTime.Add(10 * time.Hour)},
		11: {ID: 11, Name: "Running Shoes", Category: "sports", Price: 129.99, InStock: true, CreatedAt: baseTime.Add(11 * time.Hour)},
		12: {ID: 12, Name: "Yoga Mat", Category: "sports", Price: 29.99, InStock: true, CreatedAt: baseTime.Add(12 * time.Hour)},
		13: {ID: 13, Name: "Dumbbells Set", Category: "sports", Price: 199.99, InStock: false, CreatedAt: baseTime.Add(13 * time.Hour)},
		14: {ID: 14, Name: "Tennis Racket", Category: "sports", Price: 89.99, InStock: true, CreatedAt: baseTime.Add(14 * time.Hour)},
		15: {ID: 15, Name: "Basketball", Category: "sports", Price: 24.99, InStock: true, CreatedAt: baseTime.Add(15 * time.Hour)},
		16: {ID: 16, Name: "Cooking Pan Set", Category: "kitchen", Price: 149.99, InStock: true, CreatedAt: baseTime.Add(16 * time.Hour)},
		17: {ID: 17, Name: "Coffee Maker", Category: "kitchen", Price: 79.99, InStock: true, CreatedAt: baseTime.Add(17 * time.Hour)},
		18: {ID: 18, Name: "Blender", Category: "kitchen", Price: 59.99, InStock: false, CreatedAt: baseTime.Add(18 * time.Hour)},
		19: {ID: 19, Name: "Microwave Oven", Category: "kitchen", Price: 199.99, InStock: true, CreatedAt: baseTime.Add(19 * time.Hour)},
		20: {ID: 20, Name: "Electric Kettle", Category: "kitchen", Price: 39.99, InStock: true, CreatedAt: baseTime.Add(20 * time.Hour)},
	}

	return &ProductStorage{products: products}
}

// Get возвращает продукт по ID
func (s *ProductStorage) Get(id int) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	product, exists := s.products[id]
	return product, exists
}

// GetAll возвращает все продукты
func (s *ProductStorage) GetAll() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}
	return products
}

// Глобальное хранилище
var storage = NewProductStorage()

// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, status int, data any) {
	// TODO: реализуй функцию
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("fail: %v", err)
	}
}

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, status int, message string) {
	// TODO: реализуй функцию
	writeJSON(w, status, ErrorResponse{Error: message})
}

// ParseFilter извлекает параметры фильтрации из запроса
func ParseFilter(r *http.Request) ProductFilter {
	// TODO: реализуй функцию
	// query := r.URL.Query()
	//
	// filter := ProductFilter{
	//     Category: query.Get("category"),
	//     Search:   query.Get("search"),
	// }
	//
	// // Парсинг min_price
	// if minStr := query.Get("min_price"); minStr != "" {
	//     if min, err := strconv.ParseFloat(minStr, 64); err == nil {
	//         filter.MinPrice = min
	//     }
	// }
	//
	// // Парсинг max_price
	// if maxStr := query.Get("max_price"); maxStr != "" {
	//     if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
	//         filter.MaxPrice = max
	//     }
	// }
	//
	// // Парсинг in_stock (только если параметр передан)
	// if _, ok := query["in_stock"]; ok {
	//     if inStock, err := strconv.ParseBool(query.Get("in_stock")); err == nil {
	//         filter.InStock = &inStock
	//     }
	// }
	//
	// return filter
	query := r.URL.Query()

	filter := ProductFilter{
		Category: query.Get("category"),
		Search:   query.Get("search"),
	}

	if minStr := query.Get("min_price"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.MinPrice = min
		}
	}

	if maxStr := query.Get("max_price"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.MaxPrice = max
		}
	}

	if _, ok := query["in_stock"]; ok {
		if inStock, err := strconv.ParseBool(query.Get("in_stock")); err == nil {
			filter.InStock = &inStock
		}
	}

	return filter
}

// ParseSort извлекает параметры сортировки из запроса
func ParseSort(r *http.Request) SortParams {
	// TODO: реализуй функцию
	// query := r.URL.Query()
	//
	// params := SortParams{
	//     Field: query.Get("sort"),
	//     Order: query.Get("order"),
	// }
	//
	// // Валидация поля сортировки
	// validFields := map[string]bool{"name": true, "price": true, "created_at": true}
	// if !validFields[params.Field] {
	//     params.Field = ""
	// }
	//
	// // Значение по умолчанию для order
	// if params.Order != "desc" {
	//     params.Order = "asc"
	// }
	//
	// return params
	query := r.URL.Query()

	params := SortParams{
		Field: query.Get("sort"),
		Order: query.Get("order"),
	}

	validFields := map[string]bool{"name": true, "price": true, "created_at": true}
	if !validFields[params.Field] {
		params.Field = ""
	}

	if params.Order != "desc" {
		params.Order = "asc"
	}

	return params
}

// ParsePagination извлекает параметры пагинации из запроса
func ParsePagination(r *http.Request) PaginationParams {
	// TODO: реализуй функцию
	// query := r.URL.Query()
	//
	// page, _ := strconv.Atoi(query.Get("page"))
	// if page < 1 {
	//     page = 1
	// }
	//
	// perPage, _ := strconv.Atoi(query.Get("per_page"))
	// if perPage < 1 {
	//     perPage = 10
	// }
	// if perPage > 100 {
	//     perPage = 100
	// }
	//
	// return PaginationParams{Page: page, PerPage: perPage}
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(query.Get("per_page"))
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100

	}

	return PaginationParams{Page: page, PerPage: perPage}
}

// FilterProducts применяет фильтры к списку продуктов
func FilterProducts(products []Product, filter ProductFilter) []Product {
	// TODO: реализуй функцию
	// result := make([]Product, 0)
	//
	// for _, p := range products {
	//     // Фильтр по категории
	//     if filter.Category != "" && p.Category != filter.Category {
	//         continue
	//     }
	//
	//     // Фильтр по минимальной цене
	//     if filter.MinPrice > 0 && p.Price < filter.MinPrice {
	//         continue
	//     }
	//
	//     // Фильтр по максимальной цене
	//     if filter.MaxPrice > 0 && p.Price > filter.MaxPrice {
	//         continue
	//     }
	//
	//     // Фильтр по наличию
	//     if filter.InStock != nil && p.InStock != *filter.InStock {
	//         continue
	//     }
	//
	//     // Поиск по названию (регистронезависимый)
	//     if filter.Search != "" {
	//         if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(filter.Search)) {
	//             continue
	//         }
	//     }
	//
	//     result = append(result, p)
	// }
	//
	// return result
	result := make([]Product, 0)

	for _, p := range products {
		if filter.Category != "" && p.Category != filter.Category {
			continue
		}

		if filter.MinPrice > 0 && p.Price < filter.MinPrice {
			continue
		}

		if filter.MaxPrice > 0 && p.Price > filter.MaxPrice {
			continue
		}

		if filter.InStock != nil && p.InStock != *filter.InStock {
			continue
		}

		if filter.Search != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(filter.Search)) {
				continue
			}
		}

		result = append(result, p)
	}

	return result
}

// SortProducts сортирует продукты
func SortProducts(products []Product, params SortParams) {
	// TODO: реализуй функцию
	// if params.Field == "" {
	//     return
	// }
	//
	// sort.Slice(products, func(i, j int) bool {
	//     var less bool
	//
	//     switch params.Field {
	//     case "name":
	//         less = products[i].Name < products[j].Name
	//     case "price":
	//         less = products[i].Price < products[j].Price
	//     case "created_at":
	//         less = products[i].CreatedAt.Before(products[j].CreatedAt)
	//     default:
	//         return false
	//     }
	//
	//     if params.Order == "desc" {
	//         return !less
	//     }
	//     return less
	// })
	if params.Field == "" {
		return
	}

	sort.SliceStable(products, func(i, j int) bool {
		var less bool

		switch params.Field {
		case "name":
			less = products[i].Name < products[j].Name
		case "price":
			less = products[i].Price < products[j].Price
		case "created_at":
			less = products[i].CreatedAt.Before(products[j].CreatedAt)
		default:
			return false
		}

		if params.Order == "desc" {
			return !less
		}

		return less
	})
}

// Paginate возвращает срез для указанной страницы
func Paginate(products []Product, params PaginationParams) []Product {
	// TODO: реализуй функцию
	// start := (params.Page - 1) * params.PerPage
	// if start >= len(products) {
	//     return []Product{}
	// }
	//
	// end := start + params.PerPage
	// if end > len(products) {
	//     end = len(products)
	// }
	//
	// return products[start:end]
	start := (params.Page - 1) * params.PerPage
	if start >= len(products) {
		return []Product{}
	}

	end := start + params.PerPage
	if end > len(products) {
		end = len(products)
	}

	return products[start:end]
}

// CalculateTotalPages вычисляет количество страниц
func CalculateTotalPages(total, perPage int) int {
	// TODO: реализуй функцию
	// if perPage <= 0 {
	//     return 0
	// }
	// return (total + perPage - 1) / perPage
	if perPage <= 0 {
		return 0
	}
	return (total + perPage - 1) / perPage
}

// listProductsHandler обрабатывает GET /api/products
func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи все продукты: products := storage.GetAll()
	// 2. Распарси фильтры: filter := ParseFilter(r)
	// 3. Распарси сортировку: sortParams := ParseSort(r)
	// 4. Распарси пагинацию: pagination := ParsePagination(r)
	// 5. Применить фильтры: filtered := FilterProducts(products, filter)
	// 6. Сохрани total ДО пагинации: total := len(filtered)
	// 7. Отсортируй: SortProducts(filtered, sortParams)
	// 8. Пагинируй: page := Paginate(filtered, pagination)
	// 9. Собери ответ:
	//    response := ProductsResponse{
	//        Products:   page,
	//        Total:      total,
	//        Page:       pagination.Page,
	//        PerPage:    pagination.PerPage,
	//        TotalPages: CalculateTotalPages(total, pagination.PerPage),
	//    }
	// 10. Отправь: writeJSON(w, http.StatusOK, response)
	products := storage.GetAll()
	filter := ParseFilter(r)
	sortParams := ParseSort(r)
	pagination := ParsePagination(r)
	filtered := FilterProducts(products, filter)
	total := len(filtered)
	SortProducts(filtered, sortParams)
	page := Paginate(filtered, pagination)
	response := ProductsResponse{
		Products:   page,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: CalculateTotalPages(total, pagination.PerPage),
	}
	writeJSON(w, http.StatusOK, response)
}

// getProductHandler обрабатывает GET /api/products/{id}
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обработчик
	// 1. Получи id из пути
	// 2. Конвертируй в int
	// 3. Получи продукт или верни 404
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, exists := storage.Get(id)
	if !exists {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики
	// mux.HandleFunc("GET /api/products", listProductsHandler)
	// mux.HandleFunc("GET /api/products/{id}", getProductHandler)

	mux.HandleFunc("GET /api/products", listProductsHandler)
	mux.HandleFunc("GET /api/products/{id}", getProductHandler)

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET /api/products      - List products with filters")
	fmt.Println("  GET /api/products/{id} - Get product by ID")
	fmt.Println()
	fmt.Println("Query Parameters:")
	fmt.Println("  Filtering: category, min_price, max_price, in_stock, search")
	fmt.Println("  Sorting:   sort (name|price|created_at), order (asc|desc)")
	fmt.Println("  Pagination: page, per_page (max 100)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  /api/products?category=electronics&in_stock=true")
	fmt.Println("  /api/products?min_price=100&max_price=500&sort=price&order=desc")
	fmt.Println("  /api/products?search=phone&page=1&per_page=5")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
