#!/bin/bash
set -ex

source /home/isucon/env.sh

export MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
export MYSQL_PORT=${MYSQL_PORT:-3306}
export MYSQL_USER=${MYSQL_USER:-isucon}
export MYSQL_DBNAME=${MYSQL_DBNAME:-isucondition}
export MYSQL_PWD=${MYSQL_PASS:-isucon}
export LANG="C.UTF-8"

if [[ ! -f /tmp/isucon-icon/completed ]]; then
  mkdir -p /tmp/isucon-icon
  cd /tmp/isucon-icon
  echo "SELECT CONCAT('http://localhost:3000/api/icon_for_devonly/', jia_isu_uuid) FROM isu" |  mysql --defaults-file=/dev/null -N -B -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p$MYSQL_PWD $MYSQL_DBNAME | xargs wget
  touch /tmp/isucon-icon/completed
  cd -
fi
