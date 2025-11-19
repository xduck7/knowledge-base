package main

import (
	"fmt"
	"runtime"
	"time"
)

// ===== GC (GARBAGE COLLECTOR) В GO =====
// GC автоматически управляет памятью, удаляя объекты,
// на которые больше нет ссылок. Прямого управления, как в C, нет.

// ===== ПРИМЕР: АЛЛОКАЦИИ, СБОРКА, СТАТИСТИКА =====
func allocate() {
	slice := make([]byte, 10*1024*1024) // 10 MB
	for i := range slice {
		slice[i] = byte(i % 256)
	}
	fmt.Println("Выделено 10MB")
}

func showMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func main() {
	fmt.Println("=== Старт программы ===")
	showMemStats()

	allocate()
	showMemStats()

	// Освобождаем ссылку
	runtime.GC() // Принудительный вызов GC (не рекомендуется в продакшене)
	fmt.Println("После ручного GC:")
	showMemStats()

	// ===== ПРИМЕР С УТЕЧКОЙ =====
	// Даже если объект не используется, но ссылка где-то хранится — GC его не тронет
	var leak []*int
	for i := 0; i < 1_000_000; i++ {
		x := i
		leak = append(leak, &x)
	}
	fmt.Println("Создано 1,000,000 ссылок (эмулируем утечку)")
	showMemStats()

	// Удаляем утечку и снова вызываем GC
	leak = nil
	runtime.GC()
	fmt.Println("После очистки и GC:")
	showMemStats()

	// ===== ТОНКИЙ МОМЕНТ: ЗАДЕРЖКА GC =====
	// GC не мгновенный: он асинхронен и работает в несколько фаз (mark-and-sweep).
	// Чтобы показать, что память возвращается не сразу:
	fmt.Println("Ожидание GC...")
	time.Sleep(2 * time.Second)
	runtime.GC()
	showMemStats()

	// ===== СИСТЕМНЫЕ НАСТРОЙКИ GC =====
	// Можно изменить частоту GC через GOGC (процент роста heap до запуска GC)
	// runtime/debug.SetGCPercent(50) // запустить GC при росте на 50%
	// debug.FreeOSMemory() // принудительно вернуть память ОС
}
