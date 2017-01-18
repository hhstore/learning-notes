#!/usr/bin/env python
# -*- coding: utf-8 -*-
import time
import threading
from random import randint

"""
5. Event

- 一个线程发送/传递事件，
- 另外的线程等待事件的触发。
- 同样用「生产者/消费者」模型举例:
    - 可以看到事件, 被2个消费者, 比较平均的接收并处理了。
    - 如果使用了wait方法，线程就会等待我们设置事件，这有助于保证任务完成。


- 处理过程:
    - 生产者产生数据:
        - 产生数据, 并发给消费者(append 到缓冲区)
    - 消费者处理数据:
        - 消费者监听事件, 不断轮询.
        - 接受到数据, 就处理.(pop 取出丢掉)


"""

TIMEOUT = 3


# 消费者:
def consumer(event, data):
    thr = threading.currentThread()
    fail_num = 0

    # 死循环
    #   - 连续5次接收不到数据, 就结束运行
    #
    while True:
        set_event = event.wait(TIMEOUT)

        if set_event:
            try:
                digit = data.pop()
                print("\t[{}]: receive {} , and handle.".format(thr.name, digit))
                time.sleep(2)  # 模拟处理的慢, 易于平均分配.
                event.clear()
            except IndexError:
                pass
        else:
            fail_num += 1
            print("\t[{}]: receive nothing... [{}]".format(thr.name, fail_num))

        if fail_num >= 5:
            print("[{}]: thread is done.".format(thr.name))
            break


# 生产者:
def producer(event, data):
    thr = threading.currentThread()

    for i in range(1, 20):
        digit = randint(10, 100)
        data.append(digit)
        print("[{} - {}] --> appended {} to list.".format(i, thr.name, digit))
        event.set()
        time.sleep(1)

    print("\n[{}]: thread is done.".format(thr.name))


def run():
    event = threading.Event()
    data = []
    threads = []

    # 消费者:
    for name in ("consumer1", "consumer2"):
        thr = threading.Thread(name=name, target=consumer, args=(event, data))
        thr.start()
        threads.append(thr)

    # 生产者:
    #   - 生产者和消费者, 要分开
    #   - 合并在一个 for 里, 会报错
    #
    p = threading.Thread(name='producer1', target=producer, args=(event, data))
    p.start()
    threads.append(p)

    for thr in threads:
        thr.join()

    print("run over.")


if __name__ == '__main__':
    run()  # 生产者产生的事件, 被2个消费者, 比较平均的接受和处理.
