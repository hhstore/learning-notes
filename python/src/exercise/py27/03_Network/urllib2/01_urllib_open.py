# -*- coding:utf8 -*-

__author__ = 'hhstore'


'''
功能: 测试urllib2.

'''

import urllib2

# 测试主页
URL = "http://www.zhihu.com"

def get_url_page(url):
    request = urllib2.Request(url)           # 请求
    response = urllib2.urlopen(request)      # 打开
    result = response.read()                 # 读取内容
    return result

def save_html(data, filename="./data.html"):
    with open(filename, 'w') as f:
        f.write(data)

if __name__ == '__main__':
    page = get_url_page(URL)
    print "result = ", page

    # save_html(get_url_page(URL))   # 保存成HTML页面

