#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
可调用对象:
    - 魔法方法: __call__(self, [args...])
    - 允许类的一个实例, 像函数那样被调用
    - 允许你自己类的对象, 表现得像是函数,然后你就可以“调用”它们,把它们传递到使用函数做参数的函数中
    - 本质上这代表了 x() 和 x.__call__() 是相同的
    - 注意 __call__ 可以有多个参数,可以像定义其他任何函数一样,定义 __call__ ,喜欢用多少参数就用多少.
    - __call__ 在某些需要经常改变状态的类的实例中显得特别有用.
    - “调用”这个实例来改变它的状态,是一种更加符合直觉,也更加优雅的方法.

参考:
    - http://pyzh.readthedocs.io/en/latest/python-magic-methods-guide.html#id20



"""


class A(object):
    """表示一个实体的类.
    调用它的实例, 可以更新实体的位置

    """

    def __init__(self, x, y):
        self.x, self.y = x, y
        print "<__init__() is called.> | x={}, y={}".format(self.x, self.y)

    def __call__(self, x, y):
        """改变实体的位置
        :param x:
        :param y:
        :return:
        """
        self.x, self.y = x, y
        print "<__call__() is called.> | x={}, y={}".format(self.x, self.y)


class B(object):
    """表示一个实体的类.
    调用它的实例, 可以更新实体的位置
    """

    def __init__(self, x, y):
        self.x, self.y = x, y
        print "<__init__() is called.> | x={}, y={}".format(self.x, self.y)

    def __call__(self, *args, **kwargs):
        self.x, self.y = self.y, self.x
        print "<__call__() is called.> | x={}, y={}".format(self.x, self.y)


if __name__ == '__main__':
    m = A(2, 3)
    m(4, 8)

    n = B(1, 9)
    n()
