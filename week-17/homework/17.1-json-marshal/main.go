package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Product представляет товар в магазине
// TODO: добавь json-теги для каждого поля
// - ID -> "id"
// - Name -> "name"
// - Price -> "price"
// - InStock -> "in_stock"
// - Description -> "description" с omitempty
type Product struct {
	ID          int
	Name        string
	Price       float64
	InStock     bool
	Description string
}

// Order представляет заказ покупателя
// TODO: добавь json-теги для каждого поля
// - OrderID -> "order_id"
// - Customer -> "customer"
// - Products -> "products"
// - Total -> "total"
// - CreatedAt -> "created_at"
type Order struct {
	OrderID   string
	Customer  string
	Products  []Product
	Total     float64
	CreatedAt time.Time
}

// MarshalProduct сериализует продукт в JSON с отступами
func MarshalProduct(p Product) (string, error) {
	// TODO: используй json.MarshalIndent для сериализации
	// Параметры: (v any, prefix string, indent string)
	// prefix = "", indent = "  " (два пробела)
	// Верни строку и ошибку
	return "", nil
}

// MarshalOrder сериализует заказ в JSON с отступами
func MarshalOrder(o Order) (string, error) {
	// TODO: используй json.MarshalIndent для сериализации заказа
	// Параметры такие же: prefix = "", indent = "  "
	return "", nil
}

// CalculateTotal вычисляет общую сумму заказа
func CalculateTotal(products []Product) float64 {
	// TODO: пройди по всем продуктам и просуммируй цены
	return 0
}

func main() {
	fmt.Println("=== Сериализация JSON ===")

	// Создаём продукты
	products := []Product{
		{
			ID:      1,
			Name:    "Ноутбук",
			Price:   89999.99,
			InStock: true,
			// Description пустое - должно быть пропущено в JSON
		},
		{
			ID:          2,
			Name:        "Мышь беспроводная",
			Price:       1299.50,
			InStock:     true,
			Description: "Эргономичная мышь с подсветкой",
		},
		{
			ID:          3,
			Name:        "Клавиатура",
			Price:       3499.00,
			InStock:     false,
			Description: "Механическая клавиатура",
		},
	}

	// Сериализуем отдельный продукт
	fmt.Println("\n--- Продукт без описания ---")
	json1, err := MarshalProduct(products[0])
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println(json1)
	}

	fmt.Println("\n--- Продукт с описанием ---")
	json2, err := MarshalProduct(products[1])
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println(json2)
	}

	// Создаём заказ
	order := Order{
		OrderID:   "ORD-2024-001",
		Customer:  "Иван Петров",
		Products:  products[:2], // только первые два продукта
		Total:     CalculateTotal(products[:2]),
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	fmt.Println("\n--- Заказ ---")
	orderJSON, err := MarshalOrder(order)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println(orderJSON)
	}

	// Компактный JSON (без отступов)
	fmt.Println("\n--- Компактный JSON ---")
	compact, err := json.Marshal(products[0])
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println(string(compact))
	}
}
