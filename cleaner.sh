#!/bin/sh

# directory
DIR=`dirname $0`
cd ${DIR}

# 半年以上古いものは消す
find ./images/ -mtime +180 -exec rm -f {} \;
