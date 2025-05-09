package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Общая переменная - счетчик
	var counter int64

	// WaitGroup для синхронизации завершения горутин
	var wg sync.WaitGroup

	// Запускаем 5 горутин, каждая увеличивает счетчик 1000 раз
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Атомарно увеличиваем счетчик
				atomic.AddInt64(&counter, 1)
			}
			fmt.Printf("Горутина %d: завершила работу\n", id)
		}(i)
	}

	// Запускаем горутину для чтения счетчика
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			// Атомарно читаем значение счетчика
			value := atomic.LoadInt64(&counter)
			fmt.Printf("Читатель: текущее значение счетчика = %d\n", value)
			time.Sleep(100 * time.Millisecond) // Имитация работы
		}
	}()

	// Ждем завершения всех горутин
	wg.Wait()

	// Выводим финальное значение счетчика
	fmt.Printf("Финальное значение счетчика: %d\n", atomic.LoadInt64(&counter))
}
