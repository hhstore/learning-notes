# -*- coding:utf-8  -*-

from flask import Flask
app = Flask(__name__)
app.config.from_object('config')    # 告诉 Flask 去读取,并使用config.py配置参数


