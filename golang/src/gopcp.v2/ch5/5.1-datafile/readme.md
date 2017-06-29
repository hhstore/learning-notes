
# 说明:



## 单元测试写法:

### `Gogland IDE` 单元测试设置 `GOPATH`:

- 默认设置一个全局 `GOPATH` 环境变量.
- 当有单元测试时, 需针对特定项目目录, 设置一个独立的局部 `project` GOPATH 环境变量.
    - 本项目下, 需要设置 一个项目的 GOPATH: 
        - `~/learning-notes/golang/`
        - 保证项目目录为 `src/` 上一层目录即可.
    - 再启动 IDE 的`单元测试`功能, 正常运行 
    


## 同步方法对比:

- 互斥锁: `sync.Mutex`
- 读写锁: `sync.RWMutex`
- 条件变量: ``sync.Cond``
    - 搭配锁使用
- 原子操作: ``atomic.LoadInt64(), atomic.CompareAndSwapInt64()``
    - 优先选择使用原子操作, 开销小