#! /bin/sh

# Wait for MySQL
until nc -z -v -w30 mysql 3306; do
 echo 'Waiting for MySQL...'
 sleep 1
done
echo "MySQL is up and running!"