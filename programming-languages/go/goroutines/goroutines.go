package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ===== ГОРУТИНЫ В GO =====
// Горутина — это легковесный поток исполнения, управляемый рантаймом Go.
// Запускается оператором `go`, планируется M:N планировщиком (много горутин на несколько системных потоков).

// worker — пример фоновой горутины
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Горутина %d запущена\n", id)
	time.Sleep(time.Duration(id) * 200 * time.Millisecond)
	fmt.Printf("Горутина %d завершена\n", id)
}

// infinite — пример «залипшей» горутины (для демонстрации планировщика)
func infinite() {
	for {
		runtime.Gosched() // уступает выполнение другим
	}
}

func main() {
	fmt.Println("=== Демонстрация горутин ===")

	// ===== ПРОСТОЙ ЗАПУСК ГОРУТИН =====
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("Все горутины завершены")

	// ===== АНОНИМНАЯ ГОРУТИНА =====
	done := make(chan struct{})
	go func() {
		fmt.Println("Анонимная горутина работает...")
		time.Sleep(time.Second)
		fmt.Println("Анонимная горутина завершена")
		close(done)
	}()
	<-done

	// ===== ДЕМОНСТРАЦИЯ ПАРАЛЛЕЛИЗМА =====
	// По умолчанию Go может выполнять горутины параллельно на нескольких ядрах
	prev := runtime.GOMAXPROCS(0)
	fmt.Printf("GOMAXPROCS (число CPU, доступных планировщику): %d\n", prev)

	// Запускаем интенсивные задачи
	wg.Add(2)
	go func() {
		defer wg.Done()
		sum := 0
		for i := 0; i < 1e7; i++ {
			sum += i
		}
		fmt.Println("Горутина A завершила вычисления")
	}()
	go func() {
		defer wg.Done()
		product := 1
		for i := 1; i < 1000; i++ {
			product *= i
		}
		fmt.Println("Горутина B завершила вычисления")
	}()
	wg.Wait()

	// ===== ПОДСЧЁТ АКТИВНЫХ ГОРУТИН =====
	fmt.Println("Активных горутин:", runtime.NumGoroutine())

	// ===== ОСТОРОЖНО С БЕСКОНЕЧНЫМИ ЦИКЛАМИ =====
	// Если запустить без yield (`runtime.Gosched()`), программа «замрёт»
	go infinite()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Активных горутин после запуска infinite():", runtime.NumGoroutine())

	// ===== ОЖИДАНИЕ С ЗАКРЫТИЕМ КАНАЛА =====
	exit := make(chan struct{})
	go func() {
		time.Sleep(time.Second)
		close(exit)
	}()
	fmt.Println("Ожидание сигнала завершения...")
	<-exit
	fmt.Println("Поступил сигнал")

	// ===== Работа с runtime.Goexit() =====

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("вызываю runtime.Goexit()")
	runtime.Goexit()

	// Почему программа продолжает выполняться и висит?
}
