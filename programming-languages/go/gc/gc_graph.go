package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"runtime/trace"
	"time"
)

const (
	elements    = 100000 // размер для массива
	bufferCount = 100    // кол-во массивных объектов
)

type Buffer struct {
	Weight []int
}

func addBuffer() *Buffer {

	// Создание тяжелого массива
	buf := make([]int, elements)
	for i := 0; i < elements; i++ {
		buf[i] = i
	}
	return &Buffer{
		Weight: buf,
	}
}

func main() {
	// Запись в trace файл.
	f, err := os.Create("trace.out")
	if err != nil {
		panic("Error creating trace")
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic("Error starting trace")
	}

	defer fmt.Println("тяжелая работа закончилась..\nВведите 'go tool trace trace.out', чтоб увидеть результат.\nВ открывшемся окне нажмите `View trace by proc`.\nЕсли выдает ошибку, то установите пакет `graphviz`.\nНапример `brew install graphviz`")
	defer trace.Stop()

	// процент, при котором стартует сборщик мусора
	// при изменении этого параметра в меньшую сторону, график будет более зубчатыми, так как очистка будет происходить чаще
	// по умолчанию GC запускается при заполненности 100%
	debug.SetGCPercent(50)

	// создаем объекты, которые в будущем очистятся GC
	for i := 0; i < bufferCount; i++ {
		buf := addBuffer() // создаем буфер
		buf.Weight[0] = 1  // чтоб не было "buf declared and not used"
		time.Sleep(10)     // искусственная задержка
	}
}
