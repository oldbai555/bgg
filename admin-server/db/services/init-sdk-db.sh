#!/bin/sh
# docker-compose 场景下初始化 sdk-rpc 自己的 admin_sdk schema（和 db/docker-init.sh
# 初始化的主 admin 库是同一个 MySQL 实例、不同 schema，见 15-service-boundaries.md
# 第 4 节"每个 RPC 服务从第一天起就有自己独立的 MySQL schema"）。
#
# 只建表，不跑 db/services/sdk/sdk/init_sdk.sql——那个文件写的是 admin_menu/
# admin_permission/admin_api（iam 拥有的表，属于主 admin 库），不是 admin_sdk 自己的数据，
# 已经在 db/services/init-dev-db.sh 跑主库时处理过了，和 init-task-db.sh 是同一个模式。
set -e

DB="admin_sdk"
MYSQL="mysql -uroot -p${MYSQL_ROOT_PASSWORD}"

echo "==> 创建 ${DB} schema（如果不存在）"
$MYSQL -e "CREATE DATABASE IF NOT EXISTS \`${DB}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

echo "==> 建表: services/sdk/sdk/create_table_sdk.sql"
$MYSQL "$DB" < /db/services/sdk/sdk/create_table_sdk.sql

echo "==> ${DB} 初始化完成"
