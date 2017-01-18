#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os

"""
- py35 练习, 涉及内容:
    - os.walk() 遍历目录树

"""


# 目录遍历:
#   - 过滤特定文件
def path_traverse(root_dir, file_tail=".py"):
    for parent, dirnames, filenames in os.walk(root_dir):
        for filename in filenames:
            if filename.endswith(file_tail):
                f_head = filename.rsplit(".")[0]
                print('FILE', os.path.abspath(os.path.join(parent, filename)))
                print("Parent folder:", parent)
                print("Filename:", filename)
                print("Filename head:", f_head)
                print("\n")


if __name__ == "__main__":
    root_dir = os.getcwd()
    path_traverse(root_dir)
