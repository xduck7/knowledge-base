package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ===== СИНХРОНИЗАЦИЯ В GO =====
// В Go синхронизация между горутинами реализуется через:
// - sync.Mutex и sync.RWMutex — взаимное исключение доступа к данным;
// - sync.WaitGroup — ожидание завершения горутин;
// - sync/atomic — атомарные операции над числами;
// - каналы — безопасная передача данных между горутинами.

// ======================= MUTEX =======================

func exampleMutex() {
	var mu sync.Mutex
	counter := 0

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				mu.Lock()
				counter++
				fmt.Printf("[Mutex] Горутина %d увеличила счётчик: %d\n", id, counter)
				mu.Unlock()
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("[Mutex] Итоговое значение:", counter)
}

// ======================= RWMutex =======================

func exampleRWMutex() {
	var rw sync.RWMutex
	data := 0

	read := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()
		rw.RLock()
		fmt.Printf("[RWMutex] Читатель %d прочитал: %d\n", id, data)
		rw.RUnlock()
	}

	write := func(val int, wg *sync.WaitGroup) {
		defer wg.Done()
		rw.Lock()
		data = val
		fmt.Printf("[RWMutex] Записано значение: %d\n", data)
		rw.Unlock()
	}

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go read(i, &wg)
	}
	wg.Add(1)
	go write(42, &wg)
	wg.Add(1)
	go read(4, &wg)
	wg.Wait()
}

// ======================= ATOMIC =======================

func exampleAtomic() {
	var counter int64

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("[Atomic] Итоговое значение:", counter)
}

// ======================= WAITGROUP =======================

func exampleWaitGroup() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("[WaitGroup] Горутина %d начала работу\n", id)
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("[WaitGroup] Горутина %d завершила работу\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("[WaitGroup] Все горутины завершены")
}

// ======================= КАНАЛЫ =======================

func exampleChannels() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("[Channel] Отправка %d\n", i)
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()

	go func() {
		for val := range ch {
			fmt.Printf("[Channel] Получено %d\n", val)
		}
		close(done)
	}()

	<-done
	fmt.Println("[Channel] Все значения обработаны")
}

// ======================= MAIN =======================

func main() {
	fmt.Println("=== Синхронизация в Go ===")

	exampleMutex()
	fmt.Println()

	exampleRWMutex()
	fmt.Println()

	exampleAtomic()
	fmt.Println()

	exampleWaitGroup()
	fmt.Println()

	exampleChannels()
	fmt.Println()
}
