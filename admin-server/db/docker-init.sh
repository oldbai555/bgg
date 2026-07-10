#!/bin/sh
# 本地 docker-compose 场景下按依赖顺序初始化数据库：
# 建表（含增量建表）→ 基础数据 → 字典增量数据 → 模块初始化数据（菜单/权限/接口）。
# 挂载方式见 docker-compose.yml：整个 db/ 目录挂到 /db（只读，不放进
# /docker-entrypoint-initdb.d，避免 MySQL 官方镜像按文件名字典序自动执行、顺序不可控），
# 本脚本本身单独挂到 /docker-entrypoint-initdb.d/00-init.sh，由官方入口点唯一调用。
set -e

DB="${MYSQL_DATABASE:-admin}"
MYSQL="mysql -uroot -p${MYSQL_ROOT_PASSWORD} ${DB}"

echo "==> [1/5] 建表: tables.sql"
$MYSQL < /db/tables.sql

echo "==> [2/5] 增量建表: migrations/create_table_*.sql"
for f in /db/migrations/create_table_*.sql; do
  [ -e "$f" ] || continue
  echo "    -> $f"
  $MYSQL < "$f"
done

echo "==> [3/5] 基础初始化数据: data.sql"
$MYSQL < /db/data.sql

echo "==> [4/5] 字典增量数据: migrations/dict_*.sql"
for f in /db/migrations/dict_*.sql; do
  [ -e "$f" ] || continue
  echo "    -> $f"
  $MYSQL < "$f"
done

echo "==> [5/5] 模块初始化数据: migrations/init_*.sql"
for f in /db/migrations/init_*.sql; do
  [ -e "$f" ] || continue
  echo "    -> $f"
  $MYSQL < "$f"
done

echo "==> 数据库初始化完成"
