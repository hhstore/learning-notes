package main

import (
	"fmt"
	"time"
)

// 全局变量:
//	- 用途: 通信的数据.
var strChan = make(chan string, 3)

func main() {
	syncChan1 := make(chan struct{}, 1) // 信号控制
	syncChan2 := make(chan struct{}, 2) // 信号控制

	// 启动 2个 子 goroutine, 并进行调度:
	go receive(strChan, syncChan1, syncChan2) // 注意实参和形参
	go send(strChan, syncChan1, syncChan2)

	//----------------------------------------------
	//         main 例程: 阻塞等子go例程结束
	//----------------------------------------------
	<-syncChan2 // 阻塞, 等结束信号(零值)
	<-syncChan2 // 阻塞, 等结束信号(零值)
}

//
// 接收者:
//	- 注意形参: 只读通道, 只写通道
//
func receive(strCh <-chan string, readCh <-chan struct{}, writeCh chan<- struct{}) {
	<-readCh // 阻塞, 等待发送者的通知信号

	//
	fmt.Println("Received a sync signal and wait a second... [receiver]")
	time.Sleep(time.Second)

	//
	// 对比 v2版本:
	// 	- 通过for循环, 迭代 chan 实现:
	//
	for elem := range strCh {
		fmt.Println("Received:", elem, "[receiver]")
	}

	//
	fmt.Println("Stopped. [receiver]")
	writeCh <- struct{}{} // 发送本 goroutine 结束信号(零值).
}

//
// 发送者:
//	- 注意形参: 单向通道, 只写
//
func send(strCh chan<- string, writeCh1 chan<- struct{}, writeCh2 chan<- struct{}) {
	for _, elem := range []string{"a", "b", "c", "d"} {
		strCh <- elem // 存入数据
		fmt.Println("Sent:", elem, "[sender]")

		//
		if elem == "c" {
			writeCh1 <- struct{}{} // 给接受者, 发送控制信号(零值)
			fmt.Println("Sent a sync signal. [sender]")
		}
	}

	//
	fmt.Println("Wait 2 seconds... [sender]")
	time.Sleep(time.Second * 2)
	close(strCh)
	writeCh2 <- struct{}{} // 发送本 goroutine 结束信号(零值).
}
