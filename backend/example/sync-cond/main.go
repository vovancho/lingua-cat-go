package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Создаем мьютекс и условие
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)

	// Переменная состояния
	ready := false

	// WaitGroup для синхронизации
	var wg sync.WaitGroup

	// Запускаем 3 горутины-ожидателя
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mutex.Lock()
			// Ожидаем, пока ready не станет true
			for !ready {
				fmt.Printf("Горутина %d: ожидает\n", id)
				cond.Wait() // Блокируется до сигнала
			}
			fmt.Printf("Горутина %d: условие выполнено, продолжаю\n", id)
			mutex.Unlock()
		}(i)
	}

	// Даем ожидающим горутинам время на запуск
	time.Sleep(100 * time.Millisecond)

	// Отправляем сигнал всем ожидающим горутинам
	fmt.Println("Основная горутина: отправляю сигнал")
	mutex.Lock()
	ready = true
	cond.Broadcast() // Разблокирует все ожидающие горутины
	mutex.Unlock()

	// Ждем завершения всех горутин
	wg.Wait()
	fmt.Println("Все горутины завершили работу")
}
