#!/usr/bin/env python
# -*- coding: utf-8 -*-

'''
功能: 进程池 测试

'''


from multiprocessing import Pool
import time

def worker():
    print "work ..."
    print "current time: {}-{}-{} {}:{}:{}:{}".format(*time.localtime())

    print time.time()
    time.sleep(5)    # 5秒
    print "after sleep : {}-{}-{} {}:{}:{}:{}\n".format(*time.localtime())

def main():
    pool = Pool(processes=4)    # 进程池
    for i in range(4):
        pool.apply_async(worker)

    pool.close()
    pool.join()

if __name__ == '__main__':
    main()
