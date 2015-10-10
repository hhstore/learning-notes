# -*- coding:utf8 -*-
__author__ = 'hhstore'

'''
功能: 解决文件读写的编码问题

依赖: codecs模块

说明:
    1. 通过codecs模块,解决文件读写的编码问题.
    2. 如下是标准写法.
    3. 对文件读,作了类型判断.
    4. 文件写,简单判断.
'''

import codecs
import os


class BetterIO(object):
    def __init__(self):
        self.data = None

    def read(self, infile):
        if os.path.isfile(infile):    # 有效的文件
            with codecs.open(infile, mode="r", encoding="utf8") as f:
                self.data = f.read()
        elif isinstance(infile, basestring):     # 字符串
            self.data = infile

        print "file content:\n", self.data

    def write(self, outfile):
        if self.data:
            with codecs.open(outfile, mode="wb", encoding="utf8") as f:
                f.write(self.data)


def test():
    data = u"小米, 我是中文,我是中文..."
    filename = "cn.txt"

    writer = BetterIO()
    writer.read(data)
    writer.write(filename)


if __name__ == '__main__':
    test()
