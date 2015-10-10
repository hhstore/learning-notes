# -*- coding:utf8 -*-
__author__ = 'hhstore'


'''
功能: urllib2测试
说明:
    1. 测试urllib2的各种方法.
    2. 获取HTTP头信息.
    3. 设置浏览器代理.

'''


import urllib2

# 测试主页
URL = 'http://www.douban.com'

# 浏览器代理
BROWSER_AGENT = 'Mozilla/5.0 (Windows NT 5.1; rv:20.0) Gecko/20100101 Firefox/20.0'


def get_url_page_use_agent(url):
    opener = urllib2.build_opener()
    opener.addheaders = [('User-agent', BROWSER_AGENT)]    # 设置浏览器代理
    result = opener.open(URL)
    return result

def get_url_page_data(url_data):    # 获取页面内容
    return url_data.read() if url_data else None

####################################################


def print_url_header(url_data):    # HTTP头信息
    if url_data:
        print "response header:"
        print url_data.headers

        print "\nresponse header content:", type(url_data)
        for header in url_data.headers.headers:
            print "    ", header,


def print_url_data_struct(url_data):    # 数据结构
    if url_data:
        for k, v in url_data.__dict__.items():
            print k, ":\t", v, "\n"

####################################################

def test_url():
    data = get_url_page_use_agent(URL)

    print_url_data_struct(data)        # URL数据结构
    data = get_url_page_data(data)     # 获取主页内容信息
    print "page content: ", data

def test_header():
    data = get_url_page_use_agent(URL)
    print_url_header(data)              # 测试HTTTP 头信息


if __name__ == '__main__':
    #test_url()
    test_header()
