#!/usr/bin/env python
# -*- coding: utf-8 -*-
import time
from random import random
from threading import Thread, Semaphore


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

Semaphore 方式, 示例代码:
    - 信号量, 控制可访问资源数目

"""

# 同时能访问资源的数量为3
semaphore = Semaphore(3)


def foo(tid):
    with semaphore:
        print("{} acquire semaphore.".format(tid))
        wait_t = random() * 2
        time.sleep(wait_t)

    print("{} release semaphore.".format(tid))


def run_by_thread():
    threads = []

    for i in range(5):
        thr = Thread(target=foo, args=(i,))
        threads.append(thr)
        thr.start()

    for thr in threads:
        thr.join()


if __name__ == '__main__':
    run_by_thread()
