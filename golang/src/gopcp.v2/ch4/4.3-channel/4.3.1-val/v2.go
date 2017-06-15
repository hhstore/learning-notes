package main

import (
	"fmt"
	"time"
)

type Counter struct {
	count int
}

var mapChan = make(chan map[string]Counter, 1)

func main() {
	syncChan := make(chan struct{}, 2) // 同步专用通道, 用于调度

	// 接收操作:
	go func() {
		for {
			if elem, ok := <-mapChan; ok {
				counter := elem["count"]
				counter.count++
			} else {
				break
			}
		}
		//
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{} // 执行结束信号
	}()

	// 发送操作:
	go func() {
		countMap := map[string]Counter{
			"count": Counter{},
		}

		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("the count map: %v. [sender]\n", countMap)
		}

		//
		close(mapChan)
		syncChan <- struct{}{} // 执行结束信号

	}()

	// 阻塞等:
	<-syncChan
	<-syncChan
}
