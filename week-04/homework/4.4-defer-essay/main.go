package main

import (
	"fmt"
	"io"
	"net/http"
)

/*
TODO: Напиши своё объяснение defer здесь

1. Что такое defer?
команда, которая откладывает выполение функции и не выполняется до выходу окружающей функции

2. В каком порядке выполняются несколько defer?
Несколько defer работают по приниципу последний зашел - первый вышел (LIFO), то есть в обратном порядке


3. Примеры использования defer в реальном коде:
   - Пример 1:
	Одно из основных применений defer это управление ресурсами,
	которые необходимо закрывать или освобождать после использования.
	Такие операции, как закрытие файлов или сетевых соединений, идеально
	подходят для defer, потому что они должны быть выполнены независимо от того,
	как завершится функция.

	file, err := os.Open("data.txt")
	if err != nil {
		return err
	}
	defer file.Close()

   - Пример 2:
	func fetchData(url string) error {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("ошибка получения данных: %v", err)
		}
		defer resp.Body.Close()

		// обработка данных из ответа
		// ...

		return nil
	}
	defer resp.Body.Close() обеспечивает закрытие
	HTTP-соединения после того, как работа с ним будет закончена.

   - Пример 3:
	функция fetchData()

*/

func fetchData(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка получения данных: %v", err)
	}
	defer resp.Body.Close() // закрываем поток при выходе из функции

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	fmt.Println(string(body)) // выводим данные
	return nil
}

func main() {
	// Демонстрация работы defer
	fmt.Println("Начало функции main")

	defer fmt.Println("Это выполнится последним (defer 1)")
	defer fmt.Println("Это выполнится предпоследним (defer 2)")
	defer fmt.Println("Это выполнится третьим с конца (defer 3)")

	fmt.Println("Конец функции main (но до defer)")

	err := fetchData("https://api.github.com")
	if err != nil {
		fmt.Println(err)
	}

}
