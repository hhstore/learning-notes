package main

import (
	"fmt"
	"time"
)

/****************************************************
	  通道 + 定时器

说明:
	- 设定超时时间

****************************************************/

func main() {
	intChan := make(chan int, 1)

	// 发送数据:
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)

			// 写入数据到通道
			intChan <- i
		}

		//
		close(intChan)
	}()

	//
	timeout := time.Millisecond * 500
	var timer *time.Timer

	// 接收数据:
	for {
		if timer == nil {
			timer = time.NewTimer(timeout)
		} else {
			timer.Reset(timeout)
		}

		//
		select {
		case e, ok := <-intChan:
			if !ok {
				fmt.Println("End.")
				return
			}
			fmt.Printf("Received: %v\n", e)
		case <-timer.C:
			fmt.Println("Timeout!")
		}
	}

}
