package main

import (
	"fmt"
	"io"
	"strings"
)

// rot13 применяет ROT13 шифр к одному байту
// A-Z сдвигаются на 13 позиций, a-z сдвигаются на 13 позиций
// Остальные символы не меняются
func rot13(b byte) byte {
	// TODO: реализуй функцию
	// Если b в диапазоне 'A'-'Z':
	//   return 'A' + (b - 'A' + 13) % 26
	// Если b в диапазоне 'a'-'z':
	//   return 'a' + (b - 'a' + 13) % 26
	// Иначе return b
	if b >= 'A' && b <= 'Z' {
		return 'A' + (b-'A'+13)%26
	}
	if b >= 'a' && b <= 'z' {
		return 'a' + (b-'a'+13)%26
	}
	return b
}

// Rot13Reader — Reader-обёртка, применяющая ROT13
type Rot13Reader struct {
	reader io.Reader
}

// NewRot13Reader создаёт новый Rot13Reader
func NewRot13Reader(r io.Reader) *Rot13Reader {
	// TODO: реализуй функцию
	return &Rot13Reader{reader: r}
}

// Read читает данные и применяет ROT13 к каждому байту
func (r *Rot13Reader) Read(p []byte) (n int, err error) {
	// TODO: реализуй метод
	// 1. Прочитай данные из r.reader в p
	// 2. Примени rot13 к каждому прочитанному байту
	// 3. Верни количество прочитанных байт и ошибку
	n, err = r.reader.Read(p)

	for i := 0; i < n; i++ {
		p[i] = rot13(p[i])
	}

	return n, err
}

func main() {
	fmt.Println("=== ROT13 Reader ===")

	// Тест функции rot13
	fmt.Println("\n--- Тест rot13() ---")
	testChars := []byte{'A', 'N', 'Z', 'a', 'n', 'z', '1', '!'}
	for _, c := range testChars {
		fmt.Printf("rot13('%c') = '%c'\n", c, rot13(c))
	}

	// Кодирование
	fmt.Println("\n--- Кодирование ---")
	original := "Hello, World!"
	r := NewRot13Reader(strings.NewReader(original))
	encoded, _ := io.ReadAll(r)
	fmt.Printf("Оригинал: %s\n", original)
	fmt.Printf("ROT13:    %s\n", string(encoded))

	// Декодирование (ROT13 симметричен)
	fmt.Println("\n--- Декодирование ---")
	r2 := NewRot13Reader(strings.NewReader(string(encoded)))
	decoded, _ := io.ReadAll(r2)
	fmt.Printf("Закодировано: %s\n", string(encoded))
	fmt.Printf("Декодировано: %s\n", string(decoded))

	// Проверка с секретным сообщением
	fmt.Println("\n--- Секретное сообщение ---")
	secret := "Tb yrnea Tb!"
	r3 := NewRot13Reader(strings.NewReader(secret))
	message, _ := io.ReadAll(r3)
	fmt.Printf("Секрет: %s\n", secret)
	fmt.Printf("Расшифровано: %s\n", string(message))
}
