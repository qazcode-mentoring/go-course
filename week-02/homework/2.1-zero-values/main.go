package main

import "fmt"

func main() {
	// TODO: Объяви переменные всех базовых типов БЕЗ инициализации
	// и выведи их значения
	var str string
	var number int
	var f float64
	var t bool

	fmt.Println(str, number, f, t)

	// Целые числа
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64

	fmt.Println("=== Целые числа (signed) ===")
	fmt.Printf("int: %d\n", i)
	fmt.Printf("int8: %d\n", i8)
	fmt.Printf("int16: %d\n", i16)
	fmt.Printf("int32: %d\n", i32)
	fmt.Printf("int64: %d\n", i64)

	// TODO: Добавь остальные типы:
	// - uint, uint8, uint16, uint32, uint64
	// - float32, float64
	// - bool
	// - string
	var ui uint
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64

	fmt.Printf("uint: %d\n", ui)
	fmt.Printf("uint8: %d\n", ui8)
	fmt.Printf("uint16: %d\n", ui16)
	fmt.Printf("uint32: %d\n", ui32)
	fmt.Printf("uint64: %d\n", ui64)

	var f32 float32
	var f64 float64

	fmt.Printf("float32: %.f\n", f32)
	fmt.Printf("float64: %.f\n", f64)

	var b bool
	var s string

	fmt.Printf("bool: %v\n", b)
	fmt.Printf("string: %s\n", s)

	// TODO: Напиши ответ на вопрос здесь:
	// Почему в Go нет null/nil для базовых типов? Какие проблемы это решает?
	//
	// Ответ: Чтобы избежать лишних проверок на отсутсвие значения в переменной,
	// в Go базовые типы имеют значения по умолчанию.
	// Решает проблему забытых проверок на null/nil
}
