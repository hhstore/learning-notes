__author__ = 'hhstore'

CELERY_IMPORT = ("tasks",)
BROKER_URL = "amqp://guest@localhost:5672//"
CELERY_RESULT_BACKEND = "amqp://"
