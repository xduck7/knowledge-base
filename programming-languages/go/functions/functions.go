package main

import "fmt"

// ===== ОБЫЧНАЯ ФУНКЦИЯ =====
// Определение функции: ключевое слово `func`, имя, параметры, тип возвращаемого значения.
func add(a int, b int) int {
	return a + b
}

// Параметры одного типа можно объединять
func multiply(a, b int) int {
	return a * b
}

// ===== НЕСКОЛЬКО ВОЗВРАЩАЕМЫХ ЗНАЧЕНИЙ =====
func divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

// ===== ИМЕНОВАННЫЕ РЕЗУЛЬТАТЫ =====
// Можно задать имена возвращаемых значений — они ведут себя как локальные переменные
func rectangle(w, h int) (area int, perimeter int) {
	area = w * h
	perimeter = 2 * (w + h)
	return // возвращает area, perimeter автоматически
}

// ===== ВАРИАДИЧЕСКИЕ ФУНКЦИИ =====
// Принимают произвольное число аргументов (аналог ...args)
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// ===== ФУНКЦИЯ КАК ЗНАЧЕНИЕ =====
// Можно присвоить функцию переменной и вызывать её
func operate(f func(int, int) int, x, y int) int {
	return f(x, y)
}

// ===== АНОНИМНЫЕ ФУНКЦИИ =====
// Создаются прямо в коде без имени
var anon = func(msg string) {
	fmt.Println("Anon says:", msg)
}

// ===== РЕКУРСИЯ =====
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// ===== ЗАМЫКАНИЯ (CLOSURES) =====
// Функции могут замыкать внешние переменные
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	// ===== ВЫЗОВ ПРОСТЫХ ФУНКЦИЙ =====
	fmt.Println("Add:", add(3, 5))
	fmt.Println("Multiply:", multiply(4, 7))

	// ===== НЕСКОЛЬКО РЕЗУЛЬТАТОВ =====
	res, ok := divide(10, 2)
	fmt.Println("Divide:", res, "ok:", ok)

	_, fail := divide(10, 0)
	fmt.Println("Divide by zero ok:", fail)

	// ===== ИМЕНОВАННЫЕ РЕЗУЛЬТАТЫ =====
	a, p := rectangle(5, 3)
	fmt.Println("Area:", a, "Perimeter:", p)

	// ===== ВАРИАДИК =====
	fmt.Println("Sum:", sum(1, 2, 3, 4, 5))

	// ===== ФУНКЦИЯ КАК ПАРАМЕТР =====
	fmt.Println("Operate(add):", operate(add, 10, 20))
	fmt.Println("Operate(multiply):", operate(multiply, 6, 9))

	// ===== АНОНИМНАЯ =====
	anon("Hello, world!")

	func(name string) {
		fmt.Println("Inline anonymous for:", name)
	}("Gopher")

	// ===== РЕКУРСИЯ =====
	fmt.Println("Factorial 5:", factorial(5))

	// ===== ЗАМЫКАНИЕ =====
	next := counter()
	fmt.Println("Counter:", next())
	fmt.Println("Counter:", next())
	fmt.Println("Counter:", next())
}
