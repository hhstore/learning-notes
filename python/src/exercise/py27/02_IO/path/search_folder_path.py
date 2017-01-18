# -*- coding: utf-8 -*-
'''

功能: 在当前代码执行路径树中,寻找并获取某个folder的路径.
说明:
    1. 递归搜索.
    2. 递归至 "/" (Unix根目录)下,仍未匹配,则结束.未对windows平台作判断.需注意.
    3. 可根据需要指定搜索路径.
'''

import os

class FolderSearcher(object):
    def __init__(self, folder_name, match_dir=None):
        self.match_dir = match_dir if match_dir else os.getcwd()    # 当前路径
        print "match_dir: ", self.match_dir
        self.folder_name = folder_name

    def search(self):              # 搜索匹配
        def helper(current_dir):   # 辅助子函数, 递归调用
            pardir, subdir = os.path.split(current_dir)

            if subdir == self.folder_name:                         # 已匹配
                return os.path.join(pardir, self.folder_name)      # 返回 文件夹所处在的路径
            # 递归调用
            return helper(pardir) if not pardir == "/" else "Not match..."   # Unix根目录.未对其他平台判断.必要判断.
        return helper(self.match_dir)


def test():
    folder_name = "iPyScript"
    searcher = FolderSearcher(folder_name)
    print "folder path: ", searcher.search()


if __name__ == '__main__':
    test()
