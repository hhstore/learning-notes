#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
- py27 练习, 涉及内容:
    - globals(), 理解 global 内容
    - 文件保存


"""


def print_globals():
    data = globals()
    print "type(globals()): ", type(data)

    save_file("result.txt", data)


def save_file(filename, data):
    with open(filename, "w") as f:
        for k, v in data.iteritems():
            result = "{}: {}\n".format(k, v)
            print result,
            f.write(result)


if __name__ == '__main__':
    print_globals()

"""

type(globals()):  <type 'dict'>
__builtins__: <module '__builtin__' (built-in)>
save_file: <function save_file at 0x1006ff758>
__file__: /Users/hhstore/iGit/iSpace/iPy/2015/06/test_globals.py
watch_globals: <function watch_globals at 0x1006ff6e0>
__name__: __main__
__package__: None
__doc__: None


"""
