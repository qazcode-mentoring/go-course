package main

import "fmt"

// IPAddr представляет IPv4 адрес
type IPAddr [4]byte

// String возвращает строковое представление IP адреса
func (ip IPAddr) String() string {
	// TODO: реализуй метод
	// Формат: "a.b.c.d"
	// Используй fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// Duration представляет длительность в секундах
type Duration int

// String возвращает строковое представление длительности
func (d Duration) String() string {
	// TODO: реализуй метод
	// Формат: "Xh Xm Xs"
	// Раздели на часы (d / 3600), минуты ((d % 3600) / 60), секунды (d % 60)
	h := d / 3600
	m := (d % 3600) / 60
	s := d % 60
	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}

// Book представляет информацию о книге
type Book struct {
	Title  string
	Author string
	Year   int
}

// String возвращает строковое представление книги
func (b Book) String() string {
	// TODO: реализуй метод
	// Формат: "«Title» (Author, Year)"
	return fmt.Sprintf("«%v» (%v, %v)", b.Title, b.Author, b.Year)
}

func main() {
	fmt.Println("=== Реализация Stringer ===")

	// IPAddr
	fmt.Println("\n--- IPAddr ---")
	ips := []IPAddr{
		{127, 0, 0, 1},
		{192, 168, 0, 1},
		{8, 8, 8, 8},
	}
	for _, ip := range ips {
		fmt.Printf("IP: %v\n", ip)
	}

	// Duration
	fmt.Println("\n--- Duration ---")
	durations := []Duration{
		65,    // 1 минута 5 секунд
		3665,  // 1 час 1 минута 5 секунд
		86400, // 24 часа
	}
	for _, d := range durations {
		fmt.Printf("%d секунд = %v\n", int(d), d)
	}

	// Book
	fmt.Println("\n--- Book ---")
	books := []Book{
		{"Война и мир", "Лев Толстой", 1869},
		{"Мастер и Маргарита", "Михаил Булгаков", 1967},
		{"1984", "Джордж Оруэлл", 1949},
	}
	for _, book := range books {
		fmt.Println(book)
	}
}
