-- admin-frontend Phase 1 Week 1 域目录重组后的 admin_menu.component 数据修正
-- 背景：admin-frontend/docs/02-domain-reorg-and-api-layer.md §2；仅修正 component 字段值，不新增/删除菜单记录
-- 幂等性：UPDATE ... WHERE component = '<旧值>'，旧值执行一次后不再存在，重复执行时 WHERE 天然匹配 0 行
-- 注意：views/temp/* 下的孤儿页面（temp/DailyShortSentenceList、temp/MetricList 对应的菜单记录）
--      是否清理属于 07-cleanup-and-tooling.md 的 Week 2 范围，本次不处理

-- iam 域（原 system/* 拆出）
UPDATE admin_menu SET component = 'iam/UserList' WHERE component = 'system/UserList';
UPDATE admin_menu SET component = 'iam/RoleList' WHERE component = 'system/RoleList';
UPDATE admin_menu SET component = 'iam/PermissionList' WHERE component = 'system/PermissionList';
UPDATE admin_menu SET component = 'iam/DepartmentList' WHERE component = 'system/DepartmentList';
UPDATE admin_menu SET component = 'iam/MenuList' WHERE component = 'system/MenuList';
UPDATE admin_menu SET component = 'iam/ApiList' WHERE component = 'system/ApiList';

-- monitoring 域（原 system/* 拆出）
UPDATE admin_menu SET component = 'monitoring/AuditLogList' WHERE component = 'system/AuditLogList';
UPDATE admin_menu SET component = 'monitoring/LoginLogList' WHERE component = 'system/LoginLogList';
UPDATE admin_menu SET component = 'monitoring/OperationLogList' WHERE component = 'system/OperationLogList';
UPDATE admin_menu SET component = 'monitoring/PerformanceLogList' WHERE component = 'system/PerformanceLogList';
UPDATE admin_menu SET component = 'monitoring/MonitorList' WHERE component = 'system/MonitorList';
UPDATE admin_menu SET component = 'monitoring/MetricStats' WHERE component = 'system/MetricStats';

-- task 域（原 system/TaskList 拆出，后端已是独立 task-rpc）
UPDATE admin_menu SET component = 'task/TaskList' WHERE component = 'system/TaskList';

-- content 域（原 blog/* + video/* 合并）
UPDATE admin_menu SET component = 'content/BlogTagList' WHERE component = 'blog/BlogTagList';
UPDATE admin_menu SET component = 'content/BlogArticleList' WHERE component = 'blog/BlogArticleList';
UPDATE admin_menu SET component = 'content/BlogArticleAuditList' WHERE component = 'blog/BlogArticleAuditList';
UPDATE admin_menu SET component = 'content/BlogFriendLinkList' WHERE component = 'blog/BlogFriendLinkList';
UPDATE admin_menu SET component = 'content/BlogSocialInfoList' WHERE component = 'blog/BlogSocialInfoList';
UPDATE admin_menu SET component = 'content/VideoList' WHERE component = 'video/VideoList';
UPDATE admin_menu SET component = 'content/VideoPlayer' WHERE component = 'video/VideoPlayer';

-- chat 域（原 chatroom/* 改名）
UPDATE admin_menu SET component = 'chat/ChatList' WHERE component = 'chatroom/ChatList';
UPDATE admin_menu SET component = 'chat/ChatMessageList' WHERE component = 'chatroom/ChatMessageList';
UPDATE admin_menu SET component = 'chat/ChatGroupList' WHERE component = 'chatroom/ChatGroupList';
