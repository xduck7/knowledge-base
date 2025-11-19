package main

import (
	"errors"
	"fmt"
)

// ===== ОБРАБОТКА ОШИБОК В GO =====
// Ошибки в Go — это значения типа error, которые возвращаются как часть результата функции.

// ===== СОЗДАНИЕ ПРОСТОЙ ОШИБКИ =====
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("деление на ноль")
	}
	return a / b, nil
}

// ===== СОЗДАНИЕ ПОЛЬЗОВАТЕЛЬСКОГО ТИПА ОШИБКИ =====
type MyError struct {
	Code int
	Msg  string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Ошибка [%d]: %s", e.Code, e.Msg)
}

func riskyOperation(ok bool) error {
	if !ok {
		return &MyError{Code: 404, Msg: "ресурс не найден"}
	}
	return nil
}

// ===== ОБЕРТЫВАНИЕ ОШИБОК (WRAPPING) =====
func doSomething() error {
	err := riskyOperation(false)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении операции: %w", err)
	}
	return nil
}

// ===== РАСПАКОВКА И ПРОВЕРКА ТИПОВ ОШИБОК =====
func handleError(err error) {
	if err == nil {
		fmt.Println("Ошибки нет")
		return
	}

	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Println("Поймана MyError:", myErr.Code, myErr.Msg)
		return
	}

	fmt.Println("Неизвестная ошибка:", err)
}

func main() {
	// ===== ПРОСТАЯ ПРОВЕРКА ОШИБКИ =====
	res, err := divide(10, 2)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", res)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Поймана ошибка деления:", err)
	}

	// ===== ПОЛЬЗОВАТЕЛЬСКИЕ ОШИБКИ =====
	err = riskyOperation(false)
	if err != nil {
		fmt.Println("Произошла ошибка:", err)
	}

	// ===== WRAPPING И ПРОВЕРКА =====
	err = doSomething()
	if err != nil {
		fmt.Println("Обёрнутая ошибка:", err)
	}

	// проверка типа
	handleError(err)

	// ===== СОЗДАНИЕ И СРАВНЕНИЕ ОШИБОК =====
	ErrNotFound := errors.New("not found")
	ErrInvalid := errors.New("invalid input")

	err = ErrNotFound
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Ошибка типа ErrNotFound")
	}

	if !errors.Is(err, ErrInvalid) {
		fmt.Println("Это не ErrInvalid")
	}

	// ===== INLINE СОЗДАНИЕ И ВОЗВРАТ ОШИБКИ =====
	check := func(x int) error {
		if x < 0 {
			return fmt.Errorf("отрицательное значение: %d", x)
		}
		return nil
	}

	if err := check(-10); err != nil {
		fmt.Println("Inline ошибка:", err)
	}
}
