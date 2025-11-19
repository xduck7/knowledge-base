package main

import "fmt"

// ===== УСЛОВНЫЕ КОНСТРУКЦИИ В GO =====
// if, else if, else, switch, и оператор с короткой инициализацией

func main() {
	// ===== ПРОСТОЙ IF =====
	x := 10
	if x > 5 {
		fmt.Println("x больше 5")
	}

	// ===== IF / ELSE =====
	y := -3
	if y > 0 {
		fmt.Println("y положительное")
	} else {
		fmt.Println("y отрицательное или ноль")
	}

	// ===== IF С КОРОТКОЙ ИНИЦИАЛИЗАЦИЕЙ =====
	// объявление переменной только внутри области видимости if
	if n := 7; n%2 == 0 {
		fmt.Println(n, "чётное")
	} else {
		fmt.Println(n, "нечётное")
	}

	// fmt.Println(n) // ошибка: n не видно вне if

	// ===== ВЛОЖЕННЫЕ IF =====
	a := 15
	if a > 0 {
		if a < 20 {
			fmt.Println("a положительное и меньше 20")
		}
	}

	// ===== SWITCH =====
	color := "green"
	switch color {
	case "red":
		fmt.Println("Цвет — красный")
	case "green":
		fmt.Println("Цвет — зелёный")
	case "blue":
		fmt.Println("Цвет — синий")
	default:
		fmt.Println("Неизвестный цвет")
	}

	// ===== SWITCH БЕЗ УСЛОВИЯ =====
	// можно использовать switch как цепочку if-else
	age := 25
	switch {
	case age < 18:
		fmt.Println("Несовершеннолетний")
	case age < 65:
		fmt.Println("Взрослый")
	default:
		fmt.Println("Пожилой")
	}

	// ===== SWITCH С ПАДЕНИЕМ ЧЕРЕЗ FALLTHROUGH =====
	num := 3
	switch num {
	case 1:
		fmt.Println("Один")
	case 2:
		fmt.Println("Два")
	case 3:
		fmt.Println("Три")
		fallthrough // выполняет следующий блок
	case 4:
		fmt.Println("Следом идёт четыре (через fallthrough)")
	default:
		fmt.Println("Что-то другое")
	}

	// ===== SWITCH ПО ТИПУ =====
	var i interface{} = 3.14
	switch v := i.(type) {
	case int:
		fmt.Println("Целое число:", v)
	case float64:
		fmt.Println("Число с плавающей точкой:", v)
	case string:
		fmt.Println("Строка:", v)
	default:
		fmt.Println("Неизвестный тип")
	}
}
