# -*- coding:utf8 -*-
__author__ = 'hhstore'

'''
function:  spider

module:
    1. celery
    2. tornado
    3. rabbitmq

'''


from celery import Celery
from tornado import httpclient
from tornado.httpclient import HTTPClient


app = Celery("tasks")
app.config_from_object("celery_config")

@app.task
def get_html(url):
    http_client = HTTPClient()

    try:
        response = http_client.fetch(url, follow_redirects=True)
        print("body: {}".format(response.body))
        data = response.body
        result = str(data).encode(encoding="utf-8")
        return result
    except httpclient.HTTPError as e:
        return None
    finally:
        http_client.close()

