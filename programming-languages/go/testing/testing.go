package main

import "fmt"

// ===== TESTING В GO =====
// Пакет testing используется для написания юнит-тестов.
// Тесты обычно размещаются в файлах с суффиксом `_test.go` и функции начинаются с Test*

func main() {
	fmt.Println("=== Тестирование в Go ===")
	fmt.Println("Функции для тестирования обычно пишутся отдельно и проверяются через `go test`")
}

// ===== ПРИМЕР ФУНКЦИИ ДЛЯ ТЕСТИРОВАНИЯ =====
func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}
