# Пример использования `sync.Once` в Golang

Вот пример использования `sync.Once` в Go, который гарантирует, что определенная операция будет выполнена ровно **один раз**, даже если её вызывают из нескольких горутин.

## 🔹 Основные случаи применения:
- Инициализация ресурсов (например, подключение к БД)
- Ленивая загрузка (отложенная инициализация)
- Создание синглтона (единственного экземпляра объекта)

## 📌 Пример: Однократная инициализация
```go
package main

import (
	"fmt"
	"sync"
)

var (
	once     sync.Once
	instance *Database
)

// Database - пример структуры, которая должна быть создана один раз
type Database struct {
	connection string
}

// connect имитирует подключение к БД
func connect() *Database {
	fmt.Println("Connecting to the database...")
	return &Database{connection: "postgres://user:pass@localhost/db"}
}

// getDatabaseInstance возвращает синглтон Database
func getDatabaseInstance() *Database {
	once.Do(func() {
		instance = connect()
	})
	return instance
}

func main() {
	// Многократные вызовы getDatabaseInstance() приведут к однократному созданию Database
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := getDatabaseInstance()
			fmt.Printf("Goroutine %d: Database instance: %v\n", id, db.connection)
		}(i)
	}
	wg.Wait()
}
```

## 🔹 Вывод (примерный):
```
Connecting to the database...
Goroutine 2: Database instance: postgres://user:pass@localhost/db
Goroutine 4: Database instance: postgres://user:pass@localhost/db
Goroutine 0: Database instance: postgres://user:pass@localhost/db
Goroutine 1: Database instance: postgres://user:pass@localhost/db
Goroutine 3: Database instance: postgres://user:pass@localhost/db
```

## 📌 Пример: Ленивая загрузка конфига
```go
var (
	configOnce sync.Once
	appConfig  map[string]string
)

func loadConfig() {
	fmt.Println("Loading config...")
	appConfig = map[string]string{
		"host": "localhost",
		"port": "8080",
	}
}

func GetConfig() map[string]string {
	configOnce.Do(loadConfig)
	return appConfig
}

func main() {
	// Первый вызов загрузит конфиг
	fmt.Println(GetConfig()["host"]) // "localhost"
	// Последующие вызовы используют уже загруженный конфиг
	fmt.Println(GetConfig()["port"]) // "8080"
}
```

## 💡 Как это работает?
1. sync.Once гарантирует, что переданная в `Do()` функция выполнится только один раз, даже если `Do()` вызывается из нескольких горутин.
2. **Потокобезопасность**: Внутренняя блокировка (`mutex`) предотвращает гонки.
3. **Повторные вызовы** `Do()` не блокируются и не выполняют функцию повторно.

## ⚠️ Важные нюансы:
- Если функция в `Do()` паникует, последующие вызовы не будут её выполнять снова
- **Не используйте `sync.Once` для операций, которые должны повторяться** (например, переподключение к БД при разрыве)
- Для переинициализации нужно создать новый `sync.Once`

## 🏆 Когда применять?
- Инициализация глобальных переменных
- Кэширование тяжелых ресурсов
- Создание thread-safe синглтонов

Если вам нужно выполнять действие периодически, используйте `sync.Mutex` или другие механизмы синхронизации.
