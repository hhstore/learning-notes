# -*- coding:utf-8 -*-
__author__ = 'hhstore'

'''
功能: 读写office文档,支持docx.
依赖: docx模块
'''


from docx import Document


def parse_docx(in_file, out_file):
    doc = Document(in_file)

    for item in doc.paragraphs:   # 通过段落解析内容.
        print item.text

    doc.save(out_file)


if __name__ == '__main__':
    infile = "test.docx"
    outfile = "out.docx"
    parse_docx(infile, outfile)
