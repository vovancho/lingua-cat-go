package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Общая переменная - счетчик
	var counter int

	// Создаем RWMutex для защиты counter
	var rwMutex sync.RWMutex

	// WaitGroup для синхронизации завершения горутин
	var wg sync.WaitGroup

	// Запускаем 3 горутины-читателя
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				// Захватываем мьютекс для чтения
				rwMutex.RLock()
				fmt.Printf("Читатель %d: текущее значение счетчика = %d\n", id, counter)
				rwMutex.RUnlock()
				time.Sleep(50 * time.Millisecond) // Имитация работы
			}
		}(i)
	}

	// Запускаем 2 горутины-писателя
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				// Захватываем мьютекс для записи
				rwMutex.Lock()
				counter++
				fmt.Printf("Писатель %d: увеличил счетчик до %d\n", id, counter)
				rwMutex.Unlock()
				time.Sleep(100 * time.Millisecond) // Имитация работы
			}
		}(i)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Выводим финальное значение счетчика
	fmt.Printf("Финальное значение счетчика: %d\n", counter)
}
