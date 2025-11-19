package main

import "fmt"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("поймана паника:", r)
		}
	}()

	// ===== СОЗДАНИЕ MAP =====
	// Map — это хеш-таблица: ключ → значение.
	// Ключи должны быть сравнимыми (int, string, bool, comparable structs).

	// Объявление через литерал
	m := map[string]int{
		"apple":  5,
		"banana": 10,
		"orange": 7,
	}
	fmt.Println("Map literal:", m)

	// Создание пустой карты через make()
	ages := make(map[string]int)
	ages["Alice"] = 25
	ages["Bob"] = 30
	fmt.Println("Make map:", ages)

	// ===== ДОСТУП К ЭЛЕМЕНТАМ =====
	fmt.Println("Alice age:", ages["Alice"])

	// Если ключа нет — возвращается нулевое значение типа
	fmt.Println("Unknown key:", ages["Charlie"])

	// Проверка существования ключа
	val, exists := ages["Charlie"]
	if !exists {
		fmt.Println("Charlie not found")
	} else {
		fmt.Println("Charlie:", val)
	}

	// ===== УДАЛЕНИЕ =====
	delete(ages, "Bob")
	fmt.Println("After delete:", ages)

	// ===== ИТЕРАЦИЯ =====
	for key, value := range m {
		fmt.Printf("%s => %d\n", key, value)
	}

	// ===== ИЗМЕНЕНИЕ ЗНАЧЕНИЯ =====
	m["apple"] = 100
	fmt.Println("Updated map:", m)

	// ===== ВСТАВКА НОВОГО ЭЛЕМЕНТА =====
	m["melon"] = 15
	fmt.Println("Extended map:", m)

	// ===== ДЛИНА КАРТЫ =====
	fmt.Println("Length:", len(m))

	// ===== ВЛОЖЕННЫЕ КАРТЫ =====
	nested := map[string]map[string]int{
		"user1": {"score": 100, "level": 3},
		"user2": {"score": 80, "level": 2},
	}
	fmt.Println("Nested map:", nested)
	fmt.Println("user1 score:", nested["user1"]["score"])

	// ===== ОБНУЛЕНИЕ КАРТЫ =====
	m = map[string]int{}
	fmt.Println("Cleared map:", m)

	// ===== NIL MAP =====

	var m2 map[string]int // nil map
	fmt.Println("m == nil:", m2 == nil)

	// это безопасно (чтение из nil map):
	fmt.Println("m[\"key\"] =", m2["key"])

	// это вызовет панику, потому что вставка в nil map недопустима
	m2["key"] = 42
	fmt.Println("Done")
}
