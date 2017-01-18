# -*- coding:utf-8 -*-
__author__ = 'hhstore'

'''
功能: 解析markdown文件,并写出HTML网页.

依赖: markdown模块, markdown2模块.

说明:
    1. 2个模块,都简单测试了一下.功能差不多.
    2. 都可以写出正常的HTML,略有差异,一个多了些空行.
    3. 后续根据需要,再自行应用.

'''



import codecs
import markdown
import markdown2


md_content = """
Who am I?
===========
This is Ray Chen speaking.

Where am I from?
====================
I'm from SH, China. I love this place.

What's next?
=================
  * nothing
  * sleep
  * study
  * reading book

```python

    print 'hello'
```

+ 1.hi.
+ 2.thanks.
+ 3.fine.

"""



class Markdown2HTML(object):
    def __init__(self, data, outfile):
        self.data = data
        self.outfile = outfile
        self.parse_content = None

    def parse1(self):
        self.parse_content = markdown.markdown(self.data) if self.data else None

    def parse2(self):
        self.parse_content = markdown2.markdown(self.data) if self.data else None

    def save(self):
        if self.parse_content:
            with codecs.open(self.outfile, mode="w", encoding="utf8") as f:
                f.write(self.parse_content)


def main():
    html1 = "md1.html"
    html2 = "md2.html"

    md1 = Markdown2HTML(md_content, html1)
    md1.parse1()
    md1.save()

    md2 = Markdown2HTML(md_content, html2)
    md2.parse2()
    md2.save()



if __name__ == '__main__':
    main()

