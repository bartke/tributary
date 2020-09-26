#!/bin/sh

set -e

mysql="$1"

echo 'Waiting for mysql service start...';
while ! nc -z $mysql 3306; do sleep 1; done;
echo 'Connected to mysql!';

echo 'Creating mysql database'
mysql -h $mysql -u root --password=root --protocol=tcp -e "CREATE DATABASE IF NOT EXISTS tb_test character set utf8mb4 collate utf8mb4_general_ci;"

echo 'Running'
cd example/advanced && go run *.go $mysql