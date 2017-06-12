

# 并发编程方式:

- `多进程`编程
- `多线程`编程
- `goroutine` 编程 (go方式: goroutine + channel)



## 方式1: 多进程编程:

### IPC通信(进程间通信)方式:

- 基于管道:
    - 匿名管道: [apipe](./apipe.go)  | P63
    - 命名管道: [npipe](./npipe)     | P65
    - 内存管道: [npipe](./npipe)     | P65
    - 重构版: [基于 reflect 实现函数多态](./npipe/npipe-v2.go)
    - 重构版: [基于 interface 实现函数多态](./npipe/npipe-v4.go)
- 基于信号:
    - 信号: [signal](./signal/signal.go) | P72
- 基于socket:
    - TCP socket: [socket](./socket/tcp_socket.go) | P88
    
## 方式2: 多线程编程:

- POSIX线程:
- 线程模型:
    - 用户级线程: 多对一(M:1)
    - 内核级线程: 一对一(1:1)
    - 两级线程模型(go采用): 多对多(M:N)
    
### 线程同步:
- 原子操作:
- 临界区保护:
- 互斥量: mutex | P108
- 条件变量:
    - 生产者/消费者模型
    - 条件变量+互斥量, 组合使用
    
### 并发编程性能优化:
- 编写可重入函数, 消除同步方法
- 控制临界区的纯度+粒度+执行耗时
- 避免长时间持有互斥量
- 优先使用原子操作, 而非互斥量

## 方式3: goroutine 编程:

### goroutine 本质:

- 不要用 `共享内存` 的方式来 `通信`
- 用 `通信` 的方式, 来 `共享内存`
- go 不推荐 `共享内存` 方式传递`数据`, 推荐使用 `channel(通道)`





