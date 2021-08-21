#!/bin/bash

echo drop
curl -XPOST ${INFLUX_ADDR}'/query' --data-urlencode 'q=DROP DATABASE "isu"'
echo create
curl -XPOST ${INFLUX_ADDR}'/query' --data-urlencode 'q=CREATE DATABASE "isu"'
