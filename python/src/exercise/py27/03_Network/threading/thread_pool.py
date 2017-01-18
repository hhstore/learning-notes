# -*- coding: utf-8 -*-

__author__ = 'hhstore'

'''
功能: Python的线程池.
version: python3

ref_blog:
    http://www.open-open.com/home/space-5679-do-blog-id-3247.html
    http://www.the5fire.com/python-thread-pool.html

说明:
    1. 对并发的理解.
    2. 执行结果是交替的,很明显表明同时启动 N个 线程,在处理...当有些处理一半,另一些处理完,自然会打印输出信息.导致输出结果是无序的.
    3. 线程越多.越快.
    4. 应用场景区分: IO密集型, VS 计算密集型.
    5. 多线程,适合用于IO密集型的场景.

'''

import threading
import time
try:
    import Queue as queue    # python2.7
except ImportError:
    import queue             # python3+



# 管理者: 管理线程池和任务队列
class WorkManager(object):
    def __init__(self, work_num=1000, thread_num=4):
        self.work_queue = queue.Queue()
        self.threads = []

        self.__init_work_queue(work_num)       # 初始化
        self.__init_thread_pool(thread_num)    # 初始化

    def __init_thread_pool(self, thread_num):    # 初始化线程池
        for i in range(thread_num):
            work = Work(self.work_queue)         # 工程线程
            self.threads.append(work)            # 线程池更新

    def __init_work_queue(self, work_num):       # 初始化工作队列
        for i in range(work_num):
            self.add_task(todo_real_task, i)     # 添加具体的操作任务

    def add_task(self, func, *args):   # 单个任务添加至工作队列
        self.work_queue.put((func, list(args)))

    def check_queue(self):            # 队列剩余任务检查
        return self.work_queue.qsize()

    def wait_all_commplete(self):     # 等待所有线程,运行完毕
        for item in self.threads:
            if item.isAlive():
                item.join()           # 等待执行线程结束.


class Work(threading.Thread):
    def __init__(self, work_queue):
        super(Work, self).__init__()
        self.work_queue = work_queue
        self.start()

    def run(self):
        while True:
            try:
                do, args = self.work_queue.get(block=False)
                do(args)

                self.work_queue.task_done()    # 通知系统,任务已完成
            except Exception as e:
                print(str(e))    # 出错信息
                break

########################################################


# 具体的待处理任务
# IO密集型 VS 计算密集型
#
def todo_real_task(args):
    # print("args = {}\t\t{}".format(args, threading.current_thread()))
    time.sleep(0.1)    # 秒
    # print("\t\targs = {}\t\t{}".format(list(args), threading.current_thread()))


########################################################


# 并发性能测试. 统计
def profile_performance(task_num, thread_num):
    start = time.time()
    work_manager = WorkManager(task_num, thread_num)
    work_manager.wait_all_commplete()
    end = time.time()
    return {"task_num": task_num,
            "thread_num": thread_num,
            "cost_time": end-start}



########################################################
# 并发性能区间对比测试:
# 测试机: rMBP-i7-4核,16G
# 测试数据: 对比耗时
#
# 测试示例对比:
#
# Profile Result:
# {'thread_num': 0, 'task_num': 100, 'cost_time': 0.00037097930908203125}
# {'thread_num': 1, 'task_num': 100, 'cost_time': 10.415114879608154}
# {'thread_num': 2, 'task_num': 100, 'cost_time': 5.214631080627441}
# {'thread_num': 4, 'task_num': 100, 'cost_time': 2.598301887512207}
# {'thread_num': 8, 'task_num': 100, 'cost_time': 1.3543472290039062}
# {'thread_num': 16, 'task_num': 100, 'cost_time': 0.7277238368988037}
# {'thread_num': 32, 'task_num': 100, 'cost_time': 0.4187939167022705}
#
########################################################
def profile_total(task_num, thread_max):
    profile_result = [profile_performance(task_num, int(2**(i-1))) for i in range(thread_max+2)]
    print("Profile Result:")
    for item in profile_result:
        print(item)


def main():
    print("test start...")
    profile_total(100, 5)



if __name__ == '__main__':
    main()
