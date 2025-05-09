# Пример использования `sync.WaitGroup` в Go для синхронизации горутин
## Пример:
```go
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
```

### Вывод (примерный)
```
Worker 3 starting
Worker 1 starting
Worker 2 starting
Worker 1 done
Worker 2 done
Worker 3 done
All workers completed
```

### Как это работает
1. `wg.Add(1)` – увеличивает счетчик на 1 перед запуском каждой горутины.
2. `defer wg.Done()` – уменьшает счетчик при завершении работы горутины (defer гарантирует выполнение даже при панике).
3. `wg.Wait()` – блокирует выполнение основной программы, пока счетчик не станет 0.

### Важные моменты
- **`WaitGroup` передается по указателю (`*sync.WaitGroup`), иначе будет копия, и счетчик не изменится.**
- **`Add()` лучше вызывать в главной горутине, а не внутри worker, чтобы избежать гонки.**
- **`Done()` можно заменить на `Add(-1)`, но `Done()` читается лучше.**

Если нужно обрабатывать ошибки из горутин, можно использовать каналы или errgroup (из golang.org/x/sync/errgroup).
