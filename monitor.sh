#!/bin/bash

# 切换到当前工作目录
cd $(dirname `readlink -f "$0"`);

# 如果进程不存在了就重启
ps -fe | grep 'badwords_server' | grep -v 'grep' | grep -v 'monitor.sh'
if [ $? -ne 0 ]
then
    nohup ./badwords_server 0.0.0.0 8082 >> ./log/badwords.log 2>&1 &
    date "+%Y-%m-%d %H:%I:%S restarted" >> ./log/restart.log
fi
