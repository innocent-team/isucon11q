#!/bin/sh
set -ex

cd icon
echo "SELECT CONCAT('wget http://localhost:3000/api/isu/', jia_isu_uuid, '/icon -o ', jia_isu_uuid) FROM isu" | sudo mysql -N -B isucondition
cd -
