# -*- coding:utf8 -*-

'''
功能: XML-RPC 服务端
    1. 实现允许RPC服务, 被远程结束掉.

依赖: SimpleXMLRPCServer

说明:
    1. 命令行,先运行服务器端,再运行客户端.
    2. 服务执行一次,就自动关闭.需要手动重启.
'''


__author__ = 'hhstore'



from SimpleXMLRPCServer import SimpleXMLRPCServer

running = True   # 全局运行状态

def rpc_test_service():
    global running
    running = False    # 修改运行状态
    return "rpc_test_service() is calling..."   # 必须有返回值


def main():
    addr = ("localhost", 5000)                     # 主机名, 端口
    server = SimpleXMLRPCServer(addr)              # 创建RPC服务.在指定端口,监听请求.
    server.register_function(rpc_test_service)     # 注册函数

    while running:     # 自主管理服务,是否允许客户端结束.(非死循环)
        print "server on..."
        server.handle_request()   # 处理RPC服务请求
    else:
        print "server stop..."

if __name__ == '__main__':
    main()
