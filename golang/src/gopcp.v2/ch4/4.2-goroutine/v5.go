package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{
		"Eric", "Harry", "Robert", "Jim", "Mark",
	}

	for _, name := range names {
		// 传参方式, 创建5个子例程:
		go func(who string) {
			fmt.Printf("hello, %s!\n", who)
		}(name)
	}

	//
	time.Sleep(time.Millisecond)
}
