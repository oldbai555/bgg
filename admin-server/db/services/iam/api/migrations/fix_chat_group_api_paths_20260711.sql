-- 修正群组详情/成员列表接口路径，与 admin.api 对齐（禁止路径参数 :id）
-- 幂等：仅当旧路径仍存在时才更新

UPDATE `admin_api`
SET `path` = '/api/v1/chats/groups/detail', `updated_at` = UNIX_TIMESTAMP()
WHERE `method` = 'GET' AND `path` = '/api/v1/chats/groups/:id' AND `deleted_at` = 0;

UPDATE `admin_api`
SET `path` = '/api/v1/chats/groups/members', `updated_at` = UNIX_TIMESTAMP()
WHERE `method` = 'GET' AND `path` = '/api/v1/chats/groups/:id/members' AND `deleted_at` = 0;
