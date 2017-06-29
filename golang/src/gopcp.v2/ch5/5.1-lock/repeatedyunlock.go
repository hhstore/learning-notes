package main

import (
	"fmt"
	"sync"
)

/*

说明:
	1. 重复解锁, 导致抛出一个Panic, 且无法被手工捕捉到.
	2. 注意互斥锁的使用, 操作要配对.

*/

func main() {
	// go 1.8 之后, 重复解锁操作的Panic, 无法被捕捉.
	defer func() {
		fmt.Println("Try to recover the panic.")
		if p := recover(); p != nil {
			fmt.Printf("Recovered the panic(%#v).\n", p)
		}

	}()

	var mutex sync.Mutex // 互斥锁

	fmt.Println("Lock the lock.")
	mutex.Lock() // 加锁
	fmt.Println("The lock is locked.")

	//
	fmt.Println("Unlock the lock.")
	mutex.Unlock() // 解锁
	fmt.Println("The lock is unlocked.")

	// 重复解锁:
	// 抛出一个 Panic:
	fmt.Println("Unlock the lock again.")
	mutex.Unlock() // 解锁
}
