package main

import (
	"fmt"
)

/****************************************************

                通道数据发送和获取
说明:
	- 通道关闭后, 其中存储的数据, 依然可以被读取到.
	- 利用 select 获取数据

****************************************************/

func main() {
	intChan := make(chan int, 10)

	// 往 chan 中 发送数据:
	for i := 0; i < 10; i++ {
		intChan <- i
	}

	close(intChan) // 通道关闭后, 其中存储的值, 依然可以被读取到.

	// 同步用的信号通道.
	syncChan := make(chan struct{}, 1)

	// go 数据接收子例程:
	go func() {
	Loop:

		for {
			select {
			case e, ok := <-intChan:
				if !ok {
					fmt.Println("End.")
					break Loop
				}
				fmt.Printf("Received: %v\n", e)
			}
		}
		//
		syncChan <- struct{}{}

	}()

	//
	<-syncChan
}
