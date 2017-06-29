package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/***************************************************************
                    只执行一次操作的类型: sync.Once
 说明:
 	- sync.Once
 		- once.Do() 该方法有效调用次数, 永远为1
	- 应用场景:
		- 仅需执行一次的任务
		- 数据块连接池的建立
		- 全局变量的延迟初始化

***************************************************************/

func main() {
	var count int
	var once sync.Once // 只执行一次
	max := rand.Intn(100)

	// 循环:
	for i := 0; i < max; i++ {
		// 只执行一次:
		once.Do(func() {
			count++
		})
	}

	fmt.Println("Count:", count)

}
