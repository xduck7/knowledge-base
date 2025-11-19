package main

import (
	"fmt"
)

func main() {
	// ===== МАССИВЫ =====
	// Массив в Go имеет фиксированную длину, которая является частью его типа.
	// Объявление массива на 3 элемента типа int:
	var arr [3]int
	arr[0] = 10
	arr[1] = 20
	arr[2] = 30
	fmt.Println("Array:", arr)

	// инициализация массива при объявлении
	arr2 := [4]string{"a", "b", "c", "d"}
	fmt.Println("Array literal:", arr2)

	// автоматическое определение длины массива через "..."
	arr3 := [...]int{1, 2, 3, 4, 5}
	fmt.Println("Array auto length:", arr3, "len:", len(arr3))

	// массив - это значение. При передаче в функцию он копируется.
	copyArr := arr3
	copyArr[0] = 999
	fmt.Println("Original array:", arr3)
	fmt.Println("Copied array:", copyArr)

	// ===== СРЕЗЫ =====
	// Срез - это "динамический" массив
	// Он хранит ссылку на массив, длину и ёмкость (capacity).

	// type slice struct {
	// 	array unsafe.Pointer
	// 	len int
	// 	cap int
	// }

	slice := []int{10, 20, 30}
	fmt.Println("Slice:", slice, "len:", len(slice), "cap:", cap(slice))

	// создание среза из массива
	part := arr3[1:4]
	fmt.Println("Sub-slice:", part)

	// изменение среза меняет исходный массив
	part[0] = 777
	fmt.Println("Modified sub-slice:", part)
	fmt.Println("Underlying array changed:", arr3)

	// создание среза через make()
	s := make([]string, 3) // длина 3, ёмкость 3
	s[0] = "x"
	s[1] = "y"
	s[2] = "z"
	fmt.Println("Make slice:", s, "len:", len(s), "cap:", cap(s))

	// добавление элементов через append()
	s = append(s, "new")
	fmt.Println("After append:", s, "len:", len(s), "cap:", cap(s))

	// слияние срезов
	more := []string{"a", "b"}
	s = append(s, more...)
	fmt.Println("After merge:", s)

	// копирование срезов
	dst := make([]int, len(slice))
	copy(dst, slice)
	fmt.Println("Copied slice:", dst)

	// многомерный срез
	matrix := [][]int{
		{1, 2},
		{3, 4, 5},
		{6},
	}
	fmt.Println("Matrix:", matrix)
}
