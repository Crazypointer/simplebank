#!/bin/sh

set -e

source /app/app.env
echo "查看环境变量中是否有DB_SOURCE"
cat /app/app.env
echo $DB_SOURCE
echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"