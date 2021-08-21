#!/bin/sh
set -ex

cd icon
echo "SELECT CONCAT('http://localhost:3000/api/icon_for_devonly/', jia_isu_uuid) FROM isu" | sudo mysql -N -B isucondition | xargs wget
cd -
