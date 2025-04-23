package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Создаем OnceValue для ленивой инициализации
	getConfig := sync.OnceValue(func() map[string]string {
		fmt.Println("Initializing config...")
		// Имитация тяжелой инициализации
		time.Sleep(500 * time.Millisecond)
		return map[string]string{
			"host":    "localhost",
			"port":    "8080",
			"timeout": "30s",
		}
	})

	// Многопоточный доступ к конфигу
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Все горутины получат один и тот же результат
			cfg := getConfig()
			fmt.Printf("Goroutine %d got config: %v\n", id, cfg)
		}(i)
	}

	wg.Wait()

	// Повторный вызов - инициализация не повторяется
	fmt.Println("Second call:", getConfig()["host"])
}
