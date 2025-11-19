package main

import "fmt"

func main() {
	// целые числа (i8, i16, i32, i64)
	var i8 int8 = 127
	var u16 uint16 = 65000
	var i int = -100
	var u uint = 100

	// числа с плавающей точкой (float32, float64)
	var f32 float32 = 3.14
	var f64 float64 = 2.7182818284

	// комплексные числа
	var complex64 complex64 = 1 + 2i
	var complex128 complex128 = 1 + 2i

	// булевы значения (bool)
	var truth bool = true
	var lie bool = false

	// строки (string)
	var s1 string = "hello"

	// символы (rune - alias для int32)
	var r rune = 'Ж' // rune — это int32, хранит Unicode-символ

	// байты (byte - alias для uint8)
	var b byte = 'A'

	// составные типы
	arr := [3]int{1, 2, 3}               // массив фиксированной длины
	slc := []string{"go", "is", "fast"}  // срез
	mp := map[string]int{"a": 1, "b": 2} // map
	ptr := &i                            // указатель
	st := struct {
		Name string
		Age  int
	}{"Alice", 30} // структура

	fmt.Println("Integers:", i8, u16, i, u)
	fmt.Println("Floats:", f32, f64)
	fmt.Println("Complex:", complex64, complex128)
	fmt.Println("Booleans:", truth, lie)
	fmt.Println("Strings:", s1)
	fmt.Println("Rune:", r, string(r))
	fmt.Println("Byte:", b, string(b))
	fmt.Println("Array:", arr)
	fmt.Println("Slice:", slc)
	fmt.Println("Map:", mp)
	fmt.Println("Pointer:", ptr, *ptr)
	fmt.Println("Struct:", st)

	// интерфейсы и пустой интерфейс (any)
	var anyValue any = "любой тип"
	fmt.Println("Any:", anyValue)

	// функции как значения
	fn := func(x, y int) int { return x + y }
	fmt.Println("Func result:", fn(2, 3))
}
