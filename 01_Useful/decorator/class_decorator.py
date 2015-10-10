# -*- coding:utf-8 -*-
__author__ = 'hhstore'

class MyDecorator(object):
    def __init__(self, fn):
        self.fn = fn
        print "{}.__init__() is calling...".format(self.__class__.__name__)

    def __call__(self):
        self.fn()
        print "{}.__call__() is calling...".format(self.__class__.__name__)


#############################################


@MyDecorator
def test_foo():
    print "inside aFunction()"

#############################################



if __name__ == '__main__':
    test_foo()


