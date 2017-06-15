package main

import (
	"fmt"
	"time"
)

/****************************************************
		缓冲通道 和 非缓冲通道

说明:
	- 两种通道区别:
		- 容量为0, 为 非缓冲通道
		- 容量不为0, 为 缓冲通道
	- 本例为非缓冲通道:
		- 注意 chan 的长度 为 0.
		- 利用 channel 通信
		- 发送 go 例程
		- 接收部分, 利用 select{}, 实现对 channel 数据收发.

****************************************************/

func main() {
	sendingInterval := time.Second       // 发送间隔时间
	receptionInterval := time.Second * 2 // 接收间隔时间
	intChan := make(chan int, 0)         // 非缓冲通道, 用于数据传输. [整型, size = 0, 对比 string, 奇怪?]

	fmt.Println("chan size:", len(intChan)) // ? [size = 0]

	//----------------------------------------------
	//            向 chan 发送数据
	//----------------------------------------------
	go func() {
		var ts0, ts1 int64

		for i := 1; i <= 5; i++ {
			intChan <- i // 向 chan 发送数据
			//
			ts1 = time.Now().Unix()

			//
			if ts0 == 0 {
				fmt.Println("\tSent:", i)

			} else {
				fmt.Printf("\tSent: %d [interval: %d s]\n", i, ts1-ts0)
			}

			//
			ts0 = time.Now().Unix()
			time.Sleep(sendingInterval)
		}
		//
		close(intChan)
	}()

	//----------------------------------------------
	//             从 chan 读取数据
	//----------------------------------------------
	var ts0, ts1 int64

Loop:
	for {
		//
		select {
		// 输出数据:
		case v, ok := <-intChan: // 从 chan 读取数据
			if !ok {
				break Loop
			}

			ts1 = time.Now().Unix()
			if ts0 == 0 {
				fmt.Println("Received:", v) // 输出数据

			} else {
				fmt.Printf("Received: %d [interval: %d s]\n", v, ts1-ts0) // 输出数据
			}
		}

		//
		ts0 = time.Now().Unix()
		time.Sleep(receptionInterval) // 接收间隔时间
	}
	//
	fmt.Println("End.")

}
