# ДЗ 17.3: Работа с файлами

## Цель

Научиться читать и записывать файлы в Go, использовать буферизованное чтение для больших файлов и корректно обрабатывать ошибки файловой системы.

## Что нужно сделать

1. Прочитать содержимое файла целиком
2. Записать данные в файл
3. Прочитать файл построчно с помощью `bufio.Scanner`
4. Добавить данные в конец существующего файла (append)
5. Скопировать файл

## Функции для реализации

```go
// ReadFile читает файл целиком и возвращает содержимое как строку
func ReadFile(filename string) (string, error)

// WriteFile записывает строку в файл (создаёт или перезаписывает)
func WriteFile(filename string, content string) error

// ReadLines читает файл и возвращает слайс строк
func ReadLines(filename string) ([]string, error)

// WriteLines записывает слайс строк в файл (каждая строка с новой строки)
func WriteLines(filename string, lines []string) error

// AppendToFile добавляет строку в конец файла
func AppendToFile(filename string, content string) error

// CopyFile копирует файл из src в dst
func CopyFile(src, dst string) error

// FileExists проверяет существование файла
func FileExists(filename string) bool

// CountLines считает количество строк в файле
func CountLines(filename string) (int, error)
```

## Пример использования

```go
func main() {
    // Запись в файл
    err := WriteFile("test.txt", "Hello, World!")
    if err != nil {
        panic(err)
    }

    // Чтение из файла
    content, err := ReadFile("test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(content) // Hello, World!

    // Запись нескольких строк
    lines := []string{"Line 1", "Line 2", "Line 3"}
    WriteLines("lines.txt", lines)

    // Чтение построчно
    readLines, _ := ReadLines("lines.txt")
    for i, line := range readLines {
        fmt.Printf("%d: %s\n", i+1, line)
    }

    // Добавление в конец
    AppendToFile("lines.txt", "Line 4")

    // Копирование
    CopyFile("lines.txt", "lines_backup.txt")
}
```

## Подсказки

### Чтение файла целиком
```go
data, err := os.ReadFile(filename)
// data — это []byte, преобразуй в string при необходимости
```

### Запись файла целиком
```go
err := os.WriteFile(filename, []byte(content), 0644)
// 0644 — права доступа (owner: rw, group: r, others: r)
```

### Чтение построчно
```go
file, err := os.Open(filename)
if err != nil {
    return nil, err
}
defer file.Close()

var lines []string
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    lines = append(lines, scanner.Text())
}
return lines, scanner.Err()
```

### Добавление в конец файла
```go
file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
if err != nil {
    return err
}
defer file.Close()
_, err = file.WriteString(content + "\n")
return err
```

### Проверка существования файла
```go
_, err := os.Stat(filename)
return !os.IsNotExist(err)
```

### Копирование файла
```go
// Способ 1: через io.Copy
src, _ := os.Open(srcFile)
dst, _ := os.Create(dstFile)
defer src.Close()
defer dst.Close()
io.Copy(dst, src)

// Способ 2: через ReadFile/WriteFile (для небольших файлов)
data, _ := os.ReadFile(srcFile)
os.WriteFile(dstFile, data, 0644)
```

## Важно: defer и закрытие файлов

**Всегда** закрывай файлы после использования! Используй `defer`:

```go
file, err := os.Open(filename)
if err != nil {
    return err
}
defer file.Close() // закроется при выходе из функции
```

Если не закрывать файлы:
- Утечка file descriptors
- Возможные ошибки при повторном открытии
- В худшем случае — падение программы

## Права доступа (permissions)

```
0644 — owner: read+write, group: read, others: read
0755 — owner: all, group: read+execute, others: read+execute
0600 — owner: read+write, никто другой
```

## Критерии выполнения

- [ ] `ReadFile` читает файл и возвращает содержимое
- [ ] `WriteFile` создаёт/перезаписывает файл
- [ ] `ReadLines` возвращает слайс строк без символов новой строки
- [ ] `WriteLines` записывает строки с переносами
- [ ] `AppendToFile` добавляет в конец, не затирая существующее содержимое
- [ ] `CopyFile` создаёт точную копию файла
- [ ] `FileExists` корректно проверяет существование
- [ ] `CountLines` считает строки (включая пустые)
- [ ] Все функции корректно обрабатывают ошибки
- [ ] Все открытые файлы закрываются через defer
- [ ] Код компилируется без ошибок
