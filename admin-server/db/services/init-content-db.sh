#!/bin/sh
# docker-compose 场景下初始化 content-rpc 自己的 admin_content schema（和 db/docker-init.sh
# 初始化的主 admin 库是同一个 MySQL 实例、不同 schema，见 15-service-boundaries.md
# 第 4 节"每个 RPC 服务从第一天起就有自己独立的 MySQL schema"）。
#
# 只建表（blog/blog_extension/video 三个模块共 7 张表），不跑对应的 init_*.sql——那几个
# 文件写的是 admin_menu/admin_permission/admin_api（iam 拥有的表，属于主 admin 库），不是
# admin_content 自己的数据，已经在 db/services/init-dev-db.sh 跑主库时处理过了，和
# init-task-db.sh/init-sdk-db.sh 是同一个模式。
set -e

DB="admin_content"
MYSQL="mysql -uroot -p${MYSQL_ROOT_PASSWORD}"

echo "==> 创建 ${DB} schema（如果不存在）"
$MYSQL -e "CREATE DATABASE IF NOT EXISTS \`${DB}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

echo "==> 建表: services/content/blog/create_table_blog.sql"
$MYSQL "$DB" < /db/services/content/blog/create_table_blog.sql

echo "==> 建表: services/content/blog_extension/create_table_blog_extension.sql"
$MYSQL "$DB" < /db/services/content/blog_extension/create_table_blog_extension.sql

echo "==> 建表: services/content/video/create_table_video.sql"
$MYSQL "$DB" < /db/services/content/video/create_table_video.sql

echo "==> ${DB} 初始化完成"
