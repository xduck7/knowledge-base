package main

import (
	"fmt"

	"modules/mathutils" // импорт локального модуля
)

// ===== МОДУЛИ В GO =====
// Модуль (Go module) — это логическая единица кода, объединённая общим go.mod-файлом.
// Каждый модуль имеет версию и имя, обычно совпадающее с путём в репозитории.
// Модуль управляет зависимостями через go.mod и go.sum.

// В этом примере мы используем локальный модуль example.com/mathutils,
// который можно создать в отдельной папке для демонстрации.

// Скачивать модулей из других ресурсов выполняется командой `go get <url>`

func main() {
	fmt.Println("=== Работа с модулями ===")

	a, b := 6, 3
	fmt.Printf("%d + %d = %d\n", a, b, mathutils.Add(a, b))
	fmt.Printf("%d * %d = %d\n", a, b, mathutils.Mul(a, b))

	fmt.Println("Модуль успешно подключен и работает.")

	fmt.Println("Константа из модуля mathutils: ", mathutils.Pi)
	// fmt.Println("Константа из модуля mathutils: ", mathutils.e) // нельзя, так как e приватная (с маленькой буквы)

	num := mathutils.IntNumber{
		Value: 111,
	}

	fmt.Println(num.Value)
	// fmt.Println(num.private.name) // нельзя, так как поля приватные
}
