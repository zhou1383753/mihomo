#!/bin/sh

killall -9 mihomo
/mihomo -d /etc/clash >> /var/log/run.log 2>&1 &