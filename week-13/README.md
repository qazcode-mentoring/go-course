# Неделя 13: Основы конкурентности (Goroutines и Channels)

## Теория

**Ссылка:** https://www.gocat.dev/tour/

Пройди уроки 90-96:
90. Горутины (Goroutines)
91. Каналы (Channels)
92. Буферизованные каналы (Buffered Channels)
93. Range и Close
94. Select
95. Default Selection
96. Упражнение: эквивалентные бинарные деревья

## Что такое конкурентность?

Конкурентность (concurrency) - это способность программы выполнять несколько задач одновременно. Go был спроектирован с поддержкой конкурентности как одной из ключевых особенностей языка.

### Goroutines (Горутины)

Горутина - это легковесный поток выполнения, управляемый Go runtime. Горутины значительно легче системных потоков (занимают ~2KB стека vs ~1MB для потока).

```go
// Запуск горутины
go func() {
    fmt.Println("Я работаю в горутине!")
}()

// Или с именованной функцией
go sayHello("World")
```

### Channels (Каналы)

Каналы - это типизированные "трубы", через которые горутины могут безопасно обмениваться данными.

```go
// Создание канала
ch := make(chan int)

// Отправка значения в канал
ch <- 42

// Получение значения из канала
value := <-ch
```

### Буферизованные каналы

Буферизованные каналы имеют внутреннюю очередь фиксированного размера.

```go
// Канал с буфером на 3 элемента
ch := make(chan int, 3)

ch <- 1  // не блокируется
ch <- 2  // не блокируется
ch <- 3  // не блокируется
ch <- 4  // БЛОКИРУЕТСЯ, буфер полон
```

### Select

Select позволяет горутине ожидать несколько операций с каналами одновременно.

```go
select {
case msg := <-ch1:
    fmt.Println("Получено из ch1:", msg)
case msg := <-ch2:
    fmt.Println("Получено из ch2:", msg)
case <-time.After(time.Second):
    fmt.Println("Таймаут!")
default:
    fmt.Println("Нет данных")
}
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 13.1 | [concurrent-hello](./homework/13.1-concurrent-hello/) | Запуск нескольких горутин, понимание порядка выполнения |
| 13.2 | [channel-ping-pong](./homework/13.2-channel-ping-pong/) | Синхронизация через каналы (ping-pong между горутинами) |
| 13.3 | [channel-generator](./homework/13.3-channel-generator/) | Генератор чисел/данных через канал |
| 13.4 | [buffered-tasks](./homework/13.4-buffered-tasks/) | Буферизованный канал для очереди задач |

## Вопросы для самопроверки

Создай файл `answers-13.txt` и напиши ответы на вопросы:

1. Чем горутина отличается от системного потока (thread)? Почему горутины считаются "легковесными"?
2. Что произойдёт, если отправить значение в небуферизованный канал, когда никто его не читает?
3. В чём разница между `make(chan int)` и `make(chan int, 5)`?
4. Для чего используется `close(ch)` и как проверить, что канал закрыт?
5. Как select помогает работать с несколькими каналами одновременно?

## Дополнительные материалы

- [Go by Example: Goroutines](https://gobyexample.com/goroutines)
- [Go by Example: Channels](https://gobyexample.com/channels)
- [Go by Example: Buffered Channels](https://gobyexample.com/channel-buffering)
- [Go by Example: Select](https://gobyexample.com/select)
- [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go Blog: Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
