# -*- coding:utf-8 -*-

__author__ = 'hhstore'

'''
note :
1. 测试 yield 关键字用法.
2. 给出一些案例.对比执行效果.


 在Python中，拥有这种能力的“函数”被称为生成器，它非常的有用。
 生成器（以及yield语句）最初的引入是为了让程序员可以更简单的编写用来产生值的序列的代码。
 以前，要实现类似随机数生成器的东西，需要实现一个类或者一个模块，在生成数据的同时保持对每次调用之间状态的跟踪。
 引入生成器之后，这变得非常简单。

为了更好的理解生成器所解决的问题，让我们来看一个例子。
在了解这个例子的过程中，请始终记住我们需要解决的问题：生成值的序列。

注意：在Python之外，最简单的生成器应该是被称为协程（coroutines）的东西。
在本文中，我将使用这个术语。请记住，在Python的概念中，这里提到的协程就是生成器。
Python正式的术语是生成器；协程只是便于讨论，在语言层面并没有正式定义。




    generator是用来产生一系列值的
    yield则像是generator函数的返回结果
    yield唯一所做的另一件事就是保存一个generator函数的状态
    generator就是一个特殊类型的迭代器（iterator）
    和迭代器相似，我们可以通过使用next()来从generator中获取下一个值
    通过隐式地调用next()来忽略一些值



'''


def generator_value():
    yield 1
    yield 2
    yield 3

######################################################

def print_generator_value(gen):
    if gen:
        for value in gen:
            print value


def print_generator_next(gen):
    if gen:
        print 'next()= ', next(gen)
        print 'next()= ', next(gen)
        print 'next()= ', next(gen)


if __name__ == '__main__':
    gen1 = generator_value()
    gen2 = generator_value()

    print_generator_value(gen1)
    print_generator_next(gen2)   # 注意为何要定义2个gen变量的原因. 生成器的特点决定.


