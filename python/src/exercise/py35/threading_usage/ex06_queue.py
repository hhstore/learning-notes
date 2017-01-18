#!/usr/bin/env python
# -*- coding: utf-8 -*-

import time
import threading

import random
import queue

"""

6. Queue

队列在并发开发中最常用的。

我们借助「生产者/消费者」模式来理解：
生产者把生产的「消息」放入队列，消费者从这个队列中对去对应的消息执行。

大家主要关心如下4个方法就好了：

put: 向队列中添加一个项。
get: 从队列中删除并返回一个项。
task_done: 当某一项任务完成时调用。
join: 阻塞直到所有的项目都被处理完。

"""

q = queue.Queue()


# 待计算的任务:
def tiny_task(n):
    return n * 2


# 生产者:
#   - 任务分发
#   - 队列处理
def producer():
    thr = threading.currentThread()
    print("[{}] producer is called.".format(thr.name))
    for i in range(1, 10):
        x = random.randint(2, 20)
        time.sleep(0.2)  # 设置短, 设置长, 观察运行差异 [0.2, 1]
        item = (tiny_task, x)
        q.put(item)
        print("[{}]: add {} to queue.".format(thr.name, x))
        # q.join()     # 加上此语句, 会等待 消费者处理完, 才产生新的数据.


#
# 消费者:
#   - 处理生产者分配的任务
#   - 注意添加 time.sleep(), 防止处理太快, 接收不到数据
#
def consumer():
    thr = threading.currentThread()
    print("[{}] consumer is called.".format(thr.name))
    fail_count = 0

    def handle_task(thr_name, task, arg):
        print("\t==>[{}] handle the task. [arg={}, result={}]".format(
            thr_name, arg, task(arg))
        )

    # time.sleep() 需要加暂停, 否则处理太快, 接收不到数据.
    while True:

        if fail_count >= 5:
            break
        if q.qsize() == 0:
            fail_count += 1
            time.sleep(2)
            print("\tfail count: [{}]".format(fail_count))
        else:
            print("\t[before] queue size={}".format(q.qsize()))
            fail_count = fail_count - 1 if fail_count > 0 else 0
            task, arg = q.get()
            handle_task(thr.name, task, arg)
            q.task_done()
            print("\t[after] queue size={}\n".format(q.qsize()))


def run():
    threads = []

    for func in (producer, consumer):
        thr = threading.Thread(target=func)
        thr.start()
        threads.append(thr)
    for thr in threads:
        print("\t->[{}] is done.".format(thr.name))
        thr.join()

    print("run over.")


if __name__ == '__main__':
    run()
