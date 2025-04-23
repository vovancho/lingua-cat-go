Вот рабочий пример использования нового `sync.OnceValue` из Go 1.21+:

```go
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
```

Вывод программы:
```
Initializing config...
Goroutine 4 got config: map[host:localhost port:8080 timeout:30s]
Goroutine 0 got config: map[host:localhost port:8080 timeout:30s]
Goroutine 1 got config: map[host:localhost port:8080 timeout:30s]
Goroutine 2 got config: map[host:localhost port:8080 timeout:30s]
Goroutine 3 got config: map[host:localhost port:8080 timeout:30s]
Second call: localhost
```

Ключевые особенности `sync.OnceValue`:
1. Гарантирует однократное выполнение функции инициализации
2. Потокобезопасен - можно вызывать из нескольких горутин
3. Возвращает функцию-геттер, которая кэширует результат
4. Эффективнее ручной реализации с `sync.Once` для простых случаев

Альтернативная реализация для версий до Go 1.21:
```go
func manualOnceValue() {
	var once sync.Once
	var config map[string]string

	getConfig := func() map[string]string {
		once.Do(func() {
			config = map[string]string{"host": "localhost"}
		})
		return config
	}
}
```

Отличия от обычного `sync.Once`:
- Более удобный API для простых случаев
- Не требует ручного управления состоянием
- Возвращает готовую функцию-геттер
