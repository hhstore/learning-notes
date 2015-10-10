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
    print "rpc_test_service(): ", server.rpc_test_service()

if __name__ == '__main__':
    main()
