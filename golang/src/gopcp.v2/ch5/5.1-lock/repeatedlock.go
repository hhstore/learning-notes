package main

import (
	"fmt"
	"sync"
	"time"
)

/*
说明:
	- 互斥锁

*/

func main() {
	var mutex sync.Mutex

	fmt.Println("[main] Lock the lock.")
	mutex.Lock() // 加锁:
	fmt.Println("[main] The lock is locked.")

	//
	for i := 1; i <= 3; i++ {
		// go
		go func(i int) {
			fmt.Printf("[g%d] Lock the lock.\n", i)
			mutex.Lock()
			fmt.Printf("[g%d] The lock is locked.\n", i)
		}(i)
	}

	//
	time.Sleep(time.Second)
	fmt.Println("[main] Unlock the lock.")
	mutex.Unlock()
	fmt.Println("[main] The lock is unlocked.")

}
