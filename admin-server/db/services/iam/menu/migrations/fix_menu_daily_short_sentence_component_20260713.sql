-- admin-frontend Phase 1 Week 2 死代码清理后的 admin_menu.component 数据修正
-- 背景：admin-frontend/docs/07-cleanup-and-tooling.md §1
-- views/temp/DailyShortSentenceList.vue 核实后发现不是死代码——它管理的"每日一言"数据正被
-- src/views/Dashboard.vue 实际展示消费，只是管理页面本身从未真正挪出"临时目录"。
-- 已 git mv 到 admin-frontend/src/views/misc/DailyShortSentenceList.vue，本次只修正 component 字段值。
-- 菜单本身仍挂在 /temp/daily_short_sentence 下（临时目录），按 00-workflow.md 脚手架约定，
-- 挪到正式分类是后续在菜单管理界面手动操作的事，不在本次 SQL 范围内。
-- 幂等性：UPDATE ... WHERE component = '<旧值>'，旧值执行一次后不再存在，重复执行时 WHERE 天然匹配 0 行

UPDATE admin_menu SET component = 'misc/DailyShortSentenceList' WHERE component = 'temp/DailyShortSentenceList';
