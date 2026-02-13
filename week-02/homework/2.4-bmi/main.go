package main

import "fmt"

// CalculateBMI вычисляет индекс массы тела
// weightKg - вес в килограммах
// heightCm - рост в сантиметрах
func CalculateBMI(weightKg, heightCm float64) float64 {
	// TODO: реализуй функцию
	// Не забудь перевести сантиметры в метры!
	// Формула: ИМТ = вес (кг) / рост (м)²
	var heightM float64 = heightCm / 100
	var indexBody float64 = weightKg / (heightM * heightM)
	return indexBody
}

// InterpretBMI возвращает категорию по значению ИМТ
func InterpretBMI(bmi float64) string {
	// TODO: реализуй функцию
	// < 18.5: "Недостаточный вес"
	// 18.5 - 24.9: "Норма"
	// 25.0 - 29.9: "Избыточный вес"
	// >= 30: "Ожирение"
	if bmi < 18.5 {
		return "Недостаточный вес"
	} else if bmi <= 24.9 {
		return "Норма"
	} else if bmi <= 29.9 {
		return "Избыточный вес"
	}
	return "Ожирение"
}

func main() {
	// Тест 1: вес 70 кг, рост 175 см → ИМТ ≈ 22.9 (Норма)
	bmi := CalculateBMI(70, 175)
	fmt.Printf("Вес: 70 кг, Рост: 175 см\n")
	fmt.Printf("ИМТ: %.1f\n", bmi)
	fmt.Printf("Категория: %s\n\n", InterpretBMI(bmi))

	// TODO: Добавь ещё несколько тестов:
	// - Недостаточный вес
	bmi1 := CalculateBMI(50, 165)
	fmt.Printf("Вес: 50, Рост: 165 см\n")
	fmt.Printf("ИМТ: %.1f\n", bmi1)
	fmt.Printf("Категория: %s\n\n", InterpretBMI(bmi1))

	// - Избыточный вес
	bmi2 := CalculateBMI(70, 165)
	fmt.Printf("Вес: 70, Рост: 165 см\n")
	fmt.Printf("ИМТ: %.1f\n", bmi2)
	fmt.Printf("Категория: %s\n\n", InterpretBMI(bmi2))
	// - Ожирение
	bmi3 := CalculateBMI(90, 165)
	fmt.Printf("Вес: 90, Рост: 165 см\n")
	fmt.Printf("ИМТ: %.1f\n", bmi3)
	fmt.Printf("Категория: %s\n\n", InterpretBMI(bmi3))
}
