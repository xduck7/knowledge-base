package main

import "fmt"

type CustomInt int8
type AliasInt int8

// ===== ОБОБЩЕННЫЕ ФУНКЦИИ С ОГРАНИЧЕНИЯМИ ТИПОВ =====

// Number — интерфейс, допускающий разные целые и вещественные типы
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type StrongNumber interface {
	int | int8 | int16 | int32 | int64
}

// Sum суммирует два значения типа T, где T удовлетворяет Number
func Sum[T Number](a, b T) T {
	return a + b
}

// Sum суммирует два значения типа T, где T удовлетворяет StrongNumber
func StrongSum[T StrongNumber](a, b T) T {
	return a + b
}

// PrintSlice печатает срез любого типа
func PrintSlice[T any](s []T) {
	for i, v := range s {
		fmt.Printf("Index %d: %v\n", i, v)
	}
	fmt.Println()
}

// ===== ОБОБЩЕННЫЕ СТРУКТУРЫ =====
type Pair[T any, U any] struct {
	First  T
	Second U
}

func (p Pair[T, U]) Print() {
	fmt.Printf("First: %v, Second: %v\n", p.First, p.Second)
}

func main() {
	// ===== ФУНКЦИИ С ~T =====
	fmt.Println("Sum examples:")
	fmt.Println("Sum[int8]:", Sum[int8](10, 20))
	fmt.Println("Sum[int32]:", Sum[int32](1000, 2000))

	// Типы данных созданные на основе T
	fmt.Println("Sum[CustomInt]:", Sum[CustomInt](CustomInt(1), CustomInt(5)))
	fmt.Println("Sum[AliasInt]:", Sum[AliasInt](AliasInt(1), AliasInt(5)))

	// Функции, где T строгое
	fmt.Println("StrongSum[CustomInt]:", StrongSum[int8](1, 5))
	// error: Cannot use AliasInt as the type StrongNumber
	// fmt.Println("StrongSum[CustomInt]:", StrongSum[CustomInt](CustomInt(1), CustomInt(5)))
	// fmt.Println("StrongSum[AliasInt]:", StrongSum[AliasInt](AliasInt(1), AliasInt(5)))

	// ===== GENERICS С ФУНКЦИЯМИ =====
	ints := []int{1, 2, 3}
	int8s := []int8{10, 20, 30}
	strs := []string{"a", "b", "c"}

	fmt.Println("PrintSlice[int]:")
	PrintSlice(ints)
	fmt.Println("PrintSlice[int8]:")
	PrintSlice(int8s)
	fmt.Println("PrintSlice[string]:")
	PrintSlice(strs)

	// ===== GENERICS С СТРУКТУРАМИ =====
	p1 := Pair[int, string]{First: 10, Second: "ten"}
	p2 := Pair[string, string]{First: "hello", Second: "world"}

	p1.Print()
	p2.Print()
}
