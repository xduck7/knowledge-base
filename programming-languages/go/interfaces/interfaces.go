package main

import "fmt"

// ===== ОБЪЯСНЕНИЕ ИНТЕРФЕЙСОВ =====
// Интерфейс описывает набор методов, которые должен реализовать тип.

type Speaker interface {
	Speak() string
}

type Greeter interface {
	Greet() string
}

// ===== РЕАЛИЗАЦИЯ ИНТЕРФЕЙСОВ =====
type Person struct {
	Name string
	Age  int
}

func (p Person) Speak() string {
	return fmt.Sprintf("Hi, I'm %s", p.Name)
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hello! Nice to meet you, I'm %s", p.Name)
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s says: Woof!", d.Name)
}

// ===== ФУНКЦИИ С ПАРАМЕТРАМИ ИНТЕРФЕЙСОВ =====
func MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func MakeGreet(g Greeter) {
	fmt.Println(g.Greet())
}

// ===== ПРАКТИКА EMPTY INTERFACE =====
// Пустой интерфейс может хранить любой тип
func PrintAny(a any) {
	fmt.Printf("Value: %v, Type: %T\n", a, a)
}

func main() {
	// ===== СОЗДАНИЕ ОБЪЕКТОВ =====
	p := Person{Name: "Alice", Age: 25}
	d := Dog{Name: "Rex"}

	// ===== ВЫЗОВ ЧЕРЕЗ ИНТЕРФЕЙС =====
	MakeSpeak(p)
	MakeSpeak(d)

	MakeGreet(p)
	// MakeGreet(d) // ошибка компиляции: Dog не реализует Greeter

	// ===== EMPTY INTERFACE =====
	var v any
	v = 42
	PrintAny(v)
	v = "hello"
	PrintAny(v)
	v = []int{1, 2, 3}
	PrintAny(v)

	// ===== ПОЛИМОРФИЗМ =====
	speakers := []Speaker{p, d}
	for _, s := range speakers {
		fmt.Println(s.Speak())
	}
}
