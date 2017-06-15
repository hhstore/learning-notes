package main

import (
	"time"
	//"runtime"
)

func main() {
	go println("Go! goroutine!")

	// 等待
	//runtime.Gosched()
	time.Sleep(time.Millisecond)

}
