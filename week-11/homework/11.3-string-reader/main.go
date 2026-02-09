package main

import (
	"fmt"
	"io"
	"strings"
)

// RepeatReader бесконечно повторяет один байт
type RepeatReader struct {
	char byte
}

// NewRepeatReader создаёт RepeatReader для указанного символа
func NewRepeatReader(char byte) *RepeatReader {
	// TODO: реализуй функцию
	return &RepeatReader{char: char}
}

// Read заполняет буфер повторяющимся символом
// Никогда не возвращает io.EOF (бесконечный поток)
func (r *RepeatReader) Read(p []byte) (n int, err error) {
	// TODO: реализуй метод
	// Заполни весь буфер p символом r.char
	// Верни len(p), nil
	for i := 0; i < len(p); i++ {
		p[i] = r.char
	}
	return len(p), nil
}

// LimitReader ограничивает количество байт из другого Reader
type LimitReader struct {
	reader io.Reader
	limit  int
	read   int
}

// NewLimitReader создаёт LimitReader
func NewLimitReader(r io.Reader, limit int) *LimitReader {
	// TODO: реализуй функцию
	return &LimitReader{reader: r, limit: limit}
}

// Read читает данные, но не более limit байт всего
func (r *LimitReader) Read(p []byte) (n int, err error) {
	// TODO: реализуй метод
	// Если уже прочитано >= limit, верни 0, io.EOF
	// Иначе читай, но не более чем осталось до limit
	if r.read >= r.limit {
		return 0, io.EOF
	}

	toRead := min(len(p), r.limit-r.read)
	n, err = r.reader.Read(p[:toRead])
	r.read += n

	return n, err
}

// CountingReader подсчитывает количество прочитанных байт
type CountingReader struct {
	reader    io.Reader
	BytesRead int
}

// NewCountingReader создаёт CountingReader
func NewCountingReader(r io.Reader) *CountingReader {
	// TODO: реализуй функцию
	return &CountingReader{reader: r}
}

// Read читает данные и обновляет счётчик
func (r *CountingReader) Read(p []byte) (n int, err error) {
	// TODO: реализуй метод
	// Прочитай из r.reader
	// Добавь n к BytesRead
	// Верни результат

	n, err = r.reader.Read(p)
	r.BytesRead += n
	return n, err
}

func main() {
	fmt.Println("=== Реализация io.Reader ===")

	// RepeatReader
	fmt.Println("\n--- RepeatReader ---")
	rr := NewRepeatReader('A')
	buf := make([]byte, 10)
	n, _ := rr.Read(buf)
	fmt.Printf("Прочитано %d байт: %s\n", n, string(buf))

	// LimitReader
	fmt.Println("\n--- LimitReader ---")
	original := strings.NewReader("Hello, World!")
	lr := NewLimitReader(original, 5)
	data, err := io.ReadAll(lr)
	fmt.Printf("Прочитано: '%s' (err: %v)\n", string(data), err)

	// Ещё пример LimitReader
	original2 := strings.NewReader("Short")
	lr2 := NewLimitReader(original2, 100)
	data2, _ := io.ReadAll(lr2)
	fmt.Printf("Прочитано (limit > len): '%s'\n", string(data2))

	// CountingReader
	fmt.Println("\n--- CountingReader ---")
	cr := NewCountingReader(strings.NewReader("Hello, World!"))
	io.ReadAll(cr)
	fmt.Printf("Всего прочитано байт: %d\n", cr.BytesRead)
}
