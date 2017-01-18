# -*- coding:utf8 -*-

'''
XML-RPC 方法调用- 客户端测试


'''


from xmlrpclib import Server

url = "http://www.oreillynet.com/meerkat/xml-rpc/server.php"

if __name__ == '__main__':
    server = Server(url)
    print server.meerkat.getItems(
        {"serach": "[Pp]ython",
         "num_items": 5,
         "descriptions": 0
         })
