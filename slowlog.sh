#! /bin/bash

set -ex

mkdir /data/log_`date -Iminute`
cat /var/log/mysql/mysql-slow.log | pt-query-digest > /data/log_`date -Iminute`/slow_`date -Iminute`.log
#sudo cat /var/log/nginx/access.log | kataribe > /data/log_`date -Iminute`/kataribe_`date -Iminute`.log
