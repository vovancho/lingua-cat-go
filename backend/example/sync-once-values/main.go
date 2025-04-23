package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Правильное использование OnceValues с двумя возвращаемыми значениями
	getConnectionInfo := sync.OnceValues(func() (string, int) {
		fmt.Println("Initializing connection...")
		time.Sleep(500 * time.Millisecond)
		return "postgres://localhost", 5432
	})

	// Многопоточный доступ
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			url, port := getConnectionInfo()
			fmt.Printf("Goroutine %d: connected to %s:%d\n", id, url, port)
		}(i)
	}

	wg.Wait()

	// Повторный вызов (использует кэшированные значения)
	url, port := getConnectionInfo()
	fmt.Printf("\nReused connection: %s:%d\n", url, port)
}
