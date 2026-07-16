-- 修正 admin_menu.component 字段：前端 Phase 1-2 重构后按 9 个业务域重组了 src/views/
-- 目录结构（system/blog/sdk/video/chatroom/public 旧分组 -> iam/content/chat/monitoring/misc
-- 等新分组），但 oldbai 库里的菜单数据还停留在旧路径，导致动态路由解析失败。
-- 幂等：每条都是精确 id + 旧值匹配才更新，重复执行安全。

UPDATE `admin_menu` SET `component` = 'iam/RoleList' WHERE `id` = 3 AND `component` = 'system/RoleList';
UPDATE `admin_menu` SET `component` = 'iam/PermissionList' WHERE `id` = 4 AND `component` = 'system/PermissionList';
UPDATE `admin_menu` SET `component` = 'iam/DepartmentList' WHERE `id` = 5 AND `component` = 'system/DepartmentList';
UPDATE `admin_menu` SET `component` = 'iam/MenuList' WHERE `id` = 6 AND `component` = 'system/MenuList';
UPDATE `admin_menu` SET `component` = 'iam/UserList' WHERE `id` = 7 AND `component` = 'system/UserList';
UPDATE `admin_menu` SET `component` = 'iam/ApiList' WHERE `id` = 8 AND `component` = 'system/ApiList';
UPDATE `admin_menu` SET `component` = 'monitoring/OperationLogList' WHERE `id` = 44 AND `component` = 'system/OperationLogList';
UPDATE `admin_menu` SET `component` = 'monitoring/LoginLogList' WHERE `id` = 46 AND `component` = 'system/LoginLogList';
UPDATE `admin_menu` SET `component` = 'monitoring/AuditLogList' WHERE `id` = 49 AND `component` = 'system/AuditLogList';
UPDATE `admin_menu` SET `component` = 'monitoring/PerformanceLogList' WHERE `id` = 51 AND `component` = 'system/PerformanceLogList';
UPDATE `admin_menu` SET `component` = 'monitoring/MonitorList' WHERE `id` = 52 AND `component` = 'system/MonitorList';
UPDATE `admin_menu` SET `component` = 'chat/ChatList' WHERE `id` = 54 AND `component` = 'chatroom/ChatList';
UPDATE `admin_menu` SET `component` = 'chat/ChatMessageList' WHERE `id` = 55 AND `component` = 'chatroom/ChatMessageList';
UPDATE `admin_menu` SET `component` = 'chat/ChatGroupList' WHERE `id` = 57 AND `component` = 'chatroom/ChatGroupList';
UPDATE `admin_menu` SET `component` = 'misc/DailyShortSentenceList' WHERE `id` = 66 AND `component` = 'temp/DailyShortSentenceList';
UPDATE `admin_menu` SET `component` = 'content/VideoList' WHERE `id` = 83 AND `component` = 'video/VideoList';
UPDATE `admin_menu` SET `component` = 'content/VideoPlayer' WHERE `id` = 87 AND `component` = 'video/VideoPlayer';
UPDATE `admin_menu` SET `component` = 'task/TaskList' WHERE `id` = 88 AND `component` = 'system/TaskList';
UPDATE `admin_menu` SET `component` = 'content/BlogTagList' WHERE `id` = 92 AND `component` = 'blog/BlogTagList';
UPDATE `admin_menu` SET `component` = 'content/BlogArticleList' WHERE `id` = 93 AND `component` = 'blog/BlogArticleList';
UPDATE `admin_menu` SET `component` = 'content/BlogArticleAuditList' WHERE `id` = 94 AND `component` = 'blog/BlogArticleAuditList';
UPDATE `admin_menu` SET `component` = 'monitoring/MetricStats' WHERE `id` = 110 AND `component` = 'system/MetricStats';
UPDATE `admin_menu` SET `component` = 'content/BlogFriendLinkList' WHERE `id` = 111 AND `component` = 'blog/BlogFriendLinkList';
UPDATE `admin_menu` SET `component` = 'content/BlogSocialInfoList' WHERE `id` = 112 AND `component` = 'blog/BlogSocialInfoList';

-- 以下两条在当前 src/views/ 下找不到对应组件文件，是孤儿菜单项，不在本次修复范围内，
-- 需要人工确认后再处理（删除菜单，或者前端补一个对应页面）：
--   id 91  博客管理     component = blog/BlogLayout
--   id 106 metric_daily_stats  component = temp/MetricList（疑似 id 110 数据统计 的重复项）
