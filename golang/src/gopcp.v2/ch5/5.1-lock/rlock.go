package main

import (
	"fmt"
	"sync"
	"time"
)

/*
说明:
	1. 读写锁的使用.
	2. 允许多读单写. 多个 读操作, 可以同时进行, 只允许单个 写操作.
	3. 写操作, 始终和 读操作 互斥.

*/

func main() {
	var rwm sync.RWMutex

	// 启动 goroutine:
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Printf("[g%d] Try to lock for reading...\n", i)
			rwm.RLock() // 读锁定
			fmt.Printf("[g%d] Locked.\n", i)

			//
			time.Sleep(time.Second * 2)
			fmt.Printf("[g%d] Try to unlock for reading...\n", i)

			//
			rwm.RUnlock() // 读解锁
			fmt.Printf("[g%d] Unlocked for reading.\n", i)
		}(i)

	}

	// 等待
	time.Sleep(time.Millisecond * 100)

	fmt.Println("[main] Try to lock for writing...")
	rwm.Lock() // 写锁定, 阻塞等待
	fmt.Println("[main] Locked for writing...")
}
