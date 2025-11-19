package main

import (
	"fmt"
	"strings"
)

func main() {
	// ===== ОСНОВЫ СТРОК =====
	// Строки в Go — неизменяемые последовательности байтов (обычно UTF-8).

	// type _string struct {
	// 	element *byte  // указатель на данные
	// 	len     int    // длина в байтах
	// }

	// Создание строки:
	str := "Hello, Go!"
	fmt.Println("String:", str)

	// длина строки в байтах (а не символах):
	fmt.Println("Length (bytes):", len(str))

	// индексация: получаем байт, не символ.
	fmt.Println("First byte value:", str[0])

	// преобразование байтов в символы (руны)
	fmt.Println("As rune:", string(str[0]))

	// ===== РАБОТА С РУНАМИ =====
	// Rune — это Unicode-символ (int32)
	runes := []rune("Привет")
	fmt.Println("Runes slice:", runes)
	fmt.Println("Length (runes):", len(runes))
	fmt.Println("Third rune:", string(runes[2]))

	// ===== КОНКАТЕНАЦИЯ =====
	a := "Go"
	b := "Lang"
	combined := a + " " + b
	fmt.Println("Concatenation:", combined)

	// ===== ФУНКЦИИ ИЗ ПАКЕТА STRINGS =====
	sample := "  apple,banana,orange  "

	fmt.Println("To upper:", strings.ToUpper(sample))
	fmt.Println("To lower:", strings.ToLower(sample))
	fmt.Println("Trim spaces:", strings.TrimSpace(sample))
	fmt.Println("Has prefix '  ap':", strings.HasPrefix(sample, "  ap"))
	fmt.Println("Has suffix 'ge  ':", strings.HasSuffix(sample, "ge  "))
	fmt.Println("Contains 'banana':", strings.Contains(sample, "banana"))
	fmt.Println("Count of 'a':", strings.Count(sample, "a"))

	// разделение строки по разделителю
	parts := strings.Split(sample, ",")
	fmt.Println("Split:", parts)

	// объединение обратно
	joined := strings.Join(parts, "|")
	fmt.Println("Join:", joined)

	// замена подстроки
	replaced := strings.ReplaceAll(sample, "banana", "lemon")
	fmt.Println("Replace:", replaced)

	// ===== ИТЕРАЦИЯ ПО СТРОКЕ =====
	text := "Go語"
	for i, r := range text {
		fmt.Printf("Index %d → rune '%c'\n", i, r)
	}

	// ===== ПРЕОБРАЗОВАНИЯ =====
	bytes := []byte("data")       // строка -> байты
	backToString := string(bytes) // байты -> строка
	fmt.Println("Bytes:", bytes)
	fmt.Println("Back to string:", backToString)
}
