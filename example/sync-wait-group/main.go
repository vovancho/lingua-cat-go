package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second) // Имитация работы
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // Увеличиваем счетчик перед запуском горутины
		go worker(i, &wg)
	}

	wg.Wait() // Ждем завершения всех горутин
	fmt.Println("All workers completed")
}
