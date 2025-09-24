#!/bin/sh

# 启动后端应用，日志直接输出到控制台
/app/datax-admin &

# 启动 nginx，后台运行，不输出日志
nginx -g 'daemon off;' > /dev/null 2>&1