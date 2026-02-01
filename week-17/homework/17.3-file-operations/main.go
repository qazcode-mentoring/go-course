package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadFile читает файл целиком и возвращает содержимое как строку
func ReadFile(filename string) (string, error) {
	// TODO: используй os.ReadFile для чтения файла
	// Преобразуй []byte в string перед возвратом
	return "", nil
}

// WriteFile записывает строку в файл (создаёт или перезаписывает)
func WriteFile(filename string, content string) error {
	// TODO: используй os.WriteFile
	// Права доступа: 0644
	return nil
}

// ReadLines читает файл и возвращает слайс строк (без символов новой строки)
func ReadLines(filename string) ([]string, error) {
	// TODO:
	// 1. Открой файл с os.Open
	// 2. Не забудь defer file.Close()
	// 3. Создай bufio.Scanner
	// 4. В цикле scanner.Scan() собирай строки через scanner.Text()
	// 5. Проверь scanner.Err() перед возвратом

	// Подавляем предупреждение о неиспользуемом импорте
	_ = bufio.NewScanner

	return nil, nil
}

// WriteLines записывает слайс строк в файл
func WriteLines(filename string, lines []string) error {
	// TODO:
	// 1. Создай файл с os.Create
	// 2. Не забудь defer file.Close()
	// 3. Создай bufio.Writer для буферизованной записи
	// 4. Запиши каждую строку + "\n"
	// 5. Не забудь writer.Flush() перед выходом!
	return nil
}

// AppendToFile добавляет строку в конец файла
func AppendToFile(filename string, content string) error {
	// TODO:
	// 1. Открой файл с os.OpenFile и флагами:
	//    os.O_APPEND|os.O_WRONLY|os.O_CREATE
	// 2. Права доступа: 0644
	// 3. Запиши content + "\n"
	return nil
}

// CopyFile копирует файл из src в dst
func CopyFile(src, dst string) error {
	// TODO:
	// 1. Открой исходный файл для чтения
	// 2. Создай целевой файл
	// 3. Используй io.Copy(dst, src)
	// 4. Не забудь закрыть оба файла!

	// Подавляем предупреждение об неиспользуемых импортах
	_ = io.Copy
	return nil
}

// FileExists проверяет существование файла
func FileExists(filename string) bool {
	// TODO: используй os.Stat
	// os.IsNotExist(err) вернёт true если файл не существует
	return false
}

// CountLines считает количество строк в файле
func CountLines(filename string) (int, error) {
	// TODO: используй ReadLines или bufio.Scanner напрямую
	return 0, nil
}

// GetFileInfo возвращает информацию о файле
func GetFileInfo(filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("  Ошибка: %v\n", err)
		return
	}

	fmt.Printf("  Имя: %s\n", info.Name())
	fmt.Printf("  Размер: %d байт\n", info.Size())
	fmt.Printf("  Права: %s\n", info.Mode())
	fmt.Printf("  Время изменения: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("  Директория: %v\n", info.IsDir())
}

func main() {
	fmt.Println("=== Работа с файлами ===")

	// Имена тестовых файлов
	testFile := "test_output.txt"
	linesFile := "lines_output.txt"
	copyFile := "copy_output.txt"

	// Очистка в конце (опционально)
	defer func() {
		// Раскомментируй для автоматической очистки:
		// os.Remove(testFile)
		// os.Remove(linesFile)
		// os.Remove(copyFile)
	}()

	// 1. Запись в файл
	fmt.Println("\n--- 1. Запись в файл ---")
	content := "Привет, Go!\nЭто тестовый файл.\nТретья строка."
	err := WriteFile(testFile, content)
	if err != nil {
		fmt.Printf("Ошибка записи: %v\n", err)
	} else {
		fmt.Printf("Файл '%s' создан\n", testFile)
	}

	// 2. Чтение из файла
	fmt.Println("\n--- 2. Чтение из файла ---")
	readContent, err := ReadFile(testFile)
	if err != nil {
		fmt.Printf("Ошибка чтения: %v\n", err)
	} else {
		fmt.Printf("Содержимое:\n%s\n", readContent)
	}

	// 3. Проверка существования
	fmt.Println("\n--- 3. Проверка существования ---")
	fmt.Printf("'%s' существует: %v\n", testFile, FileExists(testFile))
	fmt.Printf("'nonexistent.txt' существует: %v\n", FileExists("nonexistent.txt"))

	// 4. Информация о файле
	fmt.Println("\n--- 4. Информация о файле ---")
	GetFileInfo(testFile)

	// 5. Запись строк
	fmt.Println("\n--- 5. Запись строк ---")
	lines := []string{
		"Первая строка",
		"Вторая строка",
		"Третья строка",
		"", // пустая строка тоже должна записаться
		"Пятая строка",
	}
	err = WriteLines(linesFile, lines)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Записано %d строк в '%s'\n", len(lines), linesFile)
	}

	// 6. Чтение строк
	fmt.Println("\n--- 6. Чтение строк ---")
	readLines, err := ReadLines(linesFile)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		for i, line := range readLines {
			if line == "" {
				fmt.Printf("  %d: (пустая строка)\n", i+1)
			} else {
				fmt.Printf("  %d: %s\n", i+1, line)
			}
		}
	}

	// 7. Подсчёт строк
	fmt.Println("\n--- 7. Подсчёт строк ---")
	count, err := CountLines(linesFile)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Количество строк в '%s': %d\n", linesFile, count)
	}

	// 8. Добавление в конец файла
	fmt.Println("\n--- 8. Добавление в файл ---")
	err = AppendToFile(linesFile, "Добавленная строка 1")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	err = AppendToFile(linesFile, "Добавленная строка 2")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	newCount, _ := CountLines(linesFile)
	fmt.Printf("Строк после добавления: %d\n", newCount)

	// 9. Копирование файла
	fmt.Println("\n--- 9. Копирование файла ---")
	err = CopyFile(linesFile, copyFile)
	if err != nil {
		fmt.Printf("Ошибка копирования: %v\n", err)
	} else {
		fmt.Printf("Файл скопирован: '%s' -> '%s'\n", linesFile, copyFile)

		// Проверяем что копия идентична
		original, _ := ReadFile(linesFile)
		copied, _ := ReadFile(copyFile)
		if original == copied {
			fmt.Println("Проверка: копия идентична оригиналу")
		} else {
			fmt.Println("Проверка: ОШИБКА - файлы различаются!")
		}
	}

	// 10. Обработка ошибок
	fmt.Println("\n--- 10. Обработка ошибок ---")
	_, err = ReadFile("this_file_does_not_exist.txt")
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}

	// Бонус: фильтрация строк
	fmt.Println("\n--- Бонус: поиск строк ---")
	allLines, _ := ReadLines(linesFile)
	fmt.Println("Строки содержащие 'строка':")
	for _, line := range allLines {
		if strings.Contains(line, "строка") {
			fmt.Printf("  - %s\n", line)
		}
	}
}
