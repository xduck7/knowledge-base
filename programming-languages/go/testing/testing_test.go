package main

import (
	"fmt"
	"testing"
)

// ===== ЮНИТ-ТЕСТЫ =====
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; ожидалось %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	result := Multiply(3, 4)
	expected := 12
	if result != expected {
		t.Errorf("Multiply(3, 4) = %d; ожидалось %d", result, expected)
	}
}

// ===== EXAMPLES =====
// Пример использования функции для документации
func ExampleAdd() {
	result := Add(1, 2)
	fmt.Println(result)
	// Output: 3
}

func ExampleMultiply() {
	fmt.Println(Multiply(2, 3))
	// Output: 6
}
