package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker представляет структуру для выполнения задач
type Worker struct {
	ID    int   // Идентификатор воркера
	Tasks []int // Список выполненных задач
}

func (w *Worker) Reset() {
	// Сбрасываем состояние для повторного использования
	w.ID = 0
	w.Tasks = w.Tasks[:0]
}

func main() {
	// Создаем sync.Pool для управления структурами Worker
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Создаем новый Worker")
			return &Worker{
				Tasks: make([]int, 0, 10), // Инициализируем срез с начальной емкостью
			}
		},
	}

	// WaitGroup для синхронизации завершения горутин
	var wg sync.WaitGroup

	// Запускаем 5 горутин, каждая работает с Worker из пула
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Получаем Worker из пула
			worker := pool.Get().(*Worker)
			fmt.Printf("Горутина %d: получил Worker %p с ID=%d\n", id, worker, worker.ID)

			// Устанавливаем ID и добавляем задачу
			worker.ID = id
			worker.Tasks = append(worker.Tasks, id*10)

			fmt.Printf("Горутина %d: обновил Worker %p, Tasks=%v\n", id, worker, worker.Tasks)

			// Сбрасываем состояние перед возвращением в пул
			worker.Reset()

			// Возвращаем Worker в пул
			pool.Put(worker)
			fmt.Printf("Горутина %d: вернул Worker %p в пул\n", id, worker)

			time.Sleep(50 * time.Millisecond) // Имитация работы
		}(i)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Проверяем повторное использование Worker
	fmt.Println("Проверяем повторное использование:")
	worker := pool.Get().(*Worker)
	fmt.Printf("Основная горутина: получил Worker %p с ID=%d, Tasks=%v\n", worker, worker.ID, worker.Tasks)
	pool.Put(worker)
}
