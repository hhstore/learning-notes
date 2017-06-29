package concurrent

import (
	"math/rand"
	"sync"
	"testing"
)

//***************************************************************

// 测试运行:
func TestConcurerntArray(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		// 测试1:
		t.Run("New", testNew)

		//
		array := NewConcurerntArray(uint32(rand.Int31n(100)))
		maxI := uint32(1000)

		// 测试2:
		t.Run("Set", func(t *testing.T) {
			testSet(array, maxI, t)
		})

		// 测试3:
		t.Run("Get", func(t *testing.T) {
			testGet(array, maxI, t)
		})

	})
}

//***************************************************************

func testNew(t *testing.T) {
	expectedLen := uint32(rand.Int31n(1000))
	array := NewConcurerntArray(expectedLen)

	if array == nil {
		t.Fatalf("Unnormal int array!")
	}

	if array.Len() != expectedLen {
		t.Fatalf("Incorrect int array length!")
	}

}

//
func testSet(array ConcurrentArray, maxI uint32, t *testing.T) {
	length := array.Len()
	var wg sync.WaitGroup
	wg.Add(int(maxI))

	// 测试并发安全:
	for i := uint32(0); i < maxI; i++ {
		// 并发:
		go func(i uint32) {
			defer wg.Done()

			//
			for j := uint32(0); j < length; j++ {
				err := array.Set(j, int(j*i))
				if uint32(j) >= length && err == nil {
					t.Fatalf("Unexpected nil error! (index: %d)", j)
				} else {
					if err != nil {
						t.Fatalf("Unexpecetd error: %s (index: %d)", err, j)
					}
				}
			}
		}(i)
	}

	wg.Wait() // 阻塞等待

}

//
func testGet(array ConcurrentArray, maxI uint32, t *testing.T) {
	length := array.Len()
	max := int((maxI - 1) * (length - 1))

	//
	for i := uint32(0); i < length; i++ {
		e, err := array.Get(i) // 获取元素:

		if err != nil {
			t.Fatalf("Unexpected error: %s (index: %d)", err, i)
		}

		if e < 0 || e > max {
			t.Fatalf("Incorrect element: %d! (index: %d, expect max: %d", e, i, max)
		}

	}

}
