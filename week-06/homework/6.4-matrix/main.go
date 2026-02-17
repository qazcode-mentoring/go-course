package main

import "fmt"

// CreateMatrix создаёт матрицу размером rows×cols, заполненную значением value
func CreateMatrix(rows, cols, value int) [][]int {
	// TODO: реализуй функцию
	// Создай слайс слайсов нужного размера
	// Заполни каждый элемент значением value
	matrix := make([][]int, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = value
		}
	}

	return matrix
}

// PrintMatrix выводит матрицу в читаемом формате
func PrintMatrix(matrix [][]int) {
	// TODO: реализуй функцию
	// Выведи каждую строку на отдельной линии
	// Элементы разделяй пробелами или табуляцией
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Print(matrix[i][j], " ")
		}
		fmt.Println()
	}
}

// Transpose возвращает транспонированную матрицу
func Transpose(matrix [][]int) [][]int {
	// TODO: реализуй функцию
	// Если исходная матрица M×N, результат будет N×M
	// Элемент [i][j] становится [j][i]
	rows := len(matrix)
	cols := len(matrix[0])

	trmatrix := make([][]int, cols)

	// создаём внешний слайс с длиной cols (N*M)
	for i := 0; i < cols; i++ {
		trmatrix[i] = make([]int, rows)
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			trmatrix[j][i] = matrix[i][j]
		}
	}

	return trmatrix
}

// SumMatrix возвращает сумму всех элементов матрицы
func SumMatrix(matrix [][]int) int {
	// TODO: реализуй функцию
	// Пройди по всем строкам и столбцам
	var sum int
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

func main() {
	fmt.Println("=== Работа с матрицами ===")

	// Создание матрицы
	matrix := CreateMatrix(2, 3, 1)
	fmt.Println("Созданная матрица 2×3:")
	PrintMatrix(matrix)

	// Работа с готовой матрицей
	matrix = [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("\nИсходная матрица:")
	PrintMatrix(matrix)

	fmt.Println("\nТранспонированная матрица:")
	transposed := Transpose(matrix)
	PrintMatrix(transposed)

	fmt.Println("\nСумма элементов:", SumMatrix(matrix))

	//fmt.Println(CreateMatrix(3, 3, 3))
}
