package concurrent

import (
	"sync/atomic"

	"errors"
	"fmt"
)

/***************************************************************
                    接口: 并发安全的数组
 说明:
 	- 利用  atomic.Value 实现
	- 单元测试里, 作了并发安全测试

***************************************************************/

// 接口类型: 并发安全的数组
type ConcurrentArray interface {
	Set(index uint32, element int) (err error) // 设置指定下标的元素值
	Get(index uint32) (element int, err error) // 获取指定下标的元素值
	Len() uint32                               // 获取数组长度
}

//***************************************************************
//					接口实现: 并发安全的整型数组
// 说明:
//		-  利用  atomic.Value 实现
//
//***************************************************************

type intArray struct {
	length uint32       // 数组长度
	value  atomic.Value // 原子值类型
}

//
func NewConcurerntArray(length uint32) ConcurrentArray {
	array := intArray{}

	array.length = length
	array.value.Store(make([]int, array.length))

	return &array
}

//
func (array *intArray) Len() uint32 {
	return array.length
}

//
func (array *intArray) Set(index uint32, element int) (err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}

	if err = array.checkValue(); err != nil {
		return
	}

	// 不要这样做！否则会形成竞态条件！
	// oldArray := array.val.Load().([]int)
	// oldArray[index] = elem
	// array.val.Store(oldArray)

	newArray := make([]int, array.length)
	copy(newArray, array.value.Load().([]int)) // [关键步骤] 复制操作

	newArray[index] = element
	array.value.Store(newArray)
	return
}

//
func (array *intArray) Get(index uint32) (element int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}

	if err = array.checkValue(); err != nil {
		return
	}
	element = array.value.Load().([]int)[index]
	return
}

//***************************************************************

// 检查索引的有效性:
func (array *intArray) checkIndex(index uint32) error {
	if index >= array.length {
		return fmt.Errorf("Index out of range [0, %d]!", array.length)
	}
	return nil
}

// 检查原子值中, 是否已存在值.
func (array *intArray) checkValue() error {
	v := array.value.Load()
	if v == nil {
		return errors.New("Invalid int array!")
	}
	return nil
}
