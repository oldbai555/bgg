-- 每日短句模块初始化 SQL
-- 功能组: daily_short_sentence
-- 功能名称: 每日短句

-- ============================================
-- 1. 获取父目录 ID（可配置）
-- ============================================
-- 如果提供了 ParentID，则直接使用该 ID 作为父菜单；
-- 否则根据 ParentPath 从 admin_menu 中查找父目录；
-- 如果仍未找到，默认回退到临时目录（id = 9）。
SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/temp' AND `deleted_at` = 0 LIMIT 1),
  (SELECT `id` FROM `admin_menu` WHERE `id` = 9 AND `deleted_at` = 0 LIMIT 1)
);

-- ============================================
-- 2. 插入菜单数据
-- ============================================
-- 每日短句主菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '每日短句',
    '/temp/daily_short_sentence',
    'temp/DailyShortSentenceList',
    'ele-Document',
    2, -- 类型：2 菜单
    0, -- 排序值（可根据需要调整）
    1, -- 是否可见：1 是
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

-- 获取主菜单 ID
SET @main_menu_id = LAST_INSERT_ID();

-- 每日短句新增按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @main_menu_id,
    '每日短句 新增按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否（按钮不显示在菜单中）
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @create_button_id = LAST_INSERT_ID();

-- 每日短句编辑按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @main_menu_id,
    '每日短句 编辑按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    2, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @update_button_id = LAST_INSERT_ID();

-- 每日短句删除按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @main_menu_id,
    '每日短句 删除按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    3, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @delete_button_id = LAST_INSERT_ID();

-- ============================================
-- 3. 插入权限数据
-- ============================================
-- 每日短句列表权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句列表',
    'daily_short_sentence:list',
    '查看每日短句列表',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @list_permission_id = LAST_INSERT_ID();

-- 每日短句新增权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句新增',
    'daily_short_sentence:create',
    '新增每日短句',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @create_permission_id = LAST_INSERT_ID();

-- 每日短句编辑权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句编辑',
    'daily_short_sentence:update',
    '编辑每日短句',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @update_permission_id = LAST_INSERT_ID();

-- 每日短句删除权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句删除',
    'daily_short_sentence:delete',
    '删除每日短句',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @delete_permission_id = LAST_INSERT_ID();

-- ============================================
-- 4. 插入接口数据
-- ============================================
-- 每日短句列表接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句列表',
    'GET',
    '/api/v1/daily-short-sentences',
    '获取每日短句列表',
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @list_api_id = LAST_INSERT_ID();

-- 每日短句新增接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句新增',
    'POST',
    '/api/v1/daily-short-sentences',
    '新增每日短句',
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @create_api_id = LAST_INSERT_ID();

-- 每日短句编辑接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句编辑',
    'PUT',
    '/api/v1/daily-short-sentences',
    '编辑每日短句',
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @update_api_id = LAST_INSERT_ID();

-- 每日短句删除接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '每日短句删除',
    'DELETE',
    '/api/v1/daily-short-sentences',
    '删除每日短句',
    1, -- 状态：1 启用（可根据需要设置为 0 禁用）
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @delete_api_id = LAST_INSERT_ID();

-- ============================================
-- 5. 插入权限-菜单关联数据
-- ============================================
-- 每日短句列表权限 -> 每日短句主菜单
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句新增权限 -> 每日短句新增按钮
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@create_permission_id, @create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句编辑权限 -> 每日短句编辑按钮
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@update_permission_id, @update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句删除权限 -> 每日短句删除按钮
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@delete_permission_id, @delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- ============================================
-- 6. 插入权限-接口关联数据
-- ============================================
-- 每日短句列表权限 -> GET /api/v1/daily-short-sentences接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句新增权限 -> POST /api/v1/daily-short-sentences接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句编辑权限 -> PUT /api/v1/daily-short-sentences接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 每日短句删除权限 -> DELETE /api/v1/daily-short-sentences接口
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

