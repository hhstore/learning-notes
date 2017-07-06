package ccmap

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/***************************************************************
                    接口声明: 并发安全的散列段
说明:
	- 依赖: Pair
***************************************************************/

type Segment interface {
	Put(p Pair) (bool, error)                    //插入一个键-元素对
	Get(key string) Pair                         //获取一个键-元素对, 根据键计算哈希值
	GetWithHash(key string, keyHash uint64) Pair // 获取一个键-元素对
	Delete(key string) bool                      //删除指定键-元素对
	Size() uint64                                // 获取当前段尺寸
}

/***************************************************************
                    接口实现: 并发安全的散列段
说明:
	- 依赖: Pair
	- CRUD 操作: 通过互斥锁, 控制并发安全

***************************************************************/
type aSegment struct {
	buckets           []Bucket          // 散列桶切片
	bucketsLen        int               // 散列桶切片长度
	pairTotal         uint64            // 键-元素对总数
	pairRedistributor PairRedistributor //再分布器: 键-元素对的
	lock              sync.Mutex        // 互斥锁
}

// 实例创建:
func newSegment(bucketNumber int, redistributor PairRedistributor) Segment {
	if bucketNumber <= 0 { // 若传参不合法, 使用默认值
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if redistributor == nil { // 若传参为 nil, 创建默认的负载均衡器
		redistributor = newDefaultPairRedistributor(DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber)
	}

	buckets := make([]Bucket, bucketNumber) // 创建切片
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket() // 初始化
	}

	return &aSegment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: redistributor,
	}
}

// 查找:
//	- 通过 hash 值获取元素
//	- 互斥锁, 并发安全控制
//
func (s *aSegment) GetWithHash(key string, keyHash uint64) Pair {
	s.lock.Lock()                                     //加锁
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))] // 通过 hash 值, 获取元素
	s.lock.Unlock()                                   //解锁
	return b.Get(key)

}

// 查找:
//	- 通过 key 获取对于元素
//
func (s *aSegment) Get(key string) Pair {
	return s.GetWithHash(key, hash(key)) // 自定义 hash()方法, 在 utils.go 中
}

// 插入:
//	- 互斥锁, 并发安全控制
//
func (s *aSegment) Put(p Pair) (bool, error) {
	s.lock.Lock() // 加锁
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, 1) // 计数加1
		s.redistribute(newTotal, b.Size())            // 负载均衡
	}
	s.lock.Unlock() // 解锁
	return ok, err
}

// 删除:
//	- 互斥锁, 并发安全控制
//
func (s *aSegment) Delete(key string) bool {
	s.lock.Lock()
	b := s.buckets[int(hash(key)%uint64((s.bucketsLen)))]
	ok := b.Delete(key, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, ^uint64(0)) // todo: 计数修改
		s.redistribute(newTotal, b.Size())                     // 负载均衡
	}
	s.lock.Unlock() // 解锁
	return ok
}

// 查询容量:
//	- 原子操作, 并发安全
func (s *aSegment) Size() uint64 {
	return atomic.LoadUint64(&s.pairTotal) // 原子操作
}

// 散列桶的负载均衡:
//	- 检查给定参数, 并设置相应阀值和计数
//	- 对 key-value对 作负载均衡, 必要时重新分配
//	- 注意: 需要在互斥锁保护的前提下, 调用本方法.
//
func (s *aSegment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()

	//
	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)

	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := s.pairRedistributor.Redistribe(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketsLen = len(s.buckets)
	}
	return nil

}
