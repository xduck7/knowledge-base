package main

import "fmt"

// ===== ОПРЕДЕЛЕНИЕ СТРУКТУРЫ =====
// Структура — это набор полей с именами и типами.
// Она объединяет данные логически, как "объект" без классов.
type Person struct {
	Name string
	Age  int
	City string
}

// Вложенные структуры
type Company struct {
	Name    string
	Founder Person
	Size    int
}

// Анонимная структура (без имени типа)
var anonymous = struct {
	ID    int
	Valid bool
}{ID: 1, Valid: true}

// ===== МЕТОДЫ ДЛЯ СТРУКТУР =====
// Методы определяются через "receiver" (приемник)
func (p Person) Greet() {
	fmt.Printf("Hello, my name is %s, I'm %d from %s\n", p.Name, p.Age, p.City)
}

// Метод с указателем (может изменять состояние)
func (p *Person) HaveBirthday() {
	p.Age++
}

// ===== КОНСТРУКТОРНЫЙ ПОДХОД =====
// Go не имеет конструкторов, но обычно создают функцию NewType()
func NewPerson(name string, age int, city string) Person {
	return Person{Name: name, Age: age, City: city}
}

func main() {
	// ===== СОЗДАНИЕ =====
	p1 := Person{Name: "Alice", Age: 25, City: "Paris"} // полное определение
	p2 := Person{"Bob", 30, "Berlin"}                   // позиционное
	p3 := NewPerson("Charlie", 40, "Tokyo")             // через конструктор

	fmt.Println("p1:", p1)
	fmt.Println("p2:", p2)
	fmt.Println("p3:", p3)
	fmt.Println("Anonymous struct:", anonymous)

	// ===== ДОСТУП К ПОЛЯМ =====
	fmt.Println("Name of p1:", p1.Name)

	// ===== ИЗМЕНЕНИЕ ПОЛЕЙ =====
	p1.City = "London"
	fmt.Println("Updated p1:", p1)

	// ===== МЕТОДЫ =====
	p1.Greet()
	p1.HaveBirthday()
	p1.Greet()

	// ===== ВЛОЖЕННЫЕ СТРУКТУРЫ =====
	c := Company{
		Name: "GoSoft",
		Founder: Person{
			Name: "Dana",
			Age:  35,
			City: "NYC",
		},
		Size: 120,
	}

	fmt.Println("Company:", c)
	fmt.Println("Founder Name:", c.Founder.Name)

	// ===== СРЕЗЫ И МАССИВЫ СТРУКТУР =====
	team := []Person{
		{"Eve", 28, "Prague"},
		{"Frank", 32, "Rome"},
	}
	for _, member := range team {
		member.Greet()
	}

	// ===== АНОНИМНЫЕ ПОЛЯ =====
	type Point struct {
		X, Y int
	}
	type Circle struct {
		Point
		Radius int
	}
	circle := Circle{Point{10, 20}, 5}
	fmt.Println("Circle center:", circle.X, circle.Y, "radius:", circle.Radius)

	// ===== СРАВНЕНИЕ =====
	// Структуры сравнимы, если все их поля сравнимы
	pA := Person{"Ann", 22, "Oslo"}
	pB := Person{"Ann", 22, "Oslo"}
	fmt.Println("pA == pB:", pA == pB)
}
