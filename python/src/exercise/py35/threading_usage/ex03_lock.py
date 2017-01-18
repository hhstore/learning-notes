#!/usr/bin/env python
# -*- coding: utf-8 -*-

import time
import threading

"""
线程同步机制:
    - Semaphore（信号量）
    - Lock（互斥锁）, 相当于: 信号量为1
    - RLock（可重入锁）:
        - acquire() 能够不被阻塞的被同一个线程调用多次。
        - 但release()要调用与acquire()相同的次数才能释放锁。
    - Condition（条件):
        - 一个线程等待特定条件，而另一个线程发出特定条件满足的信号。
        - 最好示例:「生产者/消费者」模型

Lock 方式, 示例代码:
    - 不加锁示例
    - 加锁示例

"""

_value1, _value2 = 0, 0


def get_no_lock():
    global _value1

    new_v = _value1 + 1
    time.sleep(0.001)  # 使用sleep让线程有机会切换
    _value1 = new_v


_lock = threading.Lock()  # 加锁


def get_lock():
    global _value2

    with _lock:  # 线程锁
        new_v = _value2 + 1
        time.sleep(0.001)
        _value2 = new_v


def run_by_threads(target_func):
    threads = []

    for i in range(100):
        thr = threading.Thread(target=target_func)
        thr.start()
        threads.append(thr)

    for thr in threads:
        thr.join()

    print("[value1 = {}, value2 = {}]".format(_value1, _value2))


if __name__ == '__main__':
    run_by_threads(get_no_lock)  # 结果随机, 但远小于100
    run_by_threads(get_lock)     # 加锁运行, 结果为100
