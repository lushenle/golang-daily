package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *DriverPg
)

// DriverPg represents a PostgreSQL database driver connection
type DriverPg struct {
	conn string
}

// Connect returns the singleton instance of DriverPg
func Connect() *DriverPg {
	once.Do(func() {
		instance = &DriverPg{conn: "DriverConnectPostgres"}
	})
	return instance
}

func main() {
	// Simulate a delayed call to Connect
	go func() {
		time.Sleep(time.Millisecond * 600)
		fmt.Println(*Connect())
	}()

	// Create 100 goroutines.
	for i := 0; i < 100; i++ {
		go func(ix int) {
			time.Sleep(time.Millisecond * 60)
			fmt.Println(ix, " = ", Connect().conn)
		}(i)
	}

	fmt.Scanln()
}
