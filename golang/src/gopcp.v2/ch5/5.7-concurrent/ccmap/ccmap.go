package ccmap

import (
	"math"
	"sync/atomic"
)

/***************************************************************
                    核心模块 - 接口声明: 并发安全的字典

***************************************************************/

// 接口声明: 并发安全的字典
type ConcurrentMap interface {
	Concurrency() int                                  // 获取并发量
	Put(key string, element interface{}) (bool, error) // 插入
	Get(key string) interface{}                        // 查询
	Delete(key string) bool                            // 删除
	Len() uint64                                       // 获取字典中键值对数量
}

/***************************************************************
                    核心模块 - 接口实现: 并发安全的字典

 说明:
	- 依赖: segment.go
	- 通过多个分段, 利用分段锁保护共享资源, 降低互斥锁的开销
	- 在散列段中, 对 key-value 对, 作负载均衡

***************************************************************/

// 接口实现: 并发安全的字典
type aConcurrentMap struct {
	concurrency int
	segments    []Segment // 并发安全的散列段 [依赖: segment.go]
	total       uint64
}

//
func NewConcurrentMap(concurrency int,
	pairRedistributor PairRedistributor) (ConcurrentMap, error) {

	// 参数合法性校验:
	if concurrency <= 0 { // 过小
		return nil, newIllegalParameterError("concurrency is too small")
	}

	if concurrency > MAX_CONCURRENCY { // 过大
		return nil, newIllegalParameterError("concurrency is too large")
	}

	// 创建:
	cm := &aConcurrentMap{}
	cm.concurrency = concurrency
	cm.segments = make([]Segment, concurrency)

	for i := 0; i < concurrency; i++ { // 初始化
		cm.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)
	}

	return cm, nil

}

//
func (cm *aConcurrentMap) Concurrency() int {
	return cm.concurrency
}

//-----------------------------------------------------------------
// 查找指定参数的散列段
func (cm *aConcurrentMap) findSegment(keyHash uint64) Segment {
	if cm.concurrency == 1 {
		return cm.segments[0]
	}

	var keyHash32 uint32

	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}

	return cm.segments[int(keyHash32>>16)%(cm.concurrency-1)] // 匹配散列段

}

//-----------------------------------------------------------------
// 查找:
func (cm *aConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cm.findSegment(keyHash) // 匹配散列段
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}

	return pair.Element()
}

// 插入:
// 	- 过程:
//		- 新建一个 key-value 对
//		- 搜索 散列段, 执行插入
//	- 原子操作, 修改计数
func (cm *aConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element) // 新建一个键值对
	if err != nil {
		return false, err
	}

	s := cm.findSegment(p.Hash()) // 匹配散列段
	ok, err := s.Put(p)           // 插入
	if ok {
		atomic.AddUint64(&cm.total, 1) // 插入成功, 计数加1
	}
	return ok, err
}

// 删除:
//	- 原子操作, 修改计数
func (cm *aConcurrentMap) Delete(key string) bool {
	s := cm.findSegment(hash(key)) // 匹配散列段

	if s.Delete(key) {
		atomic.AddUint64(&cm.total, ^uint64(0)) // 删除成功, 计数减1
		return true
	}
	return false
}

// 获取字典长度:
//	- 原子操作
func (cm *aConcurrentMap) Len() uint64 {
	return atomic.LoadUint64(&cm.total)
}
