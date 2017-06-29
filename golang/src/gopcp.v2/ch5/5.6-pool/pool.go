package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

/***************************************************************
                    临时对象池: sync.Pool
 说明:
 	- sync.Pool
 		- Get() 和 Put()
 		- 用途: 独特的同步工具, 对 GC 友好
		- 临时对象池的实例, 不应该被复制

	- 手动 GC:
		- debug.SetGCPercent(100)
		- runtime.GC()

***************************************************************/

func main() {
	// 禁用GC，并保证在main函数执行结束前恢复GC。
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	var count int32

	//
	fn := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}

	// 临时对象池:
	pool := sync.Pool{
		New: fn,
	}

	//--------------------------------------------
	// New 字段值的作用。
	v1 := pool.Get()
	fmt.Println("Value 1:", v1)

	// 临时对象池的存取。
	pool.Put(10)
	pool.Put(11)
	pool.Put(12)

	//
	v2 := pool.Get()
	fmt.Println("Value 2:", v2)

	// GC: 垃圾回收对临时对象池的影响。
	debug.SetGCPercent(100)
	runtime.GC()

	// 由上述执行 GC 回收. Get()得到值为1+1
	v3 := pool.Get()
	//v3 = pool.Get()	// 再次获取, +1
	fmt.Println("Vaule 3:", v3)

	//
	pool.New = nil

	//
	v4 := pool.Get()
	fmt.Println("Vaule 4:", v4)

}
