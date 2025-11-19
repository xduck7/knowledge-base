package mathutils

const (
	Pi = 3.14159
	e  = 2.71828
)

type IntNumber struct {
	Value int
	private
}

type private struct {
	name string
}

// Add возвращает сумму двух чисел.
func Add(a, b int) int {
	return a + b
}

// Mul возвращает произведение двух чисел.
func Mul(a, b int) int {
	return a * b
}
