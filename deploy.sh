#! /bin/bash

cd $(dirname $0)

echo $HOSTNAME

set -ex

sudo cp -a ./conf/all/etc/nginx/nginx.conf /etc/nginx/nginx.conf
sudo cp -a ./conf/all/etc/nginx/sites-available/isucondition.conf /etc/nginx/sites-available/isucondition.conf
sudo nginx -t &&  sudo systemctl restart nginx

sudo systemctl daemon-reload
sudo systemctl restart isucondition.go.service
