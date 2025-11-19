package main

import (
	"fmt"
	"runtime"
	"time"
)

// ===== ПАКЕТ runtime =====
// Пакет runtime предоставляет низкоуровневые возможности Go:
// - работа с горутинами и потоками
// - информация о памяти и сборке мусора
// - управление планировщиком (GOMAXPROCS)
// - стек, номер текущего потока, версии Go и т.д.

func main() {
	fmt.Println("=== Пакет runtime ===")

	// ===== BASIC INFO =====
	fmt.Println("Версия Go:", runtime.Version())
	fmt.Println("Число CPU доступных планировщику:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// ===== АКТИВНЫЕ ГОРУТИНЫ =====
	fmt.Println("Активных горутин:", runtime.NumGoroutine())

	// ===== ОСНОВЫ GC =====
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("HeapAlloc = %v KB, NumGC = %v\n", mem.HeapAlloc/1024, mem.NumGC)

	// ===== ПРИМЕР: ПЛАНИРОВКА =====
	done := make(chan struct{})
	go func() {
		fmt.Println("Горутина работает")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Горутина завершила работу")
		close(done)
	}()

	// Показать планировщик
	fmt.Println("Перед yield: активных горутин:", runtime.NumGoroutine())
	runtime.Gosched() // уступка планировщику
	fmt.Println("После yield: активных горутин:", runtime.NumGoroutine())

	<-done
	fmt.Println("Основная программа завершена")

	// ===== STACK TRACE =====
	fmt.Println("\nПример stack trace текущей горутины:")
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	fmt.Println(string(buf[:n]))

	// ===== ПРИНУДИТЕЛЬНЫЙ GC =====
	runtime.GC()
	runtime.ReadMemStats(&mem)
	fmt.Printf("\nПосле GC: HeapAlloc = %v KB, NumGC = %v\n", mem.HeapAlloc/1024, mem.NumGC)
}
