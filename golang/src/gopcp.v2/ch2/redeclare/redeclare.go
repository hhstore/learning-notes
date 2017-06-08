package main

import "fmt"

/*

说明:
	1. 作用域范围: 全局作用域, 函数作用域, 局部作用域
	2. 就近原则, 同名屏蔽.

*/

// 全局变量
var v = "1, 2, 3"

func main() {
	fmt.Printf("v[1]: %v\n", v)

	// 函数作用域, 同名屏蔽
	v := []int{1, 2, 3}
	fmt.Printf("v[2]: %v\n", v)

	if v != nil {
		// 局部作用域, 同名屏蔽
		var v = 123
		fmt.Printf("v[3]: %v\n", v)
	}

}
