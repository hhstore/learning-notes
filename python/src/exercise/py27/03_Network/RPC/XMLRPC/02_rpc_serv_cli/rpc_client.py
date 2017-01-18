# -*- coding:utf8 -*-

'''
功能: XML-RPC 客户端-实现
依赖: xmlrpclib
说明:
    1. 命令行,先运行服务器端,再运行客户端.
    2. 测试远程调用函数效果.

'''


__author__ = 'hhstore'


from xmlrpclib import Server, Fault

uri = "http://localhost:5000"    # 远程过程调用服务地址

def main():
    server = Server(uri)        # 连接远程服务.

    # 测试远程调用函数.
    print "chop_in_half(): ", server.chop_in_half("hello server.")
    print "repeat(): ", server.repeat("hello server.", 3)
    print "_string(): ", server._string("addadd")
    # print "join(): ", server.python_string.join(["what a", "fucking day."], "-")
    try:
        server._privateFunction()
    except Fault:
        print "_privateFunction(): not supported, can't access..."


if __name__ == '__main__':
    main()
