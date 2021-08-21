#!/bin/bash

echo drop
curl -XPOST 'http://localhost:8086/query' --data-urlencode 'q=DROP DATABASE "isu"'
echo create
curl -XPOST 'http://localhost:8086/query' --data-urlencode 'q=CREATE DATABASE "isu"'
