package main

import (
	"fmt"
	"time"
)

func main() {
	name := "Eric"

	// 注意, 子例程的输出结果:
	go func() {
		fmt.Printf("hello, %s!\n", name)
	}()

	// sleep()的位置, 会影响 go 子例程输出结果.
	//time.Sleep(time.Millisecond)
	name = "Harry"
	time.Sleep(time.Millisecond)
}
