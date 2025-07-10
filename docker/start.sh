#!/bin/sh

# 启动后端应用，日志直接输出到控制台
/app/datax-admin &

# 启动 nginx，日志输出到文件
nginx -g 'daemon off;'  