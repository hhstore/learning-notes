package main

import "fmt"

func main() {
	dataChan := make(chan int, 5)
	syncChan1 := make(chan struct{}, 1)
	syncChan2 := make(chan struct{}, 2)

	//----------------------------------------------
	//           子go例程: 接收数据
	//----------------------------------------------
	go func() { // 用于演示接收操作。
		// 阻塞等
		<-syncChan1

		for {
			if elem, ok := <-dataChan; ok {
				fmt.Printf("Received: %d [receiver]\n", elem)
			} else {
				break
			}
		}
		fmt.Println("Done. [receiver]")
		syncChan2 <- struct{}{}
	}()

	//----------------------------------------------
	//           子go例程: 发送数据
	//----------------------------------------------
	go func() { // 用于演示发送操作。
		for i := 0; i < 5; i++ {
			dataChan <- i
			fmt.Printf("Sent: %d [sender]\n", i)
		}

		// 关闭channel
		close(dataChan)

		syncChan1 <- struct{}{}
		fmt.Println("Done. [sender]")
		syncChan2 <- struct{}{}
	}()

	//----------------------------------------------
	//         main 例程: 阻塞等子go例程结束
	//----------------------------------------------
	<-syncChan2
	<-syncChan2
}
