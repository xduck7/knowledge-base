package main

import "fmt"

func main() {
	// В Go есть только один оператор цикла - for

	// классический for: инициализация, условие, инкремент
	for i := 0; i < 5; i++ {
		fmt.Println("Classic for:", i)
	}

	// цикл как while
	j := 0
	for j < 3 {
		fmt.Println("While-style:", j)
		j++
	}

	// бесконечный цикл
	counter := 0
	for {
		if counter >= 3 {
			break
		}
		fmt.Println("Infinite loop, counter:", counter)
		counter++
	}

	// проход по массиву или срезу
	numbers := []int{10, 20, 30}
	for index, value := range numbers {
		fmt.Println("Range slice:", index, value)
	}

	// проход по строке (по рунам)
	str := "Go語"
	for i, r := range str {
		fmt.Printf("Range string: index=%d rune=%c\n", i, r)
	}

	// проход по map
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	for key, val := range m {
		fmt.Println("Range map:", key, val)
	}

	// пример с continue и break
	for i := 0; i < 6; i++ {
		if i == 2 {
			continue // пропустить итерацию
		}
		if i == 5 {
			break // выйти из цикла
		}
		fmt.Println("Control flow:", i)
	}
}
