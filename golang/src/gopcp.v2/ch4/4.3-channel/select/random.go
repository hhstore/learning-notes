package main

import "fmt"

/****************************************************

                随机写入 chan 数据
说明:
	- 利用 select, 实现对 chan 通道写入随机次序的数据

****************************************************/

func main() {
	chanCap := 5
	intChan := make(chan int, chanCap)

	//
	// 随机次序的把 [1, 2, 3] 写入 chan
	// 注意每次写入的次序不确定
	//
	for i := 0; i < chanCap; i++ {
		select {
		case intChan <- 1: // 把 1 写入 chan
		case intChan <- 2: // 把 2 写入 chan
		case intChan <- 3: // 把 3 写入 chan

		}
	}

	// 读取 chan 的内容:
	for i := 0; i < chanCap; i++ {
		fmt.Printf("%d\n", <-intChan)
	}

}
