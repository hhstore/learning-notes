package main

import (
	"fmt"
	"time"
)

/****************************************************
	  通道 + 定时器

说明:
	- 设定超时时间, 若超时, 则无法获取到数据.

****************************************************/

func main() {
	intChan := make(chan int, 1)

	// 发送操作:
	go func() {
		// 若注释掉 sleep() 对比结果差别.
		time.Sleep(time.Second)
		//
		intChan <- 1 // 写入数据到通道
	}()

	// 接收操作:
	select {
	case e := <-intChan:
		fmt.Printf("Received: %v\n", e)
	case <-time.NewTimer(time.Millisecond * 500).C: // 超时时间
		fmt.Println("Timeout!")
	}
}
