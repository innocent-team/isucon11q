#! /bin/bash

cd $(dirname $0)

# 1 -> ip-192-168-0-11
# 2 -> ip-192-168-0-12
# 3 -> ip-192-168-0-13
echo $HOSTNAME

if [[ "$HOSTNAME" == ip-192-168-0-11 ]]; then
  INSTANCE_NUM="1"
elif [[ "$HOSTNAME" == ip-192-168-0-12 ]]; then
  INSTANCE_NUM="2"
elif [[ "$HOSTNAME" == ip-192-168-0-13 ]]; then
  INSTANCE_NUM="3"
else
  echo "Invalid host"
  exit 1
fi

set -ex

sudo systemctl daemon-reload

echo "Restarting App"
pushd go
go build
sudo systemctl enable isucondition.go.service
sudo systemctl restart isucondition.go.service
popd

echo "Restarting nginx"
sudo cp -a ./conf/all/etc/nginx/nginx.conf /etc/nginx/nginx.conf
sudo cp -a ./conf/all/etc/nginx/sites-available/isucondition.conf /etc/nginx/sites-available/isucondition.conf
sudo systemctl enable nginx
sudo nginx -t &&  sudo systemctl restart nginx

if [[ "$INSTANCE_NUM" == 2 ]]; then
  echo "Restarting mysql"
  sudo cp -a ./conf/$INSTANCE_NUM/etc/mysql/conf.d/my.cnf /etc/mysql/conf.d/my.cnf
  sudo cp -a ./conf/$INSTANCE_NUM/etc/mysql/conf.d/mysql.cnf /etc/mysql/conf.d/mysql.cnf
  sudo cp -a ./conf/$INSTANCE_NUM/etc/mysql/conf.d/mysqldump.cnf /etc/mysql/conf.d/mysqldump.cnf
  # sudo systemctl enable mysql
  sudo systemctl restart mysql
fi

if [[ "$INSTANCE_NUM" == 3 ]]; then
  echo "Restarting influxdb"
  sudo cp -a ./conf/$INSTANCE_NUM/etc/influxdb/influxdb.conf /etc/influxdb/influxdb.conf
  sudo systemctl enable influxdb
  sudo systemctl restart influxdb
fi
