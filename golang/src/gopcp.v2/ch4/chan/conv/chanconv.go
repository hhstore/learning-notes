package main

import "fmt"

/****************************************************

           通道类型转换
说明:
	- 双向通道 和 单向通道, 并不可以强制转换类型
	- 通道定义完, 类型就不可更改

****************************************************/

func main() {
	var ok bool

	// 双向通道
	ch := make(chan int, 1)

	_, ok = interface{}(ch).(<-chan int)
	fmt.Println("chan int => <-chan int:", ok)

	_, ok = interface{}(ch).(chan<- int)
	fmt.Println("chan int => chan<- int:", ok)

	// 单向通道: 只写
	sch := make(chan<- int, 1)
	_, ok = interface{}(sch).(chan int)
	fmt.Println("chan<- int => chan int:", ok)

	// 单向通道: 只读
	rch := make(<-chan int, 1)
	_, ok = interface{}(rch).(chan int)
	fmt.Println("<-chan int => chan int:", ok)

}
