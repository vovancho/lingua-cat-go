package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Создаем sync.Map для хранения пар ключ-значение
	var m sync.Map

	// WaitGroup для синхронизации завершения горутин
	var wg sync.WaitGroup

	// Запускаем 3 горутины для записи
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 1; j <= 3; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := id * j
				m.Store(key, value)
				fmt.Printf("Писатель %d: записал %s = %d\n", id, key, value)
				time.Sleep(50 * time.Millisecond) // Имитация работы
			}
		}(i)
	}

	// Запускаем 2 горутины для чтения
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 1; j <= 3; j++ {
				key := fmt.Sprintf("key-%d-%d", j, j)
				if value, ok := m.Load(key); ok {
					fmt.Printf("Читатель %d: прочитал %s = %v\n", id, key, value)
				} else {
					fmt.Printf("Читатель %d: ключ %s еще не существует\n", id, key)
				}
				time.Sleep(100 * time.Millisecond) // Имитация работы
			}
		}(i)
	}

	// Запускаем горутину для удаления
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond) // Ждем, чтобы некоторые записи произошли
		key := "key-1-1"
		m.Delete(key)
		fmt.Printf("Удалитель: удалил ключ %s\n", key)
	}()

	// Ждем завершения всех горутин
	wg.Wait()

	// Выводим все оставшиеся пары ключ-значение
	fmt.Println("Окончательное содержимое sync.Map:")
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true
	})
}
