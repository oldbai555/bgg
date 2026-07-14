-- 修正「博客管理」根菜单排序值
-- 背景：admin-server/db/services/content/blog/init_blog.sql 当初脚手架生成时硬编码了
-- order_num=0，比「仪表盘」的 order_num=1 还小，导致侧边栏「博客管理」排到了仪表盘前面。
-- 现按用户确认的排序方案，把「博客管理」与「影视资源」「聊天室」同级（order_num=20），
-- 排在「系统管理」（10）之后。
-- 幂等性：UPDATE ... WHERE path = '/blog' AND order_num != 20，重复执行天然匹配 0 行

UPDATE admin_menu SET order_num = 20, updated_at = UNIX_TIMESTAMP()
WHERE path = '/blog' AND deleted_at = 0 AND order_num != 20;
