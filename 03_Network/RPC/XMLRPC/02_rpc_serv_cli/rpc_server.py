# -*- coding:utf8 -*-

'''
功能: XML-RPC 服务端实现
依赖: SimpleXMLRPCServer
说明:
    1. 命令行,先运行服务器端,再运行客户端.
    2. ctrl+c, 退出服务器端服务.
'''


__author__ = 'hhstore'



from SimpleXMLRPCServer import SimpleXMLRPCServer


# 自定义的类,待注册
class StringFunction(object):
    def __init__(self):
        import string
        self.python_string = string

    def _prinvateFunction(self):
        return "never get this result on the client."

    def chop_in_half(self, a_str):
        return a_str[:len(a_str)/2]

    def repeat(self, a_str, times):
        return a_str * times


def main():
    addr = ("localhost", 5000)           # 主机名, 端口
    server = SimpleXMLRPCServer(addr)    # 创建RPC服务.在指定端口,监听请求.
    server.register_instance(StringFunction())                             # 注册自定义的类
    server.register_function(lambda a_str: "_" + a_str, name="_string")    # 注册一个lambda函数,并命名为 _string()
    print "server on..."
    server.serve_forever()                                                 # 启动RPC服务,死循环.


if __name__ == '__main__':
    main()
