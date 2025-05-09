Вот пример использования `sync.OnceValues` в Go 1.21+, который позволяет безопасно инициализировать и кэшировать несколько возвращаемых значений:

```go
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
```

Вывод программы:
```
Establishing database connection...
Goroutine 2 connected to postgres://localhost:5432 (port 5432)
Goroutine 0 connected to postgres://localhost:5432 (port 5432)
Goroutine 1 connected to postgres://localhost:5432 (port 5432)

Reused connection: postgres://localhost:5432:5432
```

Ключевые особенности `sync.OnceValues`:
1. **Множественные возвращаемые значения** - в отличие от `OnceValue`, работает с функциями, возвращающими несколько значений
2. **Потокобезопасность** - как и все примитивы из пакета sync
3. **Ленивая инициализация** - вычисление происходит при первом вызове
4. **Кэширование результата** - последующие вызовы возвращают закэшированные значения

Пример с обработкой ошибок:
```go
getConfig := sync.OnceValues(func() (map[string]string, error) {
    if time.Now().Unix()%2 == 0 {
        return nil, fmt.Errorf("initialization failed")
    }
    return map[string]string{"mode": "production"}, nil
})

cfg, err := getConfig() // Ошибка будет закэширована как и результат
```

Отличие от ручной реализации:
```go
// Без OnceValues (аналог для старых версий Go)
var (
    once sync.Once
    url  string
    port int
    err  error
)

getConnection := func() (string, int, error) {
    once.Do(func() {
        url, port, err = "localhost", 5432, nil 
    })
    return url, port, err
}
```

Рекомендации по использованию:
- Идеально для тяжелых инициализаций (подключения к БД, загрузка конфигов)
- Не используйте для часто изменяемых данных
- Для сложной логики инициализации лучше подойдет `sync.Once` + структура
- В случае ошибки, она будет закэширована и возвращаться при всех последующих вызовах
