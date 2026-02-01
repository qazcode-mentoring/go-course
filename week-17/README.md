# Неделя 17: JSON и работа с файлами

## Теория

Работа с JSON и файлами — это базовые навыки, необходимые практически в любом Go-проекте. В этом модуле мы изучим пакет `encoding/json` для сериализации и десериализации данных, а также пакеты `os`, `io` и `bufio` для эффективной работы с файловой системой.

### Пакет encoding/json

Go имеет встроенную поддержку JSON через пакет `encoding/json`. Две основные операции:
- **Marshal** (сериализация) — преобразование Go-структур в JSON
- **Unmarshal** (десериализация) — преобразование JSON в Go-структуры

#### Сериализация (Marshal)

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func main() {
    user := User{
        ID:       1,
        Username: "gopher",
        Email:    "gopher@example.com",
    }

    // Компактный JSON
    data, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
    // {"id":1,"username":"gopher","email":"gopher@example.com"}

    // Форматированный JSON (с отступами)
    prettyData, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(prettyData))
    // {
    //   "id": 1,
    //   "username": "gopher",
    //   "email": "gopher@example.com"
    // }
}
```

#### Десериализация (Unmarshal)

```go
func main() {
    jsonData := `{"id":1,"username":"gopher","email":"gopher@example.com"}`

    var user User
    err := json.Unmarshal([]byte(jsonData), &user)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", user)
    // {ID:1 Username:gopher Email:gopher@example.com}
}
```

### Struct Tags (теги структур)

Теги структур управляют тем, как Go сериализует и десериализует поля:

```go
type User struct {
    // Стандартное имя в JSON
    ID int `json:"id"`

    // Поле будет пропущено, если значение пустое (zero value)
    Email string `json:"email,omitempty"`

    // Поле всегда пропускается при сериализации
    Password string `json:"-"`

    // Имя отличается от имени поля в Go
    CreatedAt time.Time `json:"created_at"`

    // Строковое представление числа
    Balance float64 `json:"balance,string"`
}
```

#### Правила omitempty

`omitempty` пропускает поле, если оно имеет "zero value":
- Числа: `0`
- Строки: `""`
- Булевы: `false`
- Указатели, слайсы, мапы: `nil`

```go
type Config struct {
    Host    string `json:"host,omitempty"`
    Port    int    `json:"port,omitempty"`
    Debug   bool   `json:"debug,omitempty"`
    Tags    []string `json:"tags,omitempty"`
}

func main() {
    cfg := Config{Host: "localhost"} // Port=0, Debug=false, Tags=nil
    data, _ := json.Marshal(cfg)
    fmt.Println(string(data))
    // {"host":"localhost"}
}
```

### Работа с динамическим JSON

Когда структура JSON заранее неизвестна, используй `map[string]any` или `json.RawMessage`:

```go
func main() {
    jsonData := `{"name":"test","value":42,"nested":{"key":"val"}}`

    // Вариант 1: map[string]any
    var data map[string]any
    json.Unmarshal([]byte(jsonData), &data)

    name := data["name"].(string)        // type assertion
    value := data["value"].(float64)     // числа всегда float64
    fmt.Printf("name=%s, value=%v\n", name, value)

    // Вариант 2: json.RawMessage для отложенного парсинга
    type Response struct {
        Type    string          `json:"type"`
        Payload json.RawMessage `json:"payload"`
    }
}
```

### Пакеты для работы с файлами

#### os — базовые операции

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Создание файла
    file, err := os.Create("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Запись в файл
    file.WriteString("Hello, World!\n")

    // Чтение файла целиком
    data, err := os.ReadFile("test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))

    // Запись файла целиком
    err = os.WriteFile("output.txt", []byte("content"), 0644)
    if err != nil {
        panic(err)
    }

    // Проверка существования файла
    if _, err := os.Stat("test.txt"); os.IsNotExist(err) {
        fmt.Println("файл не существует")
    }

    // Удаление файла
    os.Remove("test.txt")
}
```

#### Режимы открытия файлов

```go
// Только чтение
file, err := os.Open("file.txt")

// Создание (перезаписывает если существует)
file, err := os.Create("file.txt")

// Гибкое открытие с флагами
file, err := os.OpenFile("file.txt",
    os.O_RDWR|os.O_CREATE|os.O_APPEND, // флаги
    0644)                                // права доступа

// Флаги:
// os.O_RDONLY — только чтение
// os.O_WRONLY — только запись
// os.O_RDWR   — чтение и запись
// os.O_CREATE — создать если не существует
// os.O_APPEND — добавлять в конец
// os.O_TRUNC  — очистить при открытии
```

#### bufio — буферизованное чтение

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("large_file.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Чтение построчно
    scanner := bufio.NewScanner(file)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        fmt.Printf("%d: %s\n", lineNum, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
}
```

#### bufio.Writer — буферизованная запись

```go
func main() {
    file, err := os.Create("output.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    writer.WriteString("line 1\n")
    writer.WriteString("line 2\n")
    writer.Flush() // обязательно сбросить буфер!
}
```

### io — интерфейсы и утилиты

```go
package main

import (
    "io"
    "os"
    "strings"
)

func main() {
    // Копирование данных
    src, _ := os.Open("source.txt")
    dst, _ := os.Create("dest.txt")
    defer src.Close()
    defer dst.Close()

    io.Copy(dst, src)

    // Чтение из io.Reader в []byte
    reader := strings.NewReader("hello world")
    data, _ := io.ReadAll(reader)
    fmt.Println(string(data))

    // io.MultiWriter — запись в несколько мест одновременно
    multi := io.MultiWriter(os.Stdout, file)
    multi.Write([]byte("goes to both\n"))
}
```

### JSON + файлы: практический пример

```go
package main

import (
    "encoding/json"
    "os"
)

type Config struct {
    Server   string   `json:"server"`
    Port     int      `json:"port"`
    Features []string `json:"features"`
}

// Загрузка конфигурации из файла
func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// Сохранение конфигурации в файл
func SaveConfig(filename string, cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644)
}
```

### json.Encoder и json.Decoder

Для работы с потоками (файлы, HTTP) используй Encoder и Decoder:

```go
func main() {
    // Запись JSON напрямую в файл
    file, _ := os.Create("data.json")
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    encoder.Encode(User{ID: 1, Username: "test"})

    // Чтение JSON напрямую из файла
    file2, _ := os.Open("data.json")
    defer file2.Close()

    var user User
    decoder := json.NewDecoder(file2)
    decoder.Decode(&user)
}
```

### Обработка ошибок JSON

```go
func parseJSON(data []byte) error {
    var result map[string]any
    err := json.Unmarshal(data, &result)

    if err != nil {
        // Проверка типа ошибки
        var syntaxErr *json.SyntaxError
        var typeErr *json.UnmarshalTypeError

        switch {
        case errors.As(err, &syntaxErr):
            return fmt.Errorf("синтаксическая ошибка на позиции %d", syntaxErr.Offset)
        case errors.As(err, &typeErr):
            return fmt.Errorf("неверный тип для поля %s", typeErr.Field)
        default:
            return fmt.Errorf("ошибка парсинга: %w", err)
        }
    }

    return nil
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 17.1 | [json-marshal](./homework/17.1-json-marshal/) | Сериализация структур в JSON, struct tags |
| 17.2 | [json-unmarshal](./homework/17.2-json-unmarshal/) | Парсинг JSON, omitempty, работа с динамическим JSON |
| 17.3 | [file-operations](./homework/17.3-file-operations/) | Чтение и запись файлов (os, bufio) |
| 17.4 | [json-config](./homework/17.4-json-config/) | Загрузка конфигурации приложения из JSON файла |

## Вопросы для самопроверки

Создай файл `answers-17.txt` и напиши ответы на вопросы:

1. В чём разница между `json.Marshal` и `json.NewEncoder().Encode()`? Когда следует использовать каждый из них?

2. Что делает тег `omitempty` и для каких типов он работает? Приведи пример, когда его использование может привести к неожиданному поведению.

3. Почему при десериализации JSON числа всегда становятся `float64`, если использовать `map[string]any`? Как получить целое число?

4. Объясни разницу между `os.ReadFile` и использованием `bufio.Scanner`. Когда следует использовать каждый подход?

5. Зачем нужен вызов `defer file.Close()` и что произойдёт, если его забыть? Что такое file descriptor leak?

## Дополнительные материалы

- [Go Blog: JSON and Go](https://go.dev/blog/json)
- [Go Doc: encoding/json](https://pkg.go.dev/encoding/json)
- [Go Doc: os package](https://pkg.go.dev/os)
- [Go Doc: bufio package](https://pkg.go.dev/bufio)
- [Go by Example: JSON](https://gobyexample.com/json)
- [Go by Example: Reading Files](https://gobyexample.com/reading-files)
- [Go by Example: Writing Files](https://gobyexample.com/writing-files)
- [Effective Go: The blank identifier](https://go.dev/doc/effective_go#blank)
