#!/usr/bin/env python
# -*- coding: utf-8 -*-


import threading

"""
普通版本测试:
    - 多线程, 并没有非多线程版本快, 受限于 GIL 影响.


"""


# 性能分析装饰器:
def profile(func):
    def wrapper(*args, **kwargs):
        import time
        start_at = time.time()
        func(*args, **kwargs)
        end_at = time.time()

        print("<Cost: {}>".format(end_at - start_at))

    return wrapper


# 递归测试:
def fib(n):
    if n <= 2:
        return 1
    return fib(n - 1) + fib(n - 2)


@profile
def run_no_thread():
    fib(35)
    fib(35)


@profile
def run_has_thread():
    for i in range(2):
        thr = threading.Thread(target=fib, args=(35,))
        thr.start()

    main_thread = threading.currentThread()

    for thr in threading.enumerate():
        if thr is main_thread:
            continue
        thr.join()


if __name__ == '__main__':
    run_no_thread()  # < Cost: 5.113997936248779 >
    run_has_thread()  # < Cost: 5.631266832351685 >
