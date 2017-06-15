package main

import (
	"fmt"
	"time"
)

//
// 说明:
// 	- go 发送子例程中, 当 chan 已满时, 发送 同步信号, 唤醒接收 go子例程.
// 	- 接收 go 子例程, 开始时, 阻塞等 发送子例程, 直到收到 同步信号, 才开始接收数据.
//
var strChan = make(chan string, 3) // 数据 channel, string 类型, 长度为3

func main() {
	syncChan1 := make(chan struct{}, 1) // 同步控制信号 channel
	syncChan2 := make(chan struct{}, 2) // 同步控制信号 channel

	//----------------------------------------------
	//           子go例程: 接收数据
	//----------------------------------------------
	go func() {
		// 阻塞等发送例程, 发送例程会发送控制信号.
		<-syncChan1

		fmt.Println("Received a sync signal and wait a second... [receiver]")
		time.Sleep(time.Second)

		//
		for {
			if elem, ok := <-strChan; ok {
				fmt.Println("Received:", elem, "[receiver]")
			} else {
				break
			}
		}

		//
		fmt.Println("Stopped. [receiver]")

		syncChan2 <- struct{}{} // 发送空值
	}()

	//
	//----------------------------------------------
	//           子go例程: 发送数据
	//----------------------------------------------
	go func() {
		for _, elem := range []string{"a", "b", "c", "d"} {
			strChan <- elem // 迭代, 依次把 [a, b, c, d] 存入 strChan
			fmt.Println("Sent:", elem, "[sender]")

			if elem == "c" { // 当 strChan 中, 已存入 c时, 通知 接收go例程, 准备接受数据.
				syncChan1 <- struct{}{} // 传入空值, 作为 通知信号.
				fmt.Println("Sent a sync signal. [sender]")
			}
		}

		//
		fmt.Println("Wait 2 seconds... [sender]")
		time.Sleep(time.Second * 2)

		close(strChan)          // 关闭 chan
		syncChan2 <- struct{}{} // 发送空值

	}()

	//----------------------------------------------
	//         main 例程: 阻塞等子go例程结束
	//----------------------------------------------
	<-syncChan2 // 阻塞等零值信号
	<-syncChan2 // 阻塞等零值信号
}
