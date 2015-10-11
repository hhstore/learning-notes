# -*- coding:utf8 -*-
__author__ = 'hhstore'


'''
说明: Python3版本

'''


import platform

if platform.python_version_tuple()[0] == "2":  # python2.7
    from Queue import Queue
else:
    from queue import Queue           # python 3.

from bs4 import BeautifulSoup
from urllib.parse import urlparse, urljoin



import threading

from tasks import get_html

class Spider(object):
    def __init__(self):
        self.visited = {}
        self.queue = Queue()

    def process_html(self, html):
        print("html: {}".format(html))
        pass

    def _add_links_to_queue(self, url_base, html):
        soup = BeautifulSoup(html)
        links = soup.find_all("a")

        for link in links:
            try:
                url = link["href"]
            except:
                pass
            else:
                url_com = urlparse(url)
                if not url_com.netloc:
                    self.queue.put(urljoin(url_base, url))
                else:
                    self.queue.put(url_com.geturl())

    def start(self, url):
        self.queue.put(url)

        for i in range(20):
            t = threading.Thread(target=self._worker)
            t.daemon = True
            t.start()
        self.queue.join()

    def _worker(self):
        while True:
            url = self.queue.get()
            if url in self.visited:
                continue
            else:
                result = get_html.delay(url)
                html = None
                try:
                    html = result.get(timeout=5)
                except Exception as e:
                    print(url)
                    print(e)
                if html:
                    self.process_html(html)    # 处理爬取的页面.
                    self._add_links_to_queue(url, html)

                self.visited[url] = True
                self.queue.task_done()


def main():
    url = "http://cuiqingcai.com/"
    s = Spider()
    s.start(url)

if __name__ == '__main__':
    main()
