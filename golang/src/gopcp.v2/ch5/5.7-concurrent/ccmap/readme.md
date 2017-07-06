

# 并发安全的字典(Map)类型


## 源码说明:

- 阅读顺序: 自底向上, 逐层阅读实现源码
- 注意互斥锁和原子操作的大量应用

```bash

# 自底向上4层结构:
    - Pair
    - Bucket
    - Segment
    - ConcurrentMap
 
- 整体结构:    
    - base.go   // 全局常量
    - errors.go // 自定义错误
    - utils.go  // 自定义 hash()算法实现
    
    - ccmap.go  // 核心入口
      - segment.go  // 先看此模块的实现
        - pair.go
        - bucket.go
         - redistributor.go // 01.
      - redistributor.go    // 
        - bucket.go     

 
```

## 性能测试报告:


###  插入对比:

```bash

======================================================
BenchmarkCmapPutAbsent:
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

10000000	       189 ns/op
10000000	       186 ns/op
10000000	       183 ns/op
10000000	       188 ns/op
10000000	       184 ns/op
10000000	       190 ns/op
10000000	       192 ns/op
10000000	       191 ns/op
10000000	       185 ns/op
10000000	       195 ns/op
10000000	       190 ns/op
10000000	       189 ns/op
10000000	       186 ns/op
10000000	       189 ns/op
10000000	       186 ns/op
10000000	       189 ns/op
10000000	       186 ns/op
10000000	       186 ns/op
10000000	       201 ns/op
10000000	       190 ns/op
PASS

Process finished with exit code 0

======================================================
BenchmarkCmapPutPresent:
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

5000000	       228 ns/op
5000000	       238 ns/op
10000000	       247 ns/op
5000000	       229 ns/op
10000000	       224 ns/op
10000000	       240 ns/op
10000000	       255 ns/op
10000000	       238 ns/op
10000000	       235 ns/op
10000000	       255 ns/op
10000000	       245 ns/op
10000000	       239 ns/op
10000000	       232 ns/op
10000000	       236 ns/op
10000000	       235 ns/op
5000000	       241 ns/op
10000000	       242 ns/op
10000000	       246 ns/op
10000000	       250 ns/op
5000000	       251 ns/op
PASS

Process finished with exit code 0

======================================================
BenchmarkMapPut
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

50000000	        22.4 ns/op
100000000	        23.1 ns/op
100000000	        23.2 ns/op
100000000	        23.8 ns/op
50000000	        24.7 ns/op
50000000	        25.5 ns/op
50000000	        26.6 ns/op
50000000	        27.7 ns/op
100000000	        23.9 ns/op
50000000	        25.9 ns/op
PASS

Process finished with exit code 0



======================================================
BenchmarkCmapGet
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

20000000	        72.9 ns/op
20000000	        61.1 ns/op
20000000	        62.8 ns/op
20000000	        63.1 ns/op
20000000	        61.3 ns/op
20000000	        68.6 ns/op
20000000	        68.7 ns/op
20000000	        61.5 ns/op
20000000	        61.5 ns/op
20000000	        60.7 ns/op
PASS

Process finished with exit code 0


======================================================
BenchmarkMapGet
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

100000000	        10.6 ns/op
100000000	        10.8 ns/op
200000000	         9.57 ns/op
100000000	        14.7 ns/op
100000000	        11.5 ns/op
100000000	        13.9 ns/op
100000000	        10.9 ns/op
100000000	        16.0 ns/op
100000000	        10.8 ns/op
100000000	        11.5 ns/op
PASS

Process finished with exit code 0


======================================================
BenchmarkCmapDelete
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

5000000	       287 ns/op
30000000	        61.7 ns/op
20000000	        60.9 ns/op
20000000	       115 ns/op
20000000	        60.5 ns/op
20000000	       115 ns/op
20000000	       112 ns/op
10000000	       117 ns/op
20000000	       110 ns/op
30000000	        60.9 ns/op
20000000	        61.6 ns/op
20000000	       109 ns/op
20000000	        60.9 ns/op
30000000	        60.7 ns/op
20000000	        61.0 ns/op
20000000	       109 ns/op
30000000	        60.2 ns/op
30000000	        61.3 ns/op
20000000	        61.9 ns/op
20000000	        64.3 ns/op
PASS

Process finished with exit code 0

======================================================
BenchmarkMapDelete
======================================================
GOROOT=/usr/local/Cellar/go/1.8.3/libexec

100000000	        15.7 ns/op
100000000	        15.4 ns/op
100000000	        15.2 ns/op
100000000	        15.2 ns/op
100000000	        14.7 ns/op
100000000	        21.4 ns/op
100000000	        20.9 ns/op
100000000	        14.5 ns/op
100000000	        14.5 ns/op
100000000	        14.6 ns/op
100000000	        14.6 ns/op
100000000	        14.5 ns/op
100000000	        20.9 ns/op
100000000	        14.5 ns/op
100000000	        14.5 ns/op
100000000	        14.5 ns/op
100000000	        14.6 ns/op
100000000	        14.5 ns/op
100000000	        14.4 ns/op
100000000	        14.5 ns/op
PASS

Process finished with exit code 0



```
