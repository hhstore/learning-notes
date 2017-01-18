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

Queue模块还自带了PriorityQueue（带有优先级）和LifoQueue（后进先出）2种特殊队列。
我们这里展示下线程安全的优先级队列的用法，
PriorityQueue要求我们put的数据的格式是(priority_number, data)，
我们看看下面的例子：

========================================================

Queue 对象方法:
    - put: 向队列中添加一个项。
    - get: 从队列中删除并返回一个项。
    - task_done: 当某一项任务完成时调用。
    - join: 阻塞直到所有的项目都被处理完。


- 改进版:
    - 使用队列(queue) 实现多线程, 并发处理:
        - 单 master, 多 worker 模式
        - 性能测试
        - 此示例代码, 很实用


"""

q_task = queue.Queue()  # 任务队列
q_result = queue.Queue()  # 结果队列


# 性能分析装饰器:
def profile(func):
    def wrapper(*args, **kwargs):
        import time
        start_at = time.time()
        func(*args, **kwargs)
        end_at = time.time()

        print("<Cost: {}>".format(end_at - start_at))

    return wrapper


#
# 多线程处理:
#   - 加速 IO 密集型任务计算
#
class MultiWorker(object):
    def __init__(self, task_func, task_num, slave_num):
        self.tasks = self.make_task(task_func, task_num)
        self.q_task = queue.Queue()
        self.q_result = queue.Queue()
        self.slave_num = slave_num

    #
    # 构造任务:
    #
    def make_task(self, task_func, task_num):
        tasks = [
            (i, task_func, random.randint(2, 20))
            for i in range(task_num)
            ]
        return tasks

    # 生产者(主线程):
    #   - 任务分发
    #   - 队列处理
    def master(self):
        for task in self.tasks:
            self.q_task.put(task)  # 任务创建
            # self.q_task.join()    # 阻塞, 会等待 消费者处理完, 才产生新的数据.

    #
    # 消费者(工作线程):
    #   - 处理生产者分配的任务
    #   - 注意添加 time.sleep(), 防止处理太快, 接收不到数据
    #
    def slave(self):
        def handle_task(task, arg):
            return task(arg)

        while True:
            task_no, task_func, task_arg = self.q_task.get()
            result = (
                task_no, task_arg,
                handle_task(task_func, task_arg)  # 处理任务
            )
            self.q_result.put(result)  # 保存处理结果, 共享数据集, 线程安全
            self.q_task.task_done()

            if self.q_task.empty():
                break

    #
    # 启动服务:
    #   - 多线程处理任务集
    #
    def run(self):
        threads = []

        # master:
        thr = threading.Thread(target=self.master)
        thr.start()
        threads.append(thr)

        # worker:
        for i in range(self.slave_num):
            thr = threading.Thread(target=self.slave, name="worker%s" % i)
            thr.start()
            threads.append(thr)

        for thr in threads:
            thr.join()

    def task_result(self):
        result = []
        while not self.q_result.empty():
            result.append(self.q_result.get())  # 取出计算结果
        return result


@profile
def run_no_threads(task_func, task_num):
    for i in range(task_num):
        x = random.randint(2, 20)
        task_func(x)


@profile
def run_by_threads(task_func, task_num, worker_num):
    w = MultiWorker(task_func, task_num, worker_num)
    w.run()

    # result = w.task_result()
    # print("result: size={}, value={}".format(len(result), result))  # 取出计算结果


#
# 待计算的任务:
# - 加 sleep, 模拟IO密集型任务
#
def tiny_task(n):
    # time.sleep(0.0001)   # 对比:[0.001, 0.0001, 0.0001]
    time.sleep(0.001)
    return n * 2


def main():
    task_num = 500    # 对比: [500, 1000, 1500, 5000]

    # 测试参考: [rmbp15/i7/4core/16g]
    run_no_threads(tiny_task, task_num)
    run_by_threads(tiny_task, task_num, 2)  # 2线程, 运行性能对比很明显
    run_by_threads(tiny_task, task_num, 4)  # 4线程
    run_by_threads(tiny_task, task_num, 8)  # 8线程
    run_by_threads(tiny_task, task_num, 16)  # 16线程
    run_by_threads(tiny_task, task_num, 32)  # 32线程



if __name__ == '__main__':
    main()


"""

<Cost: 1.329482078552246>
<Cost: 0.6775009632110596>
<Cost: 0.3442878723144531>
<Cost: 0.17603206634521484>
<Cost: 0.08569478988647461>


<Cost: 1.2579030990600586>
<Cost: 0.6645040512084961>
<Cost: 0.34133100509643555>
<Cost: 0.17490506172180176>
<Cost: 0.0863950252532959>

"""
