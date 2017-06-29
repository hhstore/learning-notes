package main

import (
	"log"
	"sync/atomic"
)

/***************************************************************
                    原子值类型: atomic.Value
 说明:
		- 有2个对外操作接口:
			- Load() 和 Store()

***************************************************************/

func main() {
	var countVal atomic.Value

	// 第一次操作:
	countVal.Store([]int{
		1, 3, 5, 7,
	})
	log.Println("countVal =", countVal.Load())

	// 第二次操作:
	anothorStore(countVal) // 对比:
	log.Println("countVal =", countVal.Load())

}

// 第二次操作:
func anothorStore(countVal atomic.Value) {
	// 尝试修改:
	countVal.Store([]int{
		2, 4, 6, 8,
	})

}
