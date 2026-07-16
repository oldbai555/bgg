-- admin-frontend Phase 1 Week 2 死代码清理后遗留的孤儿 admin_menu 记录清理
-- 背景：admin-frontend/docs/07-cleanup-and-tooling.md §1
-- views/temp/DemoList.vue（明确是开发流程示例脚手架）和 views/temp/MetricList.vue（功能已被
-- views/monitoring/MetricStats.vue 取代）在 a3e1af5 中已被确认删除，但对应的 admin_menu 记录
-- （含各自「新增/编辑/删除」按钮子菜单）当时未同步清理，导致这两条菜单一直是指向不存在页面的
-- 孤儿记录，直接访问 /temp/demo、/temp/metric 会 404。
-- 本次仅软删除 admin_menu 记录本身，不涉及 admin_permission/admin_api 及其关联表
-- （对应后端接口是否保留是另一个决策，不在本次范围内）。
-- 幂等性：UPDATE ... WHERE deleted_at = 0，软删除后再次执行天然匹配 0 行

UPDATE admin_menu SET deleted_at = UNIX_TIMESTAMP(), updated_at = UNIX_TIMESTAMP()
WHERE deleted_at = 0 AND (
  path IN ('/temp/demo', '/temp/metric')
  OR parent_id IN (
    SELECT id FROM (SELECT id FROM admin_menu WHERE path IN ('/temp/demo', '/temp/metric')) AS t
  )
);
