#!/bin/sh
# docker-compose 场景下的数据库初始化入口。
# 挂载方式见 docker-compose.yml：整个 db/ 目录挂到 /db（只读，不放进
# /docker-entrypoint-initdb.d，避免 MySQL 官方镜像按文件名字典序自动执行、顺序不可控），
# 本脚本本身单独挂到 /docker-entrypoint-initdb.d/00-init.sh，由官方入口点唯一调用。
#
# 实际初始化顺序见 db/services/init-dev-db.sh（db/services/ 目录树 + 该顺序脚本
# 是 docs/15-service-boundaries.md 第 4 节拆分后的唯一 SQL 来源，不再有 db/tables.sql/
# db/data.sql/db/migrations/ 这些历史文件）。
set -e
/db/services/init-dev-db.sh
