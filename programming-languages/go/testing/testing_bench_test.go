package main

import "testing"

// ===== БЕНЧМАРКИ =====
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(10, 20)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(10, 20)
	}
}
