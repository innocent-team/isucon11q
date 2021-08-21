#! /bin/bash

cd $(dirname $0)

echo $HOSTNAME

set -ex

sudo systemctl daemon-reload

echo "Restarting App"
sudo systemctl restart isucondition.go.service

echo "Restarting nginx"
sudo cp -a ./conf/all/etc/nginx/nginx.conf /etc/nginx/nginx.conf
sudo cp -a ./conf/all/etc/nginx/sites-available/isucondition.conf /etc/nginx/sites-available/isucondition.conf
sudo nginx -t &&  sudo systemctl restart nginx

echo "Restarting mysql"
sudo cp -a ./conf/etc/mysql/conf.d/my.cnf /etc/mysql/conf.d/my.cnf
sudo cp -a ./conf/etc/mysql/conf.d/mysql.cnf /etc/mysql/conf.d/mysql.cnf
sudo cp -a ./conf/etc/mysql/conf.d/mysqldump.cnf /etc/mysql/conf.d/mysqldump.cnf
sudo systemctl restart mysqld
