
## 使用Python+celery+threading+rabbitmq 实现爬虫

## 主要模块


## 启动:
1. 先启动rabbitmq
运行方式: 

```bash

# mac下:
/usr/local/sbin/rabbitmq-server

重复执行该命令,显示:
$ /usr/local/sbin/rabbitmq-server

ERROR: node with name "rabbit" already running on "localhost"


```

2. 启动celery.

```
cd 到当前项目目录下:

celery -A tasks worker --loglevel=info

```


3. 执行spider.py

```bash
python spider.py
```



## 退出服务:

1. 退出rabbitmq:

```bash
新开一个终端,执行如下命令:
$ /usr/local/sbin/rabbitmqctl stop

Stopping and halting node rabbit@localhost ...

```


2. 退出celery.

```bash
ctrl + C

```
