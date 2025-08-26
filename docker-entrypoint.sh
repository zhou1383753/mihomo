#!/bin/sh

/etc/clash/start.sh

cd /root/clash-admin
/usr/share/bin/linux >> /var/log/run.log 2>&1 &
exec /sbin/tini -- sh -c "tail -f /var/log/run.log"
