# -*- coding:utf-8  -*-


CSRF_ENABLED = True      # CSRF_ENABLED 激活 跨站点请求伪造 保护。激活该配置 使得你的应用程序更安全。
SECRET_KEY = "TEST_BLOG" # SECRET_KEY 配置仅仅当 CSRF 激活时才需要，用来建立一个加密令牌，验证表单。务必设置很难被猜测到密钥
