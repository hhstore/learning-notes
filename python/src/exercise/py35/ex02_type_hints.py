#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
python3.5+ 新特性, 静态语义分析

"""

def hello(s: int) -> None:
    print(s)
    return None


# 变量类型,静态分析:
def hello2(s: str = "233") -> str:
    print(s)
    return s


if __name__ == '__main__':
    hello(222)
    hello("22")  # 提示类型不对

    hello2()
    hello2("xiaoming")
    hello2(2333)  # 提示类型不对
