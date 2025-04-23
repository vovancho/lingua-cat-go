package main

import (
	"fmt"
	"sync"
)

func main() {
	// Общая переменная - счетчик
	var counter int

	// Создаем мьютекс для защиты counter
	var mutex sync.Mutex

	// WaitGroup для синхронизации завершения горутин
	var wg sync.WaitGroup

	// Запускаем 5 горутин, каждая увеличивает счетчик 1000 раз
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Захватываем мьютекс перед доступом к counter
				mutex.Lock()
				counter++
				// Освобождаем мьютекс после обновления
				mutex.Unlock()
			}
			fmt.Printf("Горутина %d: завершила работу\n", id)
		}(i)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Выводим финальное значение счетчика
	fmt.Printf("Финальное значение счетчика: %d\n", counter)
}
