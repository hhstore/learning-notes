# -*- coding:utf-8  -*-


from server.base import app
from server.views.calculator_view import *

if __name__ == '__main__':
    app.run(debug=True, port=5000)   # 本机测试
    # app.run(debug=True, port=5000, host='0.0.0.0')  # 局域网测试.监听所有公开的 IP
