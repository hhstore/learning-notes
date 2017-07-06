package ccmap

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"unsafe"
)

/***************************************************************
                    接口声明:

***************************************************************/
// 接口声明: 并发安全的 key-value 对
type Pair interface {
	linkedPair // 内嵌

	Key() string                          // 返回键的值
	Hash() uint64                         // 返回键的哈希值
	Element() interface{}                 //返回元素的值
	SetElement(element interface{}) error // 设置元素的值
	Copy() Pair                           // 生成一个 key-value 对的副本
	String() string                       // 返回 key-value 对的字符串表示
}

// 接口声明: 单向链接的 key-value 对
type linkedPair interface {
	Next() Pair                  // 获取下一个 key-value 对, 若为 nil, 则为单链表尾.
	SetNext(nextPair Pair) error // 设置下一个 key-value 对, 用于构成单链表
}

/***************************************************************
                    接口实现:

***************************************************************/

// pair 代表 key-value 对
type pair struct {
	key     string
	hash    uint64 // 哈希值
	element unsafe.Pointer
	next    unsafe.Pointer
}

//
func newPair(key string, element interface{}) (Pair, error) {
	p := &pair{
		key:  key,
		hash: hash(key), // 自定义 hash()算法
	}

	//
	if element == nil {
		return nil, newIllegalParameterError("element is nil")

	}

	//
	p.element = unsafe.Pointer(&element)
	return p, nil

}

//

func (p *pair) Key() string {
	return p.key
}

//
func (p *pair) Hash() uint64 {
	return p.hash
}

//
func (p *pair) Element() interface{} {
	pointer := atomic.LoadPointer(&p.element)

	if pointer == nil {
		return nil
	}
	return *(*interface{})(pointer)

}

//
func (p *pair) SetElement(element interface{}) error {
	if element == nil {
		return newIllegalParameterError("element is nil")
	}

	atomic.StorePointer(&p.element, unsafe.Pointer(&element))
	return nil
}

//
func (p *pair) Next() Pair {

	pointer := atomic.LoadPointer(&p.next)

	if pointer == nil {
		return nil
	}
	return (*pair)(pointer)
}

//
func (p *pair) SetNext(nextPair Pair) error {
	if nextPair == nil {
		atomic.StorePointer(&p.next, nil)
		return nil
	}

	//
	pp, ok := nextPair.(*pair)

	if !ok {
		return newIllegalPairTypeError(nextPair)

	}

	//
	atomic.StorePointer(&p.next, unsafe.Pointer(pp))
	return nil
}

//
func (p *pair) Copy() Pair {
	pCopy, _ := newPair(p.Key(), p.Element())
	return pCopy
}

//
func (p *pair) String() string {
	return p.genString(false)
}

// 生成 key-value 对的字符串表示形式
func (p *pair) genString(nextDetail bool) string {
	var buf bytes.Buffer

	//
	buf.WriteString("pair{key:")
	buf.WriteString(p.Key())
	buf.WriteString(", hash:")
	buf.WriteString(fmt.Sprintf("%d", p.Hash()))
	buf.WriteString(", element:")
	buf.WriteString(fmt.Sprintf("%+v", p.Element()))

	//
	if nextDetail {
		buf.WriteString(", next:")
		//
		if next := p.Next(); next != nil {
			if npp, ok := next.(*pair); ok {
				buf.WriteString(npp.genString(nextDetail)) //
			} else {
				buf.WriteString("<ignore>")
			}
		}
	} else {
		buf.WriteString("nextKey:")
		if next := p.Next(); next != nil {
			buf.WriteString(next.Key())
		}
	}

	//
	buf.WriteString("}")

	return buf.String()

}
