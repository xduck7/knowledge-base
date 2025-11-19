package main

import (
	"fmt"
	"time"
)

// ===== ОСНОВЫ КАНАЛОВ =====
// Канал — это типизированная труба для обмена данными между горутинами.

func main() {
	// Создание канала для int
	ch := make(chan int)

	// запись и чтение из канала в разных горутинах
	go func() {
		fmt.Println("Горутина: отправляю 42 в канал")
		ch <- 42
	}()

	val := <-ch // блокирует выполнение, пока значение не будет получено
	fmt.Println("main: получил значение", val)

	// ===== БУФЕРИЗИРОВАННЫЕ КАНАЛЫ =====
	// канал с буфером (неблокирующая отправка, пока буфер не заполнен)
	buf := make(chan string, 2)
	buf <- "one"
	buf <- "two"

	fmt.Println("Buffered receive 1:", <-buf)
	fmt.Println("Buffered receive 2:", <-buf)

	// ===== ЗАКРЫТИЕ КАНАЛОВ =====
	close(buf)
	// после закрытия можно читать оставшиеся данные, но нельзя писать
	// чтение из закрытого канала даёт zero value
	empty := <-buf
	fmt.Println("Reading from closed channel:", empty)

	// ===== RANGE ПО КАНАЛУ =====
	nums := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			nums <- i
			time.Sleep(200 * time.Millisecond)
		}
		close(nums)
	}()
	for n := range nums {
		fmt.Println("Range got:", n)
	}

	// ===== SELECT =====
	// позволяет слушать несколько каналов одновременно
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(400 * time.Millisecond)
		ch1 <- "message from ch1"
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "message from ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}

	// ===== ДВУНАПРАВЛЕННЫЕ И ОДНОНАПРАВЛЕННЫЕ КАНАЛЫ =====
	// обычно каналы двунаправленные, но функции можно ограничивать по направлению

	sendOnly := make(chan<- int) // только отправка
	recvOnly := make(<-chan int) // только получение

	fmt.Printf("Send-only type: %T\n", sendOnly)
	fmt.Printf("Recv-only type: %T\n", recvOnly)
}
