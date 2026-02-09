package main

import "fmt"

// Ball представляет мяч в игре
type Ball struct {
	Hits int // Количество ударов
}

// player представляет игрока
// Получает мяч из inbox, увеличивает счётчик ударов,
// выводит информацию и отправляет в outbox
// maxHits - когда достигнуто это число, игрок закрывает outbox и завершается
func player(name string, inbox <-chan Ball, outbox chan<- Ball, maxHits int, done chan string) {
	// TODO: реализуй функцию
	// 1. Используй for-range для чтения из inbox
	// 2. Увеличь ball.Hits
	// 3. Выведи: "NAME ударил мяч! (удар #N)"
	// 4. Если ball.Hits >= maxHits:
	//    - Закрой канал outbox
	//    - Заверши функцию (return)
	// 5. Иначе отправь ball в outbox
	for ball := range inbox {
		ball.Hits++
		fmt.Printf("%s ударил мяч! (удар %d)\n", name, ball.Hits)

		if ball.Hits >= maxHits {
			close(outbox)
			done <- "end"
			return
		}

		outbox <- ball
	}
}

// playGame запускает игру ping-pong
// maxHits - максимальное количество ударов
func playGame(maxHits int) {
	// TODO: реализуй функцию
	// 1. Создай два канала: pingCh и pongCh
	// 2. Запусти две горутины:
	//    - player("Ping", pingCh, pongCh, maxHits)
	//    - player("Pong", pongCh, pingCh, maxHits)
	// 3. Отправь начальный мяч в pingCh
	// 4. Дождись завершения игры
	//    Подсказка: можно использовать дополнительный канал done
	//    или просто подождать, пока оба канала не закроются
	pingCh := make(chan Ball)
	pongCh := make(chan Ball)
	done := make(chan string)

	go player("Ping", pingCh, pongCh, maxHits, done)
	go player("Pong", pongCh, pingCh, maxHits, done)

	var ball Ball

	pingCh <- ball

	fmt.Println(<-done)
}

func main() {
	fmt.Println("=== Ping-Pong Game ===")

	playGame(6)

	fmt.Println("\n=== Игра на 10 ударов ===")
	playGame(10)
}
