package ccmap

import (
	"bytes"
	"sync"
	"sync/atomic"
)

/***************************************************************
                    接口声明: 并发安全的散列桶

***************************************************************/
type Bucket interface {
	Put(p Pair, lock sync.Locker) (bool, error) // 插入一个键-元素对
	Get(key string) Pair                        // 获取指定键的键-元素对
	GetFirstPair() Pair                         // 获取首个键-元素对
	Delete(key string, lock sync.Locker) bool   // 删除指定键-元素对.
	Clear(lock sync.Locker)                     // 清空当前散列桶
	Size() uint64                               // 返回当前散列桶尺寸
	String() string                             // 获取当前散列桶字符串表示
}

/***************************************************************
                    接口实现: 并发安全的散列桶

说明:
	- CRUD 操作的实现:
		- 注意 Put() 和 Delete() 稍复杂, 都是基于链表头插法实现

***************************************************************/
type aBucket struct {
	firstValue atomic.Value // key-value 对的列表表头
	size       uint64
}

//
// 占位符
//	- 原子值不可存 nil, 当散列桶为空时, 用此占位.
//
var placeholder Pair = &aPair{}

//
func newBucket() Bucket {
	b := &aBucket{}
	b.firstValue.Store(placeholder) // 初始化零值
	return b
}

// 获取首个 key-value 对
func (b *aBucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil { // bucket 为空, v 不存在
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder { // bucket 非空, v存在, 但值为零值
		return nil
	} else {
		return p //  bucket 非空, v存在, 且为非零值
	}
}

// 查找:
// 	- 指定 key, 返回其 value
//
func (b *aBucket) Get(key string) Pair {

	firstPair := b.GetFirstPair() // 链表头
	if firstPair == nil {         // 链表为空时
		return nil
	}

	// 链表非空时, 迭代搜索:
	//	- 搜索到, 则返回
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v // 搜索成功
		}
	}

	return nil
}

// 插入:
//	- 头插法, 插入到链表头部
// 	- 实现依赖:
//		- GetFirstPair()
//
func (b *aBucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}

	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	//
	firstPair := b.GetFirstPair() // 取首个键值对
	if firstPair == nil {         // 若首个键值对为空
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}

	// 若不为空
	var target Pair // 暂存变量, 用于搜索匹配时, 暂存
	key := p.Key()  // 待插入的键值对的 key 值

	for v := firstPair; v != nil; v = v.Next() { // 迭代搜索 key
		if v.Key() == key { // key 存在
			target = v // 匹配, 说明字典中, key 已存在, 需判断值是否为零值
			break
		}
	}

	// 若 target 为非零值:
	// 	- 说明上一步 for 迭代, 找到同名的 key.
	// 	- 说明原字典中, 已存在同名 key, 则插入失败
	//
	// 若原字典中, 同名 key 的 value 值非零, 则插入失败:
	if target != nil { // 判断原字典中同名 key 的 value 值, 是否为零值
		target.SetElement(p.Element()) // 若为非零值, 返回插入失败.(此步操作无意义)
		return false, nil
	}

	// 若 target 为nil:
	//	- 说明 原链表中, 无同名 key
	//	- 运行执行插入操作, 且插入到链表头
	//
	p.SetNext(firstPair)         // 将 p 插入到链表头
	b.firstValue.Store(p)        // 更新链表头为新插入的 key-value 对
	atomic.AddUint64(&b.size, 1) // 链表长度+1

	return true, nil
}

// 删除:
//	- 指定 key, 删除其键值对
//
func (b *aBucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	//
	firstPair := b.GetFirstPair() // 链表头
	if firstPair == nil {
		return false // 为空
	}

	//
	var prevPairs []Pair
	var target Pair
	var breakpoint Pair

	// 链表非空时:
	// 	-由链表头开始遍历, 搜索指定 key 对应的键值对
	//
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key { // 找到指定 key 的键值对
			target = v            // 暂存目标键值对
			breakpoint = v.Next() // 标记指定键值对的下一位
			break
		}

		// 未匹配 key 时, 把 键对应的值, 追加到暂存切片
		prevPairs = append(prevPairs, v) // 追加, 暂存
	}

	// 搜索失败:
	// 	- target 为空, 直接返回
	if target == nil {
		return false
	}

	// 搜索成功:
	//	- 因为是头插法, 以此断点处, 为头.
	//	- 所有处于被删除 key 之前的元素, 都重新用头插法插入
	//
	newFirstPair := breakpoint // 新的链表头

	// 头插法:
	//	- 逐个插入到链表头
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Copy()
		pairCopy.SetNext(newFirstPair) // 插入在链表头
		newFirstPair = pairCopy        // 更新链表头的指向
	}

	// 对链表头作判断:
	if newFirstPair != nil { // 若非零值, 则 bucket 指向该头
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder) // 若为零值, 则 bucket 存入零值
	}

	//
	atomic.AddUint64(&b.size, ^uint64(0)) // todo: ?? [值为: 0xffffffffffffffff]
	return true

}

// 清空:
func (b *aBucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	//
	atomic.StoreUint64(&b.size, 0)  // 长度清0
	b.firstValue.Store(placeholder) // 存入零值
}

// 长度:
func (b *aBucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

// 字符串表示:
func (b *aBucket) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")

	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}

	buf.WriteString("]")
	return buf.String()
}
