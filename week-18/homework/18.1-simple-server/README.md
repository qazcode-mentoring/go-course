# ДЗ 18.1: Simple HTTP Server

## Цель

Научиться создавать простой HTTP-сервер на Go, понять основы пакета `net/http`, интерфейс `http.Handler` и функцию `http.HandleFunc`.

## Что нужно сделать

1. Создать HTTP-сервер, который слушает порт 8080
2. Реализовать обработчик для корневого пути `/`, который возвращает "Hello, World!"
3. Реализовать обработчик `/health`, который возвращает статус сервера
4. Реализовать обработчик `/info`, который показывает информацию о запросе

## Сигнатуры

```go
// Обработчик для корневого пути
func helloHandler(w http.ResponseWriter, r *http.Request)

// Обработчик для проверки здоровья сервера
func healthHandler(w http.ResponseWriter, r *http.Request)

// Обработчик, показывающий информацию о запросе
func infoHandler(w http.ResponseWriter, r *http.Request)
```

## Требования к обработчикам

### GET /
- Возвращает текст "Hello, World!"
- Content-Type: text/plain
- Статус код: 200

### GET /health
- Возвращает текст "OK"
- Content-Type: text/plain
- Статус код: 200

### GET /info
- Возвращает информацию о запросе:
  - Метод запроса
  - Путь (URL Path)
  - User-Agent заголовок
  - Remote Address (IP клиента)
- Content-Type: text/plain
- Формат вывода — по одной строке на каждое поле

## Пример использования

Запуск сервера:
```bash
go run main.go
# Server starting on :8080
```

Проверка с помощью curl:
```bash
$ curl http://localhost:8080/
Hello, World!

$ curl http://localhost:8080/health
OK

$ curl http://localhost:8080/info
Method: GET
Path: /info
User-Agent: curl/8.1.2
Remote Address: 127.0.0.1:52341
```

## Подсказки

- Используй `fmt.Fprintf(w, "текст")` для записи ответа
- `r.Method` — метод запроса (GET, POST, ...)
- `r.URL.Path` — путь запроса
- `r.Header.Get("User-Agent")` — получение заголовка
- `r.RemoteAddr` — адрес клиента
- Для запуска сервера: `http.ListenAndServe(":8080", nil)`
- При использовании `nil` в качестве handler используется `http.DefaultServeMux`

## Как тестировать

1. Запусти сервер: `go run main.go`
2. В другом терминале выполни curl-команды из примеров
3. Или открой в браузере: http://localhost:8080

## Критерии выполнения

- [ ] Сервер запускается на порту 8080
- [ ] GET / возвращает "Hello, World!"
- [ ] GET /health возвращает "OK"
- [ ] GET /info возвращает информацию о запросе (метод, путь, User-Agent, RemoteAddr)
- [ ] Код компилируется без ошибок
- [ ] Логируется информация о запуске сервера
