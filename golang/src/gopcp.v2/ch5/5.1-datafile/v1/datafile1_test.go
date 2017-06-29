package v1

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

//***************************************************************
//                         单元测试
//
//***************************************************************

func removeFile(path string) error {
	file, err := os.Open(path) // 文件打开
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	file.Close()           // 文件关闭
	return os.Remove(path) // 删除文件
}

//***************************************************************
//                         主运行函数:
//***************************************************************

func TestIDataFile(t *testing.T) {

	//启动测试:
	t.Run(
		"v1/all",
		func(t *testing.T) {
			//
			dataLen := uint32(3)

			// 文件路径1:
			path1 := filepath.Join(os.TempDir(), "data_file_test_new.txt")

			defer func() {
				if err := removeFile(path1); err != nil {
					t.Error("Open file error: %s\n", err)
				}
			}()

			// 启动测试:
			t.Run(
				"New",
				func(t *testing.T) {
					testNew(path1, dataLen, t)
				},
			)

			// 文件路径2:
			path2 := filepath.Join(os.TempDir(), "data_file_test.txt")
			defer func() {
				if err := removeFile(path2); err != nil {
					t.Fatalf("Open file error: %s\n", err)
				}
			}()

			max := 100000
			// 启动测试:
			t.Run("WriteAndRead", func(t *testing.T) {
				testRW(path2, dataLen, max, t)
			})

		},
	)

}

//***************************************************************

// 子测试:
func testNew(path string, dataLen uint32, t *testing.T) {
	t.Logf("New a data file (path: %s, dataLen: %d)...\n", path, dataLen)

	f, err := NewDataFile(path, dataLen)
	if err != nil {
		t.Logf("Coudn't new a data file: %s", err)
		t.FailNow()
	}

	if f == nil {
		t.Log("Unnormal data file!")
		t.FailNow()
	}

	defer f.Close()

	if f.DataLen() != dataLen {
		t.Fatalf("Incorrect data length!")
	}

}

// 子测试:
func testRW(path string, dataLen uint32, max int, t *testing.T) {
	t.Logf("New a data file (path: %s, dataLen: %d)...\n", path, dataLen)

	f, err := NewDataFile(path, dataLen)
	if err != nil {
		t.Logf("Couldn't new a data file: %s", err)
		t.FailNow()
	}
	defer f.Close()

	//
	var wg sync.WaitGroup

	wg.Add(5)

	//---------------------------------------------
	// 写入:
	//---------------------------------------------
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()

			var prevWSN int64 = -1

			//
			for j := 0; j < max; j++ {
				data := Data{
					byte(rand.Int31n(256)),
					byte(rand.Int31n(256)),
					byte(rand.Int31n(256)),
				}
				//
				wsn, err := f.Write(data)
				if err != nil {
					t.Fatalf("Unexpected writing error: %s\n", err)
				}

				if prevWSN >= 0 && wsn <= prevWSN {
					t.Fatalf("Incorect WSN %d! (lt %d)\n", wsn, prevWSN)
				}

				prevWSN = wsn

			}

		}()
	}

	//---------------------------------------------
	// 读取:
	//---------------------------------------------

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			var prevRSN int64 = -1

			//
			for j := 0; j < max; j++ {
				rsn, data, err := f.Read()
				if err != nil {
					t.Fatalf("Unexpected writing error: %s\n", err)

				}

				if data == nil {
					t.Fatalf("Unnormal data!")
				}

				if prevRSN >= 0 && rsn <= prevRSN {
					t.Fatalf("Incorect RSN %d! (lt %d)\n", rsn, prevRSN)
				}

				//
				prevRSN = rsn

			}
		}()
	}

	//
	wg.Wait()
}
