package main

import (
	"fmt"
	"time"
)

/****************************************************
	  通道 + 定时器

说明:
	- 定时器使用

****************************************************/

func main() {
	timer := time.NewTimer(time.Second * 2)
	fmt.Printf("Present time:\t %v.\n", time.Now())

	expirationTime := <-timer.C // 阻塞

	fmt.Printf("Expiration time:\t %v.\n", expirationTime)
	fmt.Printf("Stop timer: %v.\n", timer.Stop()) // 尝试主动关闭定时器
}
