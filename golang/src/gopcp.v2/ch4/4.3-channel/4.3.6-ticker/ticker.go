package main

import (
	"fmt"
	"time"
)

/****************************************************
	  通道 + 断续器

说明:
	- 断续器的应用场景: 定时任务的触发器

****************************************************/

func main() {
	intChan := make(chan int, 1)
	ticker := time.NewTicker(time.Second)

	// 发送数据:
	go func() {
		// 随机写入:
		// ticker.C 到期后, 立即进入下一个周期:
		// 每隔1s, 发送一个数据到通道, 周而复始, 且不会主动停止
		// 接受者会控制结束收取的时间.
		//
		for _ = range ticker.C {
			// 随机次序往 chan 中写入如下值[1, 2, 3]
			select {
			case intChan <- 1: // 写入1
			case intChan <- 2: // 写入2
			case intChan <- 3: // 写入3
			}
		}
		//
		fmt.Println("End. [sender]") // 本条语句, 永远无法被执行
	}()

	// 接收数据:
	var sum int

	for e := range intChan {
		fmt.Printf("Received: %v\n", e)
		sum += e // 求和

		if sum > 10 {
			fmt.Printf("Got: %v\n", sum)
			break
		}
	}

	fmt.Println("End. [receiver]")

}
