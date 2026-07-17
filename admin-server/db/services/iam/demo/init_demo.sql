-- 演示功能模块初始化 SQL
-- 功能组: demo
-- 功能名称: 演示功能

-- ============================================
-- 2. 演示功能没有独立菜单
-- ============================================
-- views/temp/DemoList.vue 是脚手架示例页面，已在 admin-frontend Phase 1 Week 2
-- 死代码清理中删除（admin-frontend/docs/07-cleanup-and-tooling.md §1），不再为它建
-- admin_menu 记录，避免产生指向不存在页面的孤儿菜单。demo 模块本身作为
-- generate-sql.sh 脚手架产出示例保留，权限/接口数据照常初始化，仅不挂菜单。

-- ============================================
-- 3. 插入权限数据
-- ============================================
-- 演示功能列表权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('演示功能列表', 'demo:list', '查看演示功能列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 演示功能新增权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('演示功能新增', 'demo:create', '新增演示功能', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 演示功能编辑权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('演示功能编辑', 'demo:update', '编辑演示功能', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 演示功能删除权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('演示功能删除', 'demo:delete', '删除演示功能', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='demo:list' AND `deleted_at`=0 LIMIT 1);
SET @create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='demo:create' AND `deleted_at`=0 LIMIT 1);
SET @update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='demo:update' AND `deleted_at`=0 LIMIT 1);
SET @delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='demo:delete' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 4. 插入接口数据
-- ============================================
-- 演示功能接口：path 必须与 admin-server/api/admin.api 里 misc/demo 服务块实际路由一致
-- （get/post/put/delete /demos，禁止路径参数 :id，见 00-workflow.mdc 绝对禁止事项表）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('演示功能列表', 'GET', '/api/v1/demos', '获取演示功能列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('演示功能新增', 'POST', '/api/v1/demos', '新增演示功能', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('演示功能编辑', 'PUT', '/api/v1/demos', '编辑演示功能', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('演示功能删除', 'DELETE', '/api/v1/demos', '删除演示功能', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='GET' AND `path`='/api/v1/demos' AND `deleted_at`=0 LIMIT 1);
SET @create_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='POST' AND `path`='/api/v1/demos' AND `deleted_at`=0 LIMIT 1);
SET @update_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='PUT' AND `path`='/api/v1/demos' AND `deleted_at`=0 LIMIT 1);
SET @delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='DELETE' AND `path`='/api/v1/demos' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 6. 插入权限-接口关联数据
-- ============================================
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();
