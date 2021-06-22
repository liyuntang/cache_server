package main

import (
	"fmt"
	"sync"
)

func main() {
	pool := &sync.Pool{}

	pool.Put(NewConnection(1))
	pool.Put(NewConnection(2))
	pool.Put(NewConnection(3))

	connection := pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
}