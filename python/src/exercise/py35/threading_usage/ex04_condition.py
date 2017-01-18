#!/usr/bin/env python
# -*- coding: utf-8 -*-

import time
import threading

"""

Condition（条件）示例:
    - 生产者/消费者模型
    - 一个线程等待特定条件，而另一个线程发出特定条件满足的信号。

"""


# 消费者:
def consumer(cond):
    thr = threading.currentThread()

    with cond:
        # wait()方法创建了一个名为waiter的锁，并且设置锁的状态为locked
        # 这个waiter锁用于线程间的通讯
        cond.wait()

        print("{}: resource is availble to consumer.".format(thr.name))


# 生产者:
def producer(cond):
    thr = threading.currentThread()

    with cond:
        print("{}: making resource available.".format(thr.name))

        cond.notifyAll()   # 释放waiter锁，唤醒消费者


def run():
    condition = threading.Condition()

    thr_co1 = threading.Thread(name="thr_co1", target=consumer, args=(condition,))
    thr_co2 = threading.Thread(name="thr_co2", target=consumer, args=(condition,))

    thr_prod = threading.Thread(name="thr_prod", target=producer, args=(condition,))

    thr_co1.start()
    time.sleep(1)
    thr_co2.start()
    time.sleep(1)
    thr_prod.start()  # 生产者发送通知之后，消费者都收到了


if __name__ == '__main__':
    run()

