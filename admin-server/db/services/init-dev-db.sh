#!/bin/sh
# 单体阶段（Phase 2 各服务真正拆分成独立进程/独立库之前）的开发库/CI 库初始化脚本。
#
# db/tables.sql + db/data.sql + db/migrations/* + db/demo/* 已经按
# docs/15-service-boundaries.md 第 4 节拆分到 db/services/<service>/<module>/ 下，
# 但当前 admin-server 仍是单体单库，所有服务的表都建在同一个数据库里，
# 所以这里按明确的依赖顺序把全部服务的 SQL 跑一遍（而不是 glob 乱序遍历）：
#
#   iam 建表 → iam 初始化数据（按 iam 内部依赖顺序）→ 其余各服务建表 → 其余各服务初始化数据
#
# 之所以 iam 必须整体排最前：iam 拥有 admin_menu/admin_permission/admin_api/
# admin_permission_menu/admin_permission_api 等全局共享的 RBAC 表，其余服务的
# 初始化数据（如 chat 的群组管理菜单/权限关联、blog/video/sdk 的菜单权限接口）
# 都通过 SELECT 反查这些表里已存在的行，必须等 iam 初始化跑完。
#
# 各服务/模块内部固定 create_table → init → migrations/dict_*.sql（字典增量脚本，
# 幂等、全新库上跑没有副作用）的顺序自动执行。migrations/ 下面文件名不是 dict_ 前缀的，
# 是一次性存量数据迁移脚本（如 iam/notification/migrations/fix_status_semantics_20260711.sql、
# iam/api/migrations/fix_chat_group_api_paths_20260711.sql），只对已经跑过旧版种子数据、
# 积累了真实业务数据的库有意义，全新库不需要、部分甚至不是幂等的（见各脚本文件头部说明），
# 因此故意不被下面的 glob 匹配到，必须手动执行。
set -e

# 用法：init-dev-db.sh [-h<host>]
# docker-compose 场景下脚本在 MySQL 容器内跑，连本机（无需 -h）；
# CI 场景下 mysql client 跑在 runner 上，需要 -h127.0.0.1 连到 service 容器。
DB="${MYSQL_DATABASE:-admin}"
MYSQL="mysql -uroot -p${MYSQL_ROOT_PASSWORD} $1 ${DB}"
SERVICES_DIR="$(cd "$(dirname "$0")" && pwd)"

run() {
    f="$1"
    [ -e "$f" ] || return 0
    echo "    -> $f"
    $MYSQL < "$f"
}

run_module() {
    # $1 = 模块目录（如 iam/user），依次跑 create_table -> init -> migrations/*
    dir="${SERVICES_DIR}/$1"
    module="$(basename "$1")"
    run "${dir}/create_table_${module}.sql"
    run "${dir}/init_${module}.sql"
    # 只自动执行字典增量脚本（dict_*.sql，幂等）；一次性存量数据修复脚本（fix_*.sql 等）
    # 故意不匹配，必须手动执行，见上面的说明。
    for f in "${dir}"/migrations/dict_*.sql; do
        [ -e "$f" ] || continue
        run "$f"
    done
}

echo "==> [1/4] iam 建表 + 初始化数据（内部依赖顺序）"
for m in user department role permission menu api rbac config dict file \
         notice notification operation_log login_log audit_log performance_log \
         monitor metric demo daily_short_sentence; do
    run_module "iam/${m}"
done

echo "==> [2/4] content 建表 + 初始化数据"
for m in blog blog_extension video; do
    run_module "content/${m}"
done

echo "==> [3/4] chat 建表 + 初始化数据"
run_module "chat/chat"

echo "==> [4/4] task、sdk 建表 + 初始化数据"
run_module "task/task"
run_module "sdk/sdk"

echo "==> 数据库初始化完成"
