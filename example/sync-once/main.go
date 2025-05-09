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
