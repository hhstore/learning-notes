# -*- coding:utf8 -*-
__author__ = 'hhstore'

'''
功能: 单元测试 - 自定义异常

'''


import unittest


class BadAssError(TypeError):    # 自定义异常
    pass


def foo():      # 待测试函数
    raise BadAssError("bad ass error.")


class Test(unittest.TestCase):    # 测试
    def test_foo(self):
        self.assertRaises(BadAssError, foo)   # 抛出异常
        self.assertRaises(TypeError, foo)     # 抛出异常
        self.assertRaises(Exception, foo)     # 抛出异常


if __name__ == '__main__':
    unittest.main()

