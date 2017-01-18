#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
- py27 练习, 涉及内容:
    - id()
    - 格式化字符串

"""


l1 = [1, 2, 3, 4, 5]
t1 = (4, 5, 6, 7, 8)

for item in l1:
    print "id(%s) = 0x%X // %d" % (item, id(item), id(item))

l1.append(6)
print "==========="
for item in l1:
    print "id(%s) = 0x%X // %d" % (item, id(item), id(item))

print "==========="

for item in t1:
    print "id(%s) = 0x%X // %d" % (item, id(item), id(item))

print "==========="
print "id(tuple)= 0x%X // %d" % (id(t1), id(t1))
print "id(list)= 0x%X // %d" % (id(l1), id(l1))

print id("aa"), id("bb"), id("c"), id("cc")


class test(object):
    def __init__(self):
        self.a = 20

    def get(self):
        return self

t = test()
print t.get()
