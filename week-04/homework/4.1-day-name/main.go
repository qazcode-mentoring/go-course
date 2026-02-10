package main

import "fmt"

// DayName возвращает название дня недели по номеру (1-7)
func DayName(n int) string {
	// TODO: реализуй с помощью switch
	switch n {
	case 1:
		return "Понедельник"
	// TODO: добавь остальные дни
	case 2:
		return "Вторник"
	case 3:
		return "Среда"
	case 4:
		return "Четверг"
	case 5:
		return "Пятница"
	case 6:
		return "Суббота"
	case 7:
		return "Воскресене"
	default:
		return "Некорректный день"
	}
}

func main() {
	for i := 0; i <= 8; i++ {
		fmt.Printf("День %d: %s\n", i, DayName(i))
	}
}
