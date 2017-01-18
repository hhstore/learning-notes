# -*- coding:utf8 -*-

'''
功能: XML-RPC 服务端
    1. 自动重启RPC服务.
    2. 访问权限管理.

依赖: SimpleXMLRPCServer

说明:
    1. 命令行,先运行服务器端,再运行客户端.
    2. 使用socket.setsockopt() 实现,设置端口复用,快速重启服务.

'''


__author__ = 'hhstore'



from SimpleXMLRPCServer import SimpleXMLRPCServer as BaseServer

class MyRPCServer(BaseServer):
    allow_client_hosts = ("127.0.0.1", "192.168.0.15")

    def __init__(self, host, port):
        BaseServer.__init__(self, (host, port))   # 接口参数格式转换,方便传参.

    def server_bind(self):              # 重写接口,实现服务自动重启
        import socket
        self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)   # 设置端口复用.
        BaseServer.server_bind(self)    # 调用父类方法.

    def verify_request(self, request, client_address):    # 重写 访问权限控制
        return client_address[0] in self.allow_client_hosts

def rpc_test_service():
    print "rpc_test_serive() is calling..."
    return True



def main():
    server = MyRPCServer("localhost", 5000)              # 创建RPC服务.在指定端口,监听请求.
    server.register_function(rpc_test_service)     # 注册函数
    print "server on..."
    server.serve_forever()


if __name__ == '__main__':
    main()
