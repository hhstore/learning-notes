package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{
		"Eric", "Harry", "Robert", "Jim", "Mark",
	}

	// 注意输出结果:
	// go func() 是在 for 迭代结束, 才执行的.
	for _, name := range names {

		// 创建 5 个子例程:
		go func() {
			fmt.Printf("hello, %s!\n", name)
		}()

	}

	//
	time.Sleep(time.Millisecond)

}
