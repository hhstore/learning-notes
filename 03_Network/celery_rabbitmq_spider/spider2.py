from tasks import get_html
from queue import Queue
from bs4 import BeautifulSoup
from urllib.parse import urlparse, urljoin
import threading


class spider(object):
    def __init__(self):
        self.visited = {}
        self.queue = Queue()

    def process_html(self, html):
        pass
        # print(html)

    def _add_links_to_queue(self, url_base, html):
        soup = BeautifulSoup(html, from_encoding="utf8")
        links = soup.find_all('a')
        for link in links:
            try:
                url = link['href']
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
        while 1:
            url = self.queue.get()
            if url in self.visited:
                continue
            else:
                html = None
                result = get_html.delay(url)
                try:
                    html = result.get(timeout=5)
                except Exception as e:
                    print(url)
                    print(e)

                if html:
                    self.process_html(html)
                    self._add_links_to_queue(url, html)
                    self.visited[url] = True
                    self.queue.task_done()


if __name__ == '__main__':
    url = "https://github.com/hhstore"

    s = spider()
    s.start(url)
