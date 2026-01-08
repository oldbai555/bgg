-- admin-server 数据库初始化数据脚本
-- 说明：此脚本包含所有系统初始化数据，所有 INSERT 语句使用 ON DUPLICATE KEY UPDATE 确保幂等性，可重复执行
-- 注意：此脚本中的数据为系统初始化数据，不可被删除（包括软删、硬删）
-- 
-- 初始化数据的ID范围（每张表从1开始连续）：
--   admin_user: id=1-2 (1=超级管理员, 2=admin 业务管理员)
--   admin_role: id=1-2 (1=super_admin 超级管理员角色, 2=admin 业务管理员角色)
--   admin_permission: id=1-48+ (基础48个权限，含通用权限 common:xxx，后续模块会新增)
--   admin_department: id=1 (根部门)
--   admin_menu: id=1-43+ (基础13个菜单 + 30个按钮，后续模块会新增)
--   admin_api: id=1-58+ (基础58个接口，后续模块会新增)
--   admin_user_role: id=1-2 (1=super_admin绑定超级管理员, 2=admin绑定业务管理员)
--   admin_role_permission: id=1-4+ (超级管理员角色-权限关联，后续模块会新增)
--   admin_permission_menu: id=1-40+ (基础10个菜单关联 + 30个按钮关联，后续模块会新增)
--   admin_permission_api: id=1-58+ (基础58个权限-接口关联，后续模块会新增)
--   admin_dict_type: id=1-6 (6个字典类型：用户状态、性别、是否、文件存储类型、聊天配置、消息来源类型)
--   admin_dict_item: id=1-17 (17个字典项，包含emoji分页配置)
--   admin_notice: id=1 (1条初始化公告)
--   chat: id=1 (1个默认企业群组)
--   chat_user: id=1-3 (默认群组包含2个用户，1个私聊包含2个用户)

-- ============================================
-- 1. 初始化基础数据
-- ============================================
-- 初始化用户：1=超级管理员 oldbai，2=业务管理员 admin（密码：admin）
INSERT INTO `admin_user` (`id`, `username`, `password_hash`, `department_id`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (1, 'oldbai', '$2a$10$TIjB8/yhHDiyNbJn40BUPOACjxeTccaYTD4Ot3p00ZBCKzh7/sL9q', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'admin', '$2a$10$F/GioZ0D2TUl7wQX2kErU.fuqu/IJvU8yd.VtuFXbVfcAYPaZaj7S', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `username`=VALUES(`username`), 
  `password_hash`=VALUES(`password_hash`), 
  `status`=VALUES(`status`), 
  `deleted_at`=0;

-- 根部门
INSERT INTO `admin_department` (`id`, `parent_id`, `name`, `order_num`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (1, 0, '总部', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 初始化角色：1=super_admin 超级管理员角色，2=admin 业务管理员角色
INSERT INTO `admin_role` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (1, '超级管理员', 'super_admin', '系统内置最高权限角色，拥有全部权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'admin', 'admin', '系统内置业务管理员角色，示例账号使用', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 权限列表（完整，ID从1开始连续）
INSERT INTO `admin_permission` (`id`, `name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  -- 超级权限（通配）
  (1, '全部权限', '*', '超级管理员拥有的全量权限通配符', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 角色管理权限
  (2, '角色列表', 'role:list', '查看角色列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, '角色新增', 'role:create', '新增角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, '角色编辑', 'role:update', '编辑角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, '角色删除', 'role:delete', '删除角色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 权限管理权限
  (6, '权限列表', 'permission:list', '查看权限列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, '权限新增', 'permission:create', '新增权限', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, '权限编辑', 'permission:update', '编辑权限', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, '权限删除', 'permission:delete', '删除权限', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 部门管理权限
  (10, '部门树', 'department:tree', '查看部门树', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (11, '部门新增', 'department:create', '新增部门', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, '部门编辑', 'department:update', '编辑部门', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (13, '部门删除', 'department:delete', '删除部门', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 菜单管理权限
  (14, '菜单列表', 'menu:list', '查看菜单列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, '菜单新增', 'menu:create', '新增菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, '菜单编辑', 'menu:update', '编辑菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, '菜单删除', 'menu:delete', '删除菜单', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 用户管理权限
  (18, '用户列表', 'user:list', '查看用户列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (19, '用户新增', 'user:create', '新增用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (20, '用户编辑', 'user:update', '编辑用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (21, '用户删除', 'user:delete', '删除用户', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 接口管理权限
  (22, '接口列表', 'api:list', '查看接口列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (23, '接口新增', 'api:create', '新增接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (24, '接口编辑', 'api:update', '编辑接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (25, '接口删除', 'api:delete', '删除接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 系统配置权限
  (26, '系统配置列表', 'config:list', '查看系统配置列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (27, '系统配置新增', 'config:create', '新增系统配置', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (28, '系统配置编辑', 'config:update', '编辑系统配置', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (29, '系统配置删除', 'config:delete', '删除系统配置', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典类型权限
  (30, '字典类型列表', 'dict_type:list', '查看字典类型列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (31, '字典类型新增', 'dict_type:create', '新增字典类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (32, '字典类型编辑', 'dict_type:update', '编辑字典类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (33, '字典类型删除', 'dict_type:delete', '删除字典类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典项权限
  (34, '字典项列表', 'dict_item:list', '查看字典项列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (35, '字典项新增', 'dict_item:create', '新增字典项', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (36, '字典项编辑', 'dict_item:update', '编辑字典项', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (37, '字典项删除', 'dict_item:delete', '删除字典项', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 文件管理权限
  (38, '文件列表', 'file:list', '查看文件列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (39, '文件新增', 'file:create', '新增文件', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (40, '文件编辑', 'file:update', '编辑文件', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (41, '文件删除', 'file:delete', '删除文件', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 群组管理权限（需要权限控制）
  (42, '群组创建', 'chat:group:create', '创建群组', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (43, '群组编辑', 'chat:group:update', '编辑群组信息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (44, '群组删除', 'chat:group:delete', '删除群组', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (45, '群组详情', 'chat:group:detail', '查看群组详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (46, '群组成员管理', 'chat:group:member', '管理群组成员（添加/移除）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`), 
  `code`=VALUES(`code`), 
  `description`=VALUES(`description`), 
  `updated_at`=UNIX_TIMESTAMP(), 
  `deleted_at`=0;

-- 关联：用户-角色
INSERT INTO `admin_user_role` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`)
VALUES 
  (1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (2, 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 关联：角色-权限
INSERT INTO `admin_role_permission` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 菜单列表（ID从1开始连续）
INSERT INTO `admin_menu` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, 0, '仪表盘', '/dashboard', 'Dashboard', 'ele-DataBoard', 2, 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 0, '系统管理', '/system', '', 'ele-Setting', 1, 10, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, 2, '角色管理', '/system/role', 'system/RoleList', 'ele-UserFilled', 2, 11, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 2, '权限管理', '/system/permission', 'system/PermissionList', 'ele-Lock', 2, 12, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 2, '部门管理', '/system/department', 'system/DepartmentList', 'ele-OfficeBuilding', 2, 13, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, 2, '菜单管理', '/system/menu', 'system/MenuList', 'ele-Menu', 2, 14, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 2, '用户管理', '/system/user', 'system/UserList', 'ele-UserFilled', 2, 15, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, 2, '接口管理', '/system/api', 'system/ApiList', 'ele-Connection', 2, 16, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, 0, '临时目录', '/temp', '', 'ele-Folder', 1, 999, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 角色管理按钮（parent_id=3）
  (10, 3, '角色管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (11, 3, '角色管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, 3, '角色管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 权限管理按钮（parent_id=4）
  (13, 4, '权限管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, 4, '权限管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, 4, '权限管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 部门管理按钮（parent_id=5）
  (16, 5, '部门管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, 5, '部门管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (18, 5, '部门管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 菜单管理按钮（parent_id=6）
  (19, 6, '菜单管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (20, 6, '菜单管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (21, 6, '菜单管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 用户管理按钮（parent_id=7）
  (22, 7, '用户管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (23, 7, '用户管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (24, 7, '用户管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 接口管理按钮（parent_id=8）
  (25, 8, '接口管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (26, 8, '接口管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (27, 8, '接口管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 系统配置菜单和按钮（parent_id=2，系统管理下）
  (28, 2, '系统配置', '/system/config', 'system/ConfigList', 'ele-Setting', 2, 17, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (29, 28, '系统配置 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (30, 28, '系统配置 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (31, 28, '系统配置 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典类型菜单和按钮（parent_id=2，系统管理下）
  (32, 2, '字典类型', '/system/dict-type', 'system/DictTypeList', 'ele-Notebook', 2, 18, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (33, 32, '字典类型 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (34, 32, '字典类型 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (35, 32, '字典类型 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典项菜单和按钮（parent_id=2，系统管理下）
  (36, 2, '字典项', '/system/dict-item', 'system/DictItemList', 'ele-List', 2, 19, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (37, 36, '字典项 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (38, 36, '字典项 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (39, 36, '字典项 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 文件管理菜单和按钮（parent_id=2，系统管理下）
  (40, 2, '文件管理', '/system/file', 'system/FileList', 'ele-Folder', 2, 20, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (41, 40, '文件管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (42, 40, '文件管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (43, 40, '文件管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 权限-菜单关联（ID从1开始连续）
INSERT INTO `admin_permission_menu` (`id`, `permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  -- 菜单关联（菜单页面权限）
  -- 角色管理菜单(id=3) -> role:list权限(id=2)
  (1, 2, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 权限管理菜单(id=4) -> permission:list权限(id=6)
  (2, 6, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 部门管理菜单(id=5) -> department:tree权限(id=10)
  (3, 10, 5, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 菜单管理菜单(id=6) -> menu:list权限(id=14)
  (4, 14, 6, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 用户管理菜单(id=7) -> user:list权限(id=18)
  (5, 18, 7, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 接口管理菜单(id=8) -> api:list权限(id=22)
  (6, 22, 8, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  -- 按钮关联（按钮操作权限）
  -- 角色管理按钮
  (7, 3, 10, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- role:create(id=3) -> 角色管理新增按钮(id=10)
  (8, 4, 11, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- role:update(id=4) -> 角色管理编辑按钮(id=11)
  (9, 5, 12, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- role:delete(id=5) -> 角色管理删除按钮(id=12)
  -- 权限管理按钮
  (10, 7, 13, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:create(id=7) -> 权限管理新增按钮(id=13)
  (11, 8, 14, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限管理编辑按钮(id=14)
  (12, 9, 15, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:delete(id=9) -> 权限管理删除按钮(id=15)
  -- 部门管理按钮
  (13, 11, 16, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:create(id=11) -> 部门管理新增按钮(id=16)
  (14, 12, 17, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:update(id=12) -> 部门管理编辑按钮(id=17)
  (15, 13, 18, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:delete(id=13) -> 部门管理删除按钮(id=18)
  -- 菜单管理按钮
  (16, 15, 19, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:create(id=15) -> 菜单管理新增按钮(id=19)
  (17, 16, 20, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:update(id=16) -> 菜单管理编辑按钮(id=20)
  (18, 17, 21, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:delete(id=17) -> 菜单管理删除按钮(id=21)
  -- 用户管理按钮
  (19, 19, 22, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- user:create(id=19) -> 用户管理新增按钮(id=22)
  (20, 20, 23, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- user:update(id=20) -> 用户管理编辑按钮(id=23)
  (21, 21, 24, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- user:delete(id=21) -> 用户管理删除按钮(id=24)
  -- 接口管理按钮
  (22, 23, 25, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:create(id=23) -> 接口管理新增按钮(id=25)
  (23, 24, 26, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:update(id=24) -> 接口管理编辑按钮(id=26)
  (24, 25, 27, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:delete(id=25) -> 接口管理删除按钮(id=27)
  -- 系统配置菜单和按钮
  (25, 26, 28, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:list(id=26) -> 系统配置菜单(id=28)
  (26, 27, 29, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:create(id=27) -> 系统配置新增按钮(id=29)
  (27, 28, 30, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:update(id=28) -> 系统配置编辑按钮(id=30)
  (28, 29, 31, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:delete(id=29) -> 系统配置删除按钮(id=31)
  -- 数据字典类型菜单和按钮
  (29, 30, 32, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:list(id=30) -> 字典类型菜单(id=32)
  (30, 31, 33, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:create(id=31) -> 字典类型新增按钮(id=33)
  (31, 32, 34, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:update(id=32) -> 字典类型编辑按钮(id=34)
  (32, 33, 35, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:delete(id=33) -> 字典类型删除按钮(id=35)
  -- 数据字典项菜单和按钮
  (33, 34, 36, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:list(id=34) -> 字典项菜单(id=36)
  (34, 35, 37, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:create(id=35) -> 字典项新增按钮(id=37)
  (35, 36, 38, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:update(id=36) -> 字典项编辑按钮(id=38)
  (36, 37, 39, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:delete(id=37) -> 字典项删除按钮(id=39)
  -- 文件管理菜单和按钮
  (37, 38, 40, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:list(id=38) -> 文件管理菜单(id=40)
  (38, 39, 41, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:create(id=39) -> 文件管理新增按钮(id=41)
  (39, 40, 42, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:update(id=40) -> 文件管理编辑按钮(id=42)
  (40, 41, 43, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- file:delete(id=41) -> 文件管理删除按钮(id=43)
  -- 群组管理菜单和按钮的权限关联在聊天室模块初始化数据部分使用变量动态关联（见第4节）
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 接口列表（所有业务接口，ID从1开始连续）
INSERT INTO `admin_api` (`id`, `name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 认证相关接口（无需权限，但需要认证）
  (1, '登出', 'POST', '/api/v1/logout', '用户登出', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, '个人信息', 'GET', '/api/v1/profile', '获取当前用户个人信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 用户管理接口
  (3, '用户列表', 'GET', '/api/v1/users', '获取用户列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, '用户新增', 'POST', '/api/v1/users', '新增用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, '用户编辑', 'PUT', '/api/v1/users', '编辑用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, '用户删除', 'DELETE', '/api/v1/users', '删除用户', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, '用户角色列表', 'GET', '/api/v1/users/roles', '获取用户关联的角色列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, '用户角色更新', 'PUT', '/api/v1/users/roles', '更新用户关联的角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 角色管理接口
  (9, '角色列表', 'GET', '/api/v1/roles', '获取角色列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (10, '角色新增', 'POST', '/api/v1/roles', '新增角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (11, '角色编辑', 'PUT', '/api/v1/roles', '编辑角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, '角色删除', 'DELETE', '/api/v1/roles', '删除角色', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (13, '角色权限列表', 'GET', '/api/v1/roles/permissions', '获取角色关联的权限列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, '角色权限更新', 'PUT', '/api/v1/roles/permissions', '更新角色关联的权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 权限管理接口
  (15, '权限列表', 'GET', '/api/v1/permissions', '获取权限列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, '权限新增', 'POST', '/api/v1/permissions', '新增权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, '权限编辑', 'PUT', '/api/v1/permissions', '编辑权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (18, '权限删除', 'DELETE', '/api/v1/permissions', '删除权限', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (19, '权限菜单列表', 'GET', '/api/v1/permissions/menus', '获取权限关联的菜单列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (20, '权限菜单更新', 'PUT', '/api/v1/permissions/menus', '更新权限关联的菜单', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (21, '权限接口列表', 'GET', '/api/v1/permissions/apis', '获取权限关联的接口列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (22, '权限接口更新', 'PUT', '/api/v1/permissions/apis', '更新权限关联的接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 部门管理接口
  (23, '部门树', 'GET', '/api/v1/departments/tree', '获取部门树', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (24, '部门新增', 'POST', '/api/v1/departments', '新增部门', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (25, '部门编辑', 'PUT', '/api/v1/departments', '编辑部门', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (26, '部门删除', 'DELETE', '/api/v1/departments', '删除部门', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 菜单管理接口
  (27, '菜单树', 'GET', '/api/v1/menus/tree', '获取菜单树', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (28, '我的菜单树', 'GET', '/api/v1/menus/my-tree', '获取当前用户的菜单树', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (29, '菜单新增', 'POST', '/api/v1/menus', '新增菜单', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (30, '菜单编辑', 'PUT', '/api/v1/menus', '编辑菜单', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (31, '菜单删除', 'DELETE', '/api/v1/menus', '删除菜单', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 接口管理接口
  (32, '接口列表', 'GET', '/api/v1/apis', '获取接口列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (33, '接口新增', 'POST', '/api/v1/apis', '新增接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (34, '接口编辑', 'PUT', '/api/v1/apis', '编辑接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (35, '接口删除', 'DELETE', '/api/v1/apis', '删除接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 系统配置接口
  (36, '系统配置列表', 'GET', '/api/v1/configs', '获取系统配置列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (37, '系统配置查询', 'GET', '/api/v1/configs/get', '查询系统配置', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (38, '系统配置新增', 'POST', '/api/v1/configs', '新增系统配置', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (39, '系统配置编辑', 'PUT', '/api/v1/configs', '编辑系统配置', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (40, '系统配置删除', 'DELETE', '/api/v1/configs', '删除系统配置', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典类型接口
  (41, '字典类型列表', 'GET', '/api/v1/dict-types', '获取字典类型列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (42, '字典类型新增', 'POST', '/api/v1/dict-types', '新增字典类型', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (43, '字典类型编辑', 'PUT', '/api/v1/dict-types', '编辑字典类型', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (44, '字典类型删除', 'DELETE', '/api/v1/dict-types', '删除字典类型', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 数据字典项接口
  (45, '字典项列表', 'GET', '/api/v1/dict-items', '获取字典项列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (46, '字典项新增', 'POST', '/api/v1/dict-items', '新增字典项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (47, '字典项编辑', 'PUT', '/api/v1/dict-items', '编辑字典项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (48, '字典项删除', 'DELETE', '/api/v1/dict-items', '删除字典项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 字典查询接口（无需权限）
  (49, '字典查询', 'GET', '/api/v1/dict', '根据编码查询字典项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 文件管理接口
  (50, '文件列表', 'GET', '/api/v1/files', '获取文件列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (51, '文件新增', 'POST', '/api/v1/files', '新增文件', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (52, '文件编辑', 'PUT', '/api/v1/files', '编辑文件', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (53, '文件删除', 'DELETE', '/api/v1/files', '删除文件', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (54, '文件上传', 'POST', '/api/v1/files/upload', '上传文件', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (55, '文件下载', 'GET', '/api/v1/files/download', '下载文件', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 个人信息相关接口（只需登录即可访问，不需要权限关联）
  (57, '个人信息更新', 'PUT', '/api/v1/profile', '更新当前用户个人信息（头像、个性签名）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (58, '修改密码', 'POST', '/api/v1/profile/password', '修改当前用户登录密码', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 群组管理接口（需要权限控制）
  (59, '群组列表', 'GET', '/api/v1/chats/groups', '获取群组列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (60, '群组创建', 'POST', '/api/v1/chats/groups', '创建群组', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (61, '群组更新', 'PUT', '/api/v1/chats/groups', '更新群组信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (62, '群组删除', 'DELETE', '/api/v1/chats/groups', '删除群组', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (63, '群组详情', 'GET', '/api/v1/chats/groups/:id', '获取群组详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (64, '群组成员列表', 'GET', '/api/v1/chats/groups/:id/members', '获取群组成员列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (65, '群组成员添加', 'POST', '/api/v1/chats/groups/members', '添加群组成员', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (66, '群组成员移除', 'DELETE', '/api/v1/chats/groups/members', '移除群组成员', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 权限-接口关联（所有权限与接口的关联，ID从1开始连续）
INSERT INTO `admin_permission_api` (`id`, `permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  -- 用户管理权限
  (1, 18, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:list(id=18) -> 用户列表(api_id=3)
  (2, 19, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:create(id=19) -> 用户新增(api_id=4)
  (3, 20, 5, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:update(id=20) -> 用户编辑(api_id=5)
  (4, 21, 6, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:delete(id=21) -> 用户删除(api_id=6)
  (5, 20, 7, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:update(id=20) -> 用户角色列表(api_id=7)
  (6, 20, 8, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- user:update(id=20) -> 用户角色更新(api_id=8)
  -- 角色管理权限
  (7, 2, 9, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),   -- role:list(id=2) -> 角色列表(api_id=9)
  (8, 3, 10, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- role:create(id=3) -> 角色新增(api_id=10)
  (9, 4, 11, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),  -- role:update(id=4) -> 角色编辑(api_id=11)
  (10, 5, 12, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- role:delete(id=5) -> 角色删除(api_id=12)
  (11, 4, 13, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- role:update(id=4) -> 角色权限列表(api_id=13)
  (12, 4, 14, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- role:update(id=4) -> 角色权限更新(api_id=14)
  -- 权限管理权限
  (13, 6, 15, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:list(id=6) -> 权限列表(api_id=15)
  (14, 7, 16, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:create(id=7) -> 权限新增(api_id=16)
  (15, 8, 17, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限编辑(api_id=17)
  (16, 9, 18, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:delete(id=9) -> 权限删除(api_id=18)
  (17, 8, 19, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限菜单列表(api_id=19)
  (18, 8, 20, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限菜单更新(api_id=20)
  (19, 8, 21, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限接口列表(api_id=21)
  (20, 8, 22, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- permission:update(id=8) -> 权限接口更新(api_id=22)
  -- 部门管理权限
  (21, 10, 23, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:tree(id=10) -> 部门树(api_id=23)
  (22, 11, 24, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:create(id=11) -> 部门新增(api_id=24)
  (23, 12, 25, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:update(id=12) -> 部门编辑(api_id=25)
  (24, 13, 26, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- department:delete(id=13) -> 部门删除(api_id=26)
  -- 菜单管理权限
  (25, 14, 27, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:list(id=14) -> 菜单树(api_id=27)
  (26, 15, 29, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:create(id=15) -> 菜单新增(api_id=29)
  (28, 16, 30, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:update(id=16) -> 菜单编辑(api_id=30)
  (29, 17, 31, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- menu:delete(id=17) -> 菜单删除(api_id=31)
  -- 接口管理权限
  (30, 22, 32, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:list(id=22) -> 接口列表(api_id=32)
  (31, 23, 33, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:create(id=23) -> 接口新增(api_id=33)
  (32, 24, 34, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:update(id=24) -> 接口编辑(api_id=34)
  (33, 25, 35, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- api:delete(id=25) -> 接口删除(api_id=35)
  -- 系统配置权限
  (34, 26, 36, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:list(id=26) -> 系统配置列表(api_id=36)
  (35, 26, 37, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:list(id=26) -> 系统配置查询(api_id=37)
  (36, 27, 38, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:create(id=27) -> 系统配置新增(api_id=38)
  (37, 28, 39, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:update(id=28) -> 系统配置编辑(api_id=39)
  (38, 29, 40, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- config:delete(id=29) -> 系统配置删除(api_id=40)
  -- 数据字典类型权限
  (39, 30, 41, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:list(id=30) -> 字典类型列表(api_id=41)
  (40, 31, 42, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:create(id=31) -> 字典类型新增(api_id=42)
  (41, 32, 43, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:update(id=32) -> 字典类型编辑(api_id=43)
  (42, 33, 44, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_type:delete(id=33) -> 字典类型删除(api_id=44)
  -- 数据字典项权限
  (43, 34, 45, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:list(id=34) -> 字典项列表(api_id=45)
  (44, 35, 46, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:create(id=35) -> 字典项新增(api_id=46)
  (45, 36, 47, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:update(id=36) -> 字典项编辑(api_id=47)
  (46, 37, 48, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- dict_item:delete(id=37) -> 字典项删除(api_id=48)
  -- 文件管理权限
  (47, 38, 50, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:list(id=38) -> 文件列表(api_id=50)
  (48, 39, 51, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:create(id=39) -> 文件新增(api_id=51)
  (49, 40, 52, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:update(id=40) -> 文件编辑(api_id=52)
  (50, 41, 53, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- file:delete(id=41) -> 文件删除(api_id=53)
  -- 群组管理权限
  (51, 45, 59, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:detail(id=45) -> 群组列表(api_id=59)
  (52, 42, 60, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:create(id=42) -> 群组创建(api_id=60)
  (53, 43, 61, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:update(id=43) -> 群组更新(api_id=61)
  (54, 44, 62, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:delete(id=44) -> 群组删除(api_id=62)
  (55, 45, 63, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:detail(id=45) -> 群组详情(api_id=63)
  (56, 46, 64, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:member(id=46) -> 群组成员列表(api_id=64)
  (57, 46, 65, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- chat:group:member(id=46) -> 群组成员添加(api_id=65)
  (58, 46, 66, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- chat:group:member(id=46) -> 群组成员移除(api_id=66)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 2. 其他初始化业务数据（配置、字典等）
-- ============================================

-- 系统配置初始化数据
INSERT INTO `admin_config` (`id`, `group`, `key`, `value`, `type`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, 'system', 'system:app_name', '"后台管理系统"', 'string', '应用名称', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 'system', 'system:app_logo', '"/static/logo.png"', 'string', '应用Logo路径', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, 'system', 'system:app_version', '"1.0.0"', 'string', '应用版本', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 'system', 'system:timeout', '300', 'number', '会话超时时间（秒）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 'theme', 'theme:primary_color', '"#409EFF"', 'string', '主题主色', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, 'theme', 'theme:sidebar_width', '200', 'number', '侧边栏宽度（px）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 'upload', 'upload:max_size', '10485760', 'number', '最大上传文件大小（字节，默认10MB）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, 'upload', 'upload:allowed_types', '["jpg","jpeg","png","gif","pdf","doc","docx","xls","xlsx"]', 'json', '允许上传的文件类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `value`=VALUES(`value`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 数据字典类型初始化数据
INSERT INTO `admin_dict_type` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, '用户状态', 'user_status', '用户账号状态字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, '性别', 'gender', '性别字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, '是否', 'yes_no', '是否字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, '文件存储类型', 'file_storage_type', '文件存储类型字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, '聊天配置', 'chat_config', '在线聊天相关配置字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, '消息来源类型', 'notification_source_type', '消息通知的来源类型字典（chat、notice、system等）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- 数据字典项初始化数据
INSERT INTO `admin_dict_item` (`id`, `type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 用户状态字典项（枚举从 1 开始，0 预留为「全部/不筛选」）
  (1, 1, '启用', '1', 1, 1, '用户账号启用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 1, '禁用', '2', 2, 1, '用户账号禁用状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 性别字典项（枚举从 1 开始，0 不再使用）
  (3, 2, '男', '1', 1, 1, '男性', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 2, '女', '2', 2, 1, '女性', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 2, '未知', '3', 3, 1, '未知性别', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 是否字典项（枚举从 1 开始，0 不再使用）
  (6, 3, '是', '1', 1, 1, '是', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 3, '否', '2', 2, 1, '否', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 文件存储类型字典项
  (8, 4, '本地存储', 'local', 1, 1, '本地文件系统存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, 4, 'OSS存储', 'oss', 2, 1, '阿里云OSS存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (10, 4, 'S3存储', 's3', 3, 1, 'AWS S3存储', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 聊天配置字典项
  (11, 5, '聊天窗口消息数量', '30', 1, 1, '每个聊天窗口显示的最新消息数量', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, 5, '在线聊天页面路径', '/chatroom/chat', 2, 1, '在线聊天页面的前端路由路径，用于消息通知跳转', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, 5, 'Emoji每行显示数量', '8', 3, 1, 'Emoji表情选择器每行显示的表情数量（x）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, 5, 'Emoji显示行数', '3', 4, 1, 'Emoji表情选择器显示的行数（y）', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  -- 消息来源类型字典项
  (13, 6, '在线聊天', 'chat', 1, 1, '在线聊天消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, 6, '系统公告', 'notice', 2, 1, '系统公告消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, 6, '系统通知', 'system', 3, 1, '系统通知消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;

-- ============================================
-- 2.1 初始化公告数据
-- ============================================
-- 插入一条欢迎公告（已发布状态）
INSERT INTO `admin_notice` (`id`, `title`, `content`, `type`, `status`, `publish_time`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  1,
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  1, -- 类型：1 普通公告
  2, -- 状态：2 已发布
  UNIX_TIMESTAMP(), -- 发布时间：当前时间
  1, -- 创建人：超级管理员（id=1）
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE 
  `title`=VALUES(`title`),
  `content`=VALUES(`content`),
  `type`=VALUES(`type`),
  `status`=VALUES(`status`),
  `publish_time`=VALUES(`publish_time`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 2.2 初始化聊天数据
-- ============================================
-- 创建默认企业群组（id=1，type=2群组，created_by=1超级管理员）
INSERT INTO `chat` (`id`, `name`, `type`, `avatar`, `description`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  1,
  '企业群组',
  2, -- 类型：2 群组
  '',
  '默认企业群组，所有用户自动加入',
  1, -- 创建人：超级管理员（id=1）
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`), 
  `type`=VALUES(`type`), 
  `description`=VALUES(`description`), 
  `updated_at`=UNIX_TIMESTAMP(), 
  `deleted_at`=0;

-- 将初始的两个用户（id=1和id=2）都加入默认企业群组
INSERT INTO `chat_user` (`id`, `chat_id`, `user_id`, `joined_at`, `created_at`, `updated_at`)
VALUES 
  (1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- 超级管理员加入群组
  (2, 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- 业务管理员加入群组
ON DUPLICATE KEY UPDATE 
  `joined_at`=VALUES(`joined_at`), 
  `updated_at`=UNIX_TIMESTAMP();

-- 创建用户1和用户2之间的私聊（id=2，type=1私聊）
INSERT INTO `chat` (`id`, `name`, `type`, `avatar`, `description`, `created_by`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  2,
  '', -- 私聊名称为空，前端根据对方用户信息显示
  1, -- 类型：1 私聊
  '',
  '',
  0, -- 私聊创建人为0
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
)
ON DUPLICATE KEY UPDATE 
  `type`=VALUES(`type`), 
  `updated_at`=UNIX_TIMESTAMP(), 
  `deleted_at`=0;

-- 将用户1和用户2加入私聊
INSERT INTO `chat_user` (`id`, `chat_id`, `user_id`, `joined_at`, `created_at`, `updated_at`)
VALUES 
  (3, 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()), -- 用户1加入私聊
  (4, 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP())  -- 用户2加入私聊
ON DUPLICATE KEY UPDATE 
  `joined_at`=VALUES(`joined_at`), 
  `updated_at`=UNIX_TIMESTAMP();

-- 为初始化用户创建公告通知（公告已发布，需要给所有用户创建通知）
-- 注意：这里只给初始化时已存在的用户创建通知，后续新增的用户会在登录时自动获取未读公告
INSERT INTO `admin_notification` (`user_id`, `source_type`, `source_id`, `title`, `content`, `read_status`, `read_at`, `created_at`, `updated_at`, `deleted_at`)
SELECT 
  1, -- 超级管理员（id=1）
  'notice',
  1, -- 公告ID
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  0, -- 未读
  0,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_notification` 
  WHERE `user_id` = 1 AND `source_type` = 'notice' AND `source_id` = 1 AND `deleted_at` = 0
);

INSERT INTO `admin_notification` (`user_id`, `source_type`, `source_id`, `title`, `content`, `read_status`, `read_at`, `created_at`, `updated_at`, `deleted_at`)
SELECT 
  2, -- admin业务管理员（id=2）
  'notice',
  1, -- 公告ID
  '欢迎使用后台管理系统',
  '欢迎使用本后台管理系统！\n\n本系统提供了完整的权限管理、用户管理、角色管理、菜单管理、接口管理等基础功能，以及聊天室、公告管理、消息通知等业务功能。\n\n祝您使用愉快！',
  0, -- 未读
  0,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
WHERE NOT EXISTS (
  SELECT 1 FROM `admin_notification` 
  WHERE `user_id` = 2 AND `source_type` = 'notice' AND `source_id` = 1 AND `deleted_at` = 0
);

-- ============================================
-- 3. 日志与监控相关模块初始化数据（归类到系统管理）
-- ============================================

-- 获取系统管理目录菜单ID（path = '/system'）
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system' AND `deleted_at` = 0 LIMIT 1);

-- ==========================
-- 3.1 操作日志模块
-- ==========================
-- 操作日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '操作日志',
    '/system/operation-log',
    'system/OperationLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    30, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @operation_menu_id = (
  SELECT `id` FROM `admin_menu` 
  WHERE `path` = '/system/operation-log' AND `deleted_at` = 0 
  LIMIT 1
);

-- 操作日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @operation_menu_id,
    '操作日志 导出按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否（按钮不显示在菜单中）
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @operation_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @operation_menu_id 
    AND `name` = '操作日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('操作日志列表', 'operation_log:list', '查看操作日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志详情', 'operation_log:detail', '查看操作日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志导出', 'operation_log:export', '导出操作日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @operation_list_permission_id = (
  SELECT `id` FROM `admin_permission` 
  WHERE `code` = 'operation_log:list' AND `deleted_at` = 0 
  LIMIT 1
);
SET @operation_detail_permission_id = (
  SELECT `id` FROM `admin_permission` 
  WHERE `code` = 'operation_log:detail' AND `deleted_at` = 0 
  LIMIT 1
);
SET @operation_export_permission_id = (
  SELECT `id` FROM `admin_permission` 
  WHERE `code` = 'operation_log:export' AND `deleted_at` = 0 
  LIMIT 1
);

-- 操作日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('操作日志列表', 'GET', '/api/v1/operation-logs', '获取操作日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志详情', 'GET', '/api/v1/operation-logs/detail', '获取操作日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作日志导出', 'GET', '/api/v1/operation-logs/export', '导出操作日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @operation_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @operation_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/operation-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 操作日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@operation_list_permission_id, @operation_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_export_permission_id, @operation_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 操作日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@operation_list_permission_id, @operation_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_detail_permission_id, @operation_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@operation_export_permission_id, @operation_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联操作日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @operation_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @operation_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @operation_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ==========================
-- 3.2 登录日志模块
-- ==========================
-- 登录日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '登录日志',
    '/system/login-log',
    'system/LoginLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    31, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/login-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志详情按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @login_menu_id,
    '登录日志 详情按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_detail_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @login_menu_id 
    AND `name` = '登录日志 详情按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @login_menu_id,
    '登录日志 导出按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    2, -- 排序值
    0, -- 是否可见：0 否
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @login_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @login_menu_id 
    AND `name` = '登录日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('登录日志列表', 'login_log:list', '查看登录日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志详情', 'login_log:detail', '查看登录日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志导出', 'login_log:export', '导出登录日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @login_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:list' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_detail_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_export_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'login_log:export' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('登录日志列表', 'GET', '/api/v1/login-logs', '获取登录日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志详情', 'GET', '/api/v1/login-logs/detail', '获取登录日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录日志导出', 'GET', '/api/v1/login-logs/export', '导出登录日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @login_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @login_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/login-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 登录日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@login_list_permission_id, @login_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_detail_permission_id, @login_detail_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_export_permission_id, @login_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 登录日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@login_list_permission_id, @login_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_detail_permission_id, @login_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@login_export_permission_id, @login_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联登录日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @login_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @login_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @login_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ==========================
-- 3.3 审计日志模块
-- ==========================
-- 审计日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '审计日志',
    '/system/audit-log',
    'system/AuditLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    32, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @audit_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/audit-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志导出按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @audit_menu_id,
    '审计日志 导出按钮',
    '',
    '',
    '',
    3, -- 类型：3 按钮
    1, -- 排序值
    0, -- 是否可见：0 否（按钮不显示在菜单中）
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @audit_export_button_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `parent_id` = @audit_menu_id 
    AND `name` = '审计日志 导出按钮'
    AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('审计日志列表', 'audit_log:list', '查看审计日志列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志详情', 'audit_log:detail', '查看审计日志详情', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志导出', 'audit_log:export', '导出审计日志', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @audit_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:list' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_detail_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_export_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'audit_log:export' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('审计日志列表', 'GET', '/api/v1/audit-logs', '获取审计日志列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志详情', 'GET', '/api/v1/audit-logs/detail', '获取审计日志详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计日志导出', 'GET', '/api/v1/audit-logs/export', '导出审计日志', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @audit_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_detail_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs/detail' AND `deleted_at` = 0
  LIMIT 1
);
SET @audit_export_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/audit-logs/export' AND `deleted_at` = 0
  LIMIT 1
);

-- 审计日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@audit_list_permission_id, @audit_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_export_permission_id, @audit_export_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 审计日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@audit_list_permission_id, @audit_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_detail_permission_id, @audit_detail_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@audit_export_permission_id, @audit_export_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联审计日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @audit_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @audit_detail_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @audit_export_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ==========================
-- 3.4 性能监控日志模块
-- ==========================
-- 性能监控日志主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '性能监控日志',
    '/system/performance-log',
    'system/PerformanceLogList',
    'ele-Document',
    2, -- 类型：2 菜单
    33, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @performance_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/performance-log' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志列表权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '性能监控日志列表',
    'performance_log:list',
    '查看性能监控日志列表',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @performance_list_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'performance_log:list' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志列表接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '性能监控日志列表',
    'GET',
    '/api/v1/performance-logs',
    '获取性能监控日志列表',
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @performance_list_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/performance-logs' AND `deleted_at` = 0
  LIMIT 1
);

-- 性能监控日志 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@performance_list_permission_id, @performance_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 性能监控日志 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES (@performance_list_permission_id, @performance_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联性能监控日志权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES (1, @performance_list_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ==========================
-- 3.5 系统监控模块
-- ==========================
-- 系统监控主菜单（系统管理下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @system_menu_id,
    '系统监控',
    '/system/monitor',
    'system/MonitorList',
    'ele-Monitor',
    2, -- 类型：2 菜单
    34, -- 排序值
    1, -- 是否可见：1 是
    1, -- 状态：1 启用
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
    `name`=VALUES(`name`),
    `component`=VALUES(`component`),
    `icon`=VALUES(`icon`),
    `type`=VALUES(`type`),
    `order_num`=VALUES(`order_num`),
    `visible`=VALUES(`visible`),
    `status`=VALUES(`status`),
    `updated_at`=UNIX_TIMESTAMP(),
    `deleted_at`=0;
SET @monitor_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/system/monitor' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控查看权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    '系统监控查看',
    'monitor:view',
    '查看系统监控信息',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @monitor_view_permission_id = (
  SELECT `id` FROM `admin_permission`
  WHERE `code` = 'monitor:view' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('系统监控状态', 'GET', '/api/v1/monitor/status', '获取系统资源使用情况（CPU、内存、磁盘、网络）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('系统统计', 'GET', '/api/v1/monitor/stats', '获取系统统计数据（用户数、角色数、权限数等）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `status`=VALUES(`status`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
SET @monitor_status_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/monitor/status' AND `deleted_at` = 0
  LIMIT 1
);
SET @monitor_stats_api_id = (
  SELECT `id` FROM `admin_api`
  WHERE `method` = 'GET' AND `path` = '/api/v1/monitor/stats' AND `deleted_at` = 0
  LIMIT 1
);

-- 系统监控 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES (@monitor_view_permission_id, @monitor_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 系统监控 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@monitor_view_permission_id, @monitor_status_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@monitor_view_permission_id, @monitor_stats_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 超级管理员角色关联系统监控权限（role_id = 1）
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES (1, @monitor_view_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ============================================
-- 4. 聊天室模块初始化数据
-- ============================================
-- 聊天室目录
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (0, '聊天室', '/chatroom', '', 'ele-ChatDotRound', 1, 20, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chatroom_dir_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/chatroom' AND `deleted_at` = 0 LIMIT 1);

-- 在线聊天菜单（无需权限，只要登录就可以访问）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '在线聊天', '/chatroom/chat', 'chatroom/ChatList', 'ele-ChatLineRound', 2, 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 聊天记录管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '聊天记录管理', '/chatroom/chat-message', 'chatroom/ChatMessageList', 'ele-Document', 2, 2, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/chatroom/chat-message' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理删除按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chat_message_menu_id, '聊天记录管理 删除按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_message_menu_id AND `name` = '聊天记录管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('聊天记录列表', 'chat_message:list', '查看聊天记录列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('聊天记录删除', 'chat_message:delete', '删除聊天记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat_message:list' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat_message:delete' AND `deleted_at` = 0 LIMIT 1);

-- 在线聊天接口（无需权限，只需要AuthMiddleware）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('聊天消息发送', 'POST', '/api/v1/chats/messages', '发送聊天消息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('可聊天用户列表', 'GET', '/api/v1/chats/users', '获取可聊天用户列表（包含部门-角色-昵称信息）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 聊天记录管理接口（需要权限）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('聊天记录列表', 'GET', '/api/v1/chats/messages', '获取聊天记录列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('聊天记录删除', 'DELETE', '/api/v1/chats/messages', '删除聊天记录', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_message_list_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/chats/messages' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/chats/messages' AND `deleted_at` = 0 LIMIT 1);

-- 聊天记录管理 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@chat_message_list_permission_id, @chat_message_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_message_delete_permission_id, @chat_message_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 群组管理菜单（在聊天室目录下）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@chatroom_dir_id, '群组管理', '/chatroom/chat-group', 'chatroom/ChatGroupList', 'ele-ChatDotRound', 2, 3, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_group_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/chatroom/chat-group' AND `deleted_at` = 0 LIMIT 1);

-- 群组管理按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (@chat_group_menu_id, '群组管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_group_menu_id, '群组管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_group_menu_id, '群组管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @chat_group_create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 新增按钮' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 编辑按钮' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @chat_group_menu_id AND `name` = '群组管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 群组管理 权限-菜单 关联
SET @chat_group_detail_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:detail' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:create' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:update' AND `deleted_at` = 0 LIMIT 1);
SET @chat_group_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'chat:group:delete' AND `deleted_at` = 0 LIMIT 1);
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@chat_group_detail_permission_id, @chat_group_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_create_permission_id, @chat_group_create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_update_permission_id, @chat_group_update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_group_delete_permission_id, @chat_group_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- 聊天记录管理 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@chat_message_list_permission_id, @chat_message_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@chat_message_delete_permission_id, @chat_message_delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP();

-- ============================================
-- 5. 公告管理模块初始化数据
-- ============================================
SET @system_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@system_menu_id, '公告管理', '/system/notice', 'system/NoticeList', 'ele-Document', 2, 21, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/system/notice' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理按钮
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (@notice_menu_id, '公告管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_menu_id, '公告管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_menu_id, '公告管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 新增按钮' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 编辑按钮' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @notice_menu_id AND `name` = '公告管理 删除按钮' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理权限
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('公告管理列表', 'notice:list', '查看公告管理列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理新增', 'notice:create', '新增公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理编辑', 'notice:update', '编辑公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理删除', 'notice:delete', '删除公告管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:list' AND `deleted_at` = 0 LIMIT 1);
SET @notice_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:create' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:update' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'notice:delete' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理接口
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('公告管理列表', 'GET', '/api/v1/notices', '获取公告管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理新增', 'POST', '/api/v1/notices', '新增公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理编辑', 'PUT', '/api/v1/notices', '编辑公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告管理删除', 'DELETE', '/api/v1/notices', '删除公告管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;
SET @notice_list_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_create_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_update_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'PUT' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);
SET @notice_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/notices' AND `deleted_at` = 0 LIMIT 1);

-- 公告管理 权限-菜单 关联
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES 
  (@notice_list_permission_id, @notice_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_create_permission_id, @notice_create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_update_permission_id, @notice_update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_delete_permission_id, @notice_delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 公告管理 权限-接口 关联
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES 
  (@notice_list_permission_id, @notice_list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_create_permission_id, @notice_create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_update_permission_id, @notice_update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@notice_delete_permission_id, @notice_delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 6. 消息通知管理模块初始化数据
-- ============================================
-- 注意：消息通知管理只要登录就有权限，不需要权限控制

-- 消息通知管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (@system_menu_id, '消息通知管理', '/system/notification', 'system/NotificationList', 'ele-Bell', 2, 22, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`), `name`=VALUES(`name`), `path`=VALUES(`path`), `component`=VALUES(`component`), `icon`=VALUES(`icon`), `type`=VALUES(`type`), `order_num`=VALUES(`order_num`), `visible`=VALUES(`visible`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- 消息通知管理接口（只需要AuthMiddleware，不需要PermissionMiddleware）
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ('消息通知管理列表', 'GET', '/api/v1/notifications', '获取消息通知管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知全部已读', 'PUT', '/api/v1/notifications/read-all', '标记所有消息通知为已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知全部已读', 'PUT', '/api/v1/notifications/read-all', '标记所有消息通知为已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知标记已读', 'PUT', '/api/v1/notifications/read', '标记单个消息通知为已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知清除已读', 'DELETE', '/api/v1/notifications/read', '清除所有已读消息通知', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息通知删除', 'DELETE', '/api/v1/notifications', '删除消息通知', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- ============================================
-- 8. 字典数据（从迁移文件合并）
-- ============================================

-- 8.1 SDK 状态字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK状态', 'sdk_status', 'SDK Key 状态（启用/禁用）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_status' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_status_type_id, '启用', '1', 1, 1, '可用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_status_type_id, '禁用', '2', 2, 1, '停用', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.2 SDK 默认限频字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK默认限频', 'sdk_rate_limit_default', 'SDK 接口默认限频（次/分钟），可被接口自定义值覆盖', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_rate_limit_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_rate_limit_default' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_rate_limit_type_id, '默认60次/分钟', '60', 1, 1, '默认限频上限，单位：次/分钟', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.3 SDK HTTP 方法
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK HTTP 方法', 'sdk_http_method', 'SDK 接口支持的 HTTP 方法', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @sdk_http_method_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'sdk_http_method' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_http_method_type_id, 'GET', 'GET', 1, 1, 'HTTP GET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'POST', 'POST', 2, 1, 'HTTP POST', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'PUT', 'PUT', 3, 1, 'HTTP PUT', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_http_method_type_id, 'DELETE', 'DELETE', 4, 1, 'HTTP DELETE', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.4 本地存储配置字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('本地存储', 'storage_base_url', '本地存储配置，用于配置文件存储的baseURL', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'storage_base_url' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认存储地址', 'https://oldbai.top/oss', 1, 1, '文件存储的baseURL，用于生成文件完整访问路径', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.5 系统选项字典（性能日志慢查询状态、菜单类型、已读状态、操作类型、HTTP请求方法、公告类型、公告状态、登录状态、审计类型、消息类型、短句类型）
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('性能日志慢查询状态', 'performance_log_slow_status', '性能日志慢查询状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('菜单类型', 'menu_type', '菜单类型选项：目录、菜单、按钮', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('已读状态', 'read_status', '已读状态选项：未读、已读', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('操作类型', 'operation_type', '操作日志操作类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('HTTP请求方法', 'http_method', 'HTTP请求方法选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告类型', 'notice_type', '公告类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公告状态', 'notice_status', '公告状态选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('登录状态', 'login_status', '登录状态选项：成功、失败', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('审计类型', 'audit_type', '审计日志审计类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('消息类型', 'chat_message_type', '聊天消息类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('短句类型', 'daily_short_sentence_type', '每日短句类型选项', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @performance_log_slow_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'performance_log_slow_status' AND `deleted_at` = 0 LIMIT 1);
SET @menu_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'menu_type' AND `deleted_at` = 0 LIMIT 1);
SET @read_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'read_status' AND `deleted_at` = 0 LIMIT 1);
SET @operation_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'operation_type' AND `deleted_at` = 0 LIMIT 1);
SET @http_method_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'http_method' AND `deleted_at` = 0 LIMIT 1);
SET @notice_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_type' AND `deleted_at` = 0 LIMIT 1);
SET @notice_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'notice_status' AND `deleted_at` = 0 LIMIT 1);
SET @login_status_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'login_status' AND `deleted_at` = 0 LIMIT 1);
SET @audit_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'audit_type' AND `deleted_at` = 0 LIMIT 1);
SET @chat_message_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'chat_message_type' AND `deleted_at` = 0 LIMIT 1);
SET @daily_short_sentence_type_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'daily_short_sentence_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@performance_log_slow_status_type_id, 'Normal', '2', 1, 1, '正常查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@performance_log_slow_status_type_id, 'Slow', '1', 2, 1, '慢查询', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '目录', '1', 1, 1, '目录类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '菜单', '2', 2, 1, '菜单类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@menu_type_type_id, '按钮', '3', 3, 1, '按钮类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@read_status_type_id, '未读', '1', 1, 1, '未读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@read_status_type_id, '已读', '2', 2, 1, '已读状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '创建', 'create', 1, 1, '创建操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '更新', 'update', 2, 1, '更新操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '删除', 'delete', 3, 1, '删除操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '查询', 'query', 4, 1, '查询操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@operation_type_type_id, '导出', 'export', 5, 1, '导出操作', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'GET', 'GET', 1, 1, 'GET请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'POST', 'POST', 2, 1, 'POST请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'PUT', 'PUT', 3, 1, 'PUT请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@http_method_type_id, 'DELETE', 'DELETE', 4, 1, 'DELETE请求方法', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '普通公告', '1', 1, 1, '普通公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '重要公告', '2', 2, 1, '重要公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_type_type_id, '紧急公告', '3', 3, 1, '紧急公告类型', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_status_type_id, '草稿', '1', 1, 1, '草稿状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@notice_status_type_id, '已发布', '2', 2, 1, '已发布状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@login_status_type_id, '失败', '2', 1, 1, '登录失败', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@login_status_type_id, '成功', '1', 2, 1, '登录成功', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '权限分配', 'permission_assign', 1, 1, '权限分配审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '角色变更', 'role_change', 2, 1, '角色变更审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '配置修改', 'config_modify', 3, 1, '配置修改审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@audit_type_type_id, '数据删除', 'data_delete', 4, 1, '数据删除审计', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '文本', '1', 1, 1, '文本消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '图片', '2', 2, 1, '图片消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@chat_message_type_type_id, '文件', '3', 3, 1, '文件消息', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@daily_short_sentence_type_type_id, '普通', '1', 1, 1, '普通短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@daily_short_sentence_type_type_id, '文学', '2', 2, 1, '文学短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.6 视频代理地址字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频代理地址', 'video_proxy_url', '视频代理服务器地址配置，用于代理m3u8等视频流', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'video_proxy_url' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认代理', 'http://localhost:8888/api/v1/videos/proxy', 1, 1, '默认视频代理服务器地址', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.7 视频来源类型字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频来源类型', 'video_source_type', '视频来源类型字典，用于区分手动添加和采集的视频', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'video_source_type' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '手动添加', '1', 1, 1, '手动添加的视频', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@dict_type_id, '采集', '2', 2, 1, '通过采集接口添加的视频', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 8.8 WebSocket连接配置字典
INSERT INTO `admin_dict_type` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('WebSocket连接', 'websocket_base_url', 'WebSocket连接配置，用于配置WebSocket的baseURL', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @dict_type_id = (SELECT `id` FROM `admin_dict_type` WHERE `code` = 'websocket_base_url' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_dict_item` (`type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@dict_type_id, '默认WebSocket地址', 'oldbai.top/ws', 1, 1, 'WebSocket连接的baseURL，生产环境使用 wss://oldbai.top/ws，开发环境使用 ws://localhost:20000', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `label`=VALUES(`label`),
  `value`=VALUES(`value`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `remark`=VALUES(`remark`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- ============================================
-- 9. 模块初始化数据（从迁移文件合并）
-- ============================================

-- 9.1 每日短句模块初始化（菜单、权限、接口）
SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/temp' AND `deleted_at` = 0 LIMIT 1),
  (SELECT `id` FROM `admin_menu` WHERE `id` = 9 AND `deleted_at` = 0 LIMIT 1)
);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '每日短句',
    '/temp/daily_short_sentence',
    'temp/DailyShortSentenceList',
    'ele-Document',
    2,
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/temp/daily_short_sentence' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@main_menu_id, '每日短句 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '每日短句 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '每日短句 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 新增按钮' AND `deleted_at`=0 LIMIT 1);
SET @update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 编辑按钮' AND `deleted_at`=0 LIMIT 1);
SET @delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 删除按钮' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('每日短句列表', 'daily_short_sentence:list', '查看每日短句列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句新增', 'daily_short_sentence:create', '新增每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句编辑', 'daily_short_sentence:update', '编辑每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句删除', 'daily_short_sentence:delete', '删除每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:list' AND `deleted_at`=0 LIMIT 1);
SET @create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:create' AND `deleted_at`=0 LIMIT 1);
SET @update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:update' AND `deleted_at`=0 LIMIT 1);
SET @delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('每日短句列表', 'GET', '/api/v1/daily-short-sentences', '获取每日短句列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句新增', 'POST', '/api/v1/daily-short-sentences', '新增每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句编辑', 'PUT', '/api/v1/daily-short-sentences', '编辑每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句删除', 'DELETE', '/api/v1/daily-short-sentences', '删除每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='GET' AND `deleted_at`=0 LIMIT 1);
SET @create_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='POST' AND `deleted_at`=0 LIMIT 1);
SET @update_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='PUT' AND `deleted_at`=0 LIMIT 1);
SET @delete_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='DELETE' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 9.2 每日短句初始化数据
INSERT INTO `daily_short_sentence` (`id`, `type`, `content`, `img`, `literature_author`, `convert_img`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, 2, '我愿做一颗永不生锈的螺丝钉。', 'https://t10.baidu.com/it/u=2607170580,4056796944&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 2, '母爱胜于万爱。', 'https://t10.baidu.com/it/u=3325666458,3828073077&fm=58', '莎士比亚', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, 2, 'If you shed tears when you miss the sun,you also miss the stars. 如果你因为失去太阳而落泪，那么你也将失去群星。', 'https://t12.baidu.com/it/u=3583036367,4054301455&fm=58', '泰戈尔', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 2, '没有在深夜痛哭过的人，不足以谈人生。', 'https://t10.baidu.com/it/u=3894103860,4159876305&fm=58', '托马斯·卡莱尔', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 2, '那一世，你为蝴蝶，我为落花，花心已碎，蝶翼天涯，那一世，你为繁星，我为月牙，形影相错，空负年华，那一世，你为歌女，我为琵琶，乱世笙歌，深情天下，金戈铁马，水月镜花，容华一刹那，那缕传世的青烟，点缀着你我结缘的童话。不问贵贱，不顾浮华，三千华发，一生牵挂。', 'https://t10.baidu.com/it/u=2066699184,3350866713&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, 2, '鲜花打扮不出美丽的春天，一个人先进总是单枪匹马，众人先进才能移山填海。', 'https://t11.baidu.com/it/u=3681108379,232235240&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 2, '劳动是一切知识的源泉。', 'https://t11.baidu.com/it/u=3277121875,3598099677&fm=58', '陶铸', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, 2, '一朵鲜花打扮不出美丽的春天，一个人先进总是单一槍一匹马，众人先进才能移山填海。', 'https://t11.baidu.com/it/u=1167556189,1973262450&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, 2, '强悍太久，让我软弱很难。', 'https://t11.baidu.com/it/u=3439716001,3622666293&fm=58', '赵丽颖', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (10, 2, '信仰是石，擦起星星之火；信仰是火，点亮希望之灯；信仰是灯，照亮夜行的路；信仰是路，引你走向黎明。', 'https://t12.baidu.com/it/u=546752963,682448945&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (11, 2, '站立着的心，只有努力努力再努力。', 'https://t11.baidu.com/it/u=3748327216,3933540256&fm=58', '张艺兴', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, 2, '读书使人得到一种优雅和风味，这就是读书的整个目的，而只有抱着这种目的的读书才可以叫做艺术，一人读书的目的并不是要"改进心智"，因为当他开始想要改进心智的时候，一切读书的乐趣便丧失净尽了。', 'https://t11.baidu.com/it/u=4257120178,225048982&fm=58', '林语堂', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (13, 2, '去见你想见的人吧。趁阳光正好，趁微风不噪，趁繁花还未开至荼蘼，趁现在还年轻，还可以走很长很长的路，还能诉说很深很深的思念，趁世界还不那么拥挤，趁飞机还没有起飞，趁现在自己的双手还能拥抱彼此，趁我们还有呼吸。', 'https://t11.baidu.com/it/u=1093319152,1457586368&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, 2, '孤单不是与生俱来，而是由你爱上一个人的那一刻开始。', 'https://t12.baidu.com/it/u=3224252345,3273330929&fm=58', '张小娴', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, 2, '我谈过最长的恋爱，就是自恋，我爱自己，没有情敌。', 'https://t12.baidu.com/it/u=289577954,149653588&fm=58', '安东尼', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, 2, '上帝创造了整数，所有其余的数都是人造的。', 'https://t11.baidu.com/it/u=4055619370,68771096&fm=58', 'L·克隆内克', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, 2, '帝王：待我君临天下，许你四海为家；国臣：待我了无牵挂，许你浪迹天涯；将军：待我半生戎马，许你共话桑麻；书生：待我功成名达，许你花前月下；侠客：待我名满华夏，许你放歌纵马；琴师：待我弦断音垮，许你青丝白发；面首：待我不再有她，许你淡饭粗茶；情郎：待我高头大马，许你嫁衣红霞；农夫：待我富贵荣华，许你十里桃花；僧人：待我一袭袈裟，许你相思放下。', 'https://t10.baidu.com/it/u=2524492674,127412282&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (18, 2, '我是广大劳苦大众当中的一员，我能帮忙人民克服一点困难，是最幸福的。', 'https://t11.baidu.com/it/u=3991612066,1002782314&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (19, 2, '雨下给富人，也下给穷人，下给义人，也下给不义的人；其实，雨并不公道，因为下落在一个没有公道的世界上。', 'https://t10.baidu.com/it/u=578245077,202302703&fm=58', '老舍', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (20, 2, '只要能培一朵花，就不妨做做会朽的腐草。', 'https://t12.baidu.com/it/u=472244129,624373840&fm=58', '鲁迅', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (21, 2, '当你做对的时候，没有人会记得；当你做错的时候，连呼吸都是错。', 'https://t12.baidu.com/it/u=3520060508,2659517975&fm=58', '郭敬明', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (22, 2, '故事的开头总是这样，适逢其会，猝不及防。故事的结局总是这样，花开两朵，天各一方。', 'https://t11.baidu.com/it/u=4192581556,724875387&fm=58', '张嘉佳', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (23, 2, '下定决心，不怕牺牲，排除万难，去争取胜利。', 'https://t12.baidu.com/it/u=1394716907,535750628&fm=58', '毛泽东', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (24, 2, '青春是明知道错了，偏要任性到底！', 'https://t12.baidu.com/it/u=3520060508,2659517975&fm=58', '何炅', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (25, 2, '坚持自己做的事情就可以了，时间会告诉你你的选择正确与否。', 'https://t10.baidu.com/it/u=4208160840,613911946&fm=58', '金星', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (26, 2, '但愿每次回忆，对 生活都不认为负疚。', 'https://t12.baidu.com/it/u=3224252345,3273330929&fm=58', '', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (27, 2, '宁肯少些，但要好些。', 'https://t12.baidu.com/it/u=3208410425,757503539&fm=58', '列宁', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (28, 2, '也许路途很遥远，也许这条路很危险，但是我眼中的风景，是你想像不到的耀眼。', 'https://t12.baidu.com/it/u=472244129,624373840&fm=58', '杨幂', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (29, 2, '别低头，王冠会掉，别流泪，坏人会笑。', 'https://t12.baidu.com/it/u=3687046766,3009085819&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (30, 2, '如果你是一滴水，你是否滋润了一寸土地？如果你是一线阳光，你是否照亮了一分黑暗？如果你是一粒粮食，你是否哺育了有用的生命？如果你是最小的一颗螺丝钉，你是否永远坚守你生活的岗位。', 'https://t11.baidu.com/it/u=3993765373,4203836679&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (31, 2, '我们记得，马吕斯便是从这儿开始的，狂热的恋情忽然出现，并把他推到了种种无目的和无基础的幻想中，他出门仅仅为了去胡思乱想，缓慢的渍染，喧闹而淤止的深渊，并且，随着工作的减少，需要增加了，这是一条规律，处于梦想状态中的人自然是不节约、不振作的，弛懈的精神经受不住紧张的生活，在这种生活方式中，有坏处也有好处，因为慵懒固然有害，慷慨却是健康和善良的，但是不工作的人，穷而慷慨高尚，那是不可救药的，财源涸竭，费用急增， 这是一条导向绝境的下坡路，在这方面，最诚实和最稳定的人也能跟最软弱和最邪恶的人一样往下滑，一直滑到两个深坑中的一个里去：自杀或是犯罪。', 'https://t12.baidu.com/it/u=1579004585,3659949784&fm=58', '雨果', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (32, 2, '阅读的最大理由是想摆脱平庸，早一天就多一份人生的精彩；迟一天就多一天平庸的困扰。', 'https://t10.baidu.com/it/u=174649884,718480879&fm=58', '余秋雨', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `type`=VALUES(`type`),
  `content`=VALUES(`content`),
  `img`=VALUES(`img`),
  `literature_author`=VALUES(`literature_author`),
  `convert_img`=VALUES(`convert_img`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 9.3 SDK模块初始化（菜单、权限、接口）
SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk' AND `type` = 1 AND `deleted_at` = 0 LIMIT 1),
  0
);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    'SDK管理',
    '/sdk',
    '',
    'ele-Key',
    1,
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_root_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk' AND `type` = 1 AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_root_id, 'API管理', '/sdk/api-key', 'sdk/ApiKeyList', 'ele-List', 2, 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_root_id, '调用记录', '/sdk/call-log', 'sdk/SdkCallLogList', 'ele-DocumentChecked', 2, 2, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_root_id, 'SDK接口管理', '/sdk/interface', 'sdk/SdkInterfaceList', 'ele-Setting', 2, 3, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @sdk_api_key_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/api-key' AND `deleted_at` = 0 LIMIT 1);
SET @sdk_call_log_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/call-log' AND `deleted_at` = 0 LIMIT 1);
SET @sdk_interface_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/sdk/interface' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@sdk_api_key_menu_id, 'API Key 新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_api_key_menu_id, 'API Key 分配接口', '', '', '', 3, 4, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_call_log_menu_id, '调用记录导出', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_interface_menu_id, 'SDK接口新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_interface_menu_id, 'SDK接口编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@sdk_interface_menu_id, 'SDK接口删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @btn_api_key_create = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 新增' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_update = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 编辑' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_delete = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 删除' AND `deleted_at`=0 LIMIT 1);
SET @btn_api_key_bind   = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_api_key_menu_id AND `name`='API Key 分配接口' AND `deleted_at`=0 LIMIT 1);
SET @btn_call_log_export = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_call_log_menu_id AND `name`='调用记录导出' AND `deleted_at`=0 LIMIT 1);
SET @btn_interface_create = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口新增' AND `deleted_at`=0 LIMIT 1);
SET @btn_interface_update = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口编辑' AND `deleted_at`=0 LIMIT 1);
SET @btn_interface_delete = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@sdk_interface_menu_id AND `name`='SDK接口删除' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK API Key 列表', 'sdk:key:list', '查看 SDK API Key 列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 新增', 'sdk:key:create', '新增 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 编辑', 'sdk:key:update', '编辑 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 删除', 'sdk:key:delete', '删除 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 分配接口', 'sdk:key:bind_api', '分配接口给 SDK API Key', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录列表', 'sdk:call_log:list', '查看 SDK 调用记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录导出', 'sdk:call_log:export', '导出 SDK 调用记录', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口列表', 'sdk:interface:list', '查看 SDK 接口列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口新增', 'sdk:interface:create', '新增 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口编辑', 'sdk:interface:update', '编辑 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口删除', 'sdk:interface:delete', '删除 SDK 接口', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @perm_key_list      = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_create    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:create' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_update    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:update' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_delete    = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:delete' AND `deleted_at`=0 LIMIT 1);
SET @perm_key_bind      = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:key:bind_api' AND `deleted_at`=0 LIMIT 1);
SET @perm_call_list     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:call_log:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_call_export   = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:call_log:export' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_list       = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:list' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_create     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:create' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_update     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:update' AND `deleted_at`=0 LIMIT 1);
SET @perm_if_delete     = (SELECT `id` FROM `admin_permission` WHERE `code`='sdk:interface:delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('SDK API Key 列表', 'GET', '/api/v1/sdk/key/list', 'SDK API Key 列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 新增', 'POST', '/api/v1/sdk/key/create', '新增 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 编辑', 'POST', '/api/v1/sdk/key/update', '编辑 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 删除', 'POST', '/api/v1/sdk/key/delete', '删除 SDK API Key', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 接口列表', 'GET', '/api/v1/sdk/key/apis', '查看指定 key 的接口授权', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK API Key 接口保存', 'POST', '/api/v1/sdk/key/apis/save', '为 key 保存接口授权与限频', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录列表', 'GET', '/api/v1/sdk/call/log/list', 'SDK 调用记录列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 调用记录导出', 'GET', '/api/v1/sdk/call/log/export', '导出 SDK 调用记录', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口列表', 'GET', '/api/v1/sdk/interface/list', 'SDK 接口列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口新增', 'POST', '/api/v1/sdk/interface/create', '新增 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口编辑', 'POST', '/api/v1/sdk/interface/update', '编辑 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('SDK 接口删除', 'POST', '/api/v1/sdk/interface/delete', '删除 SDK 接口', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @api_key_list     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/list' AND `deleted_at`=0 LIMIT 1);
SET @api_key_create   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/create' AND `deleted_at`=0 LIMIT 1);
SET @api_key_update   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/update' AND `deleted_at`=0 LIMIT 1);
SET @api_key_delete   = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/delete' AND `deleted_at`=0 LIMIT 1);
SET @api_key_apis     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/apis' AND `deleted_at`=0 LIMIT 1);
SET @api_key_save     = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/key/apis/save' AND `deleted_at`=0 LIMIT 1);
SET @api_call_list    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/call/log/list' AND `deleted_at`=0 LIMIT 1);
SET @api_call_export  = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/call/log/export' AND `deleted_at`=0 LIMIT 1);
SET @api_if_list      = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/list' AND `deleted_at`=0 LIMIT 1);
SET @api_if_create    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/create' AND `deleted_at`=0 LIMIT 1);
SET @api_if_update    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/update' AND `deleted_at`=0 LIMIT 1);
SET @api_if_delete    = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/sdk/interface/delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@perm_key_list, @sdk_api_key_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_create, @btn_api_key_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_update, @btn_api_key_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_delete, @btn_api_key_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @btn_api_key_bind, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_list, @sdk_call_log_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_export, @btn_call_log_export, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_list, @sdk_interface_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_create, @btn_interface_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_update, @btn_interface_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_delete, @btn_interface_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@perm_key_list, @api_key_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_create, @api_key_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_update, @api_key_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_delete, @api_key_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @api_key_apis, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_key_bind, @api_key_save, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_list, @api_call_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_call_export, @api_call_export, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_list, @api_if_list, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_create, @api_if_create, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_update, @api_if_update, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@perm_if_delete, @api_if_delete, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 9.4 视频模块初始化（父目录、菜单、权限、接口）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    0,
    '影视资源',
    '/video',
    '',
    'ele-VideoPlay',
    1,
    20,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `path`=VALUES(`path`),
  `icon`=VALUES(`icon`),
  `order_num`=VALUES(`order_num`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/video' AND `deleted_at` = 0 LIMIT 1),
  (SELECT `id` FROM `admin_menu` WHERE `id` = 9 AND `deleted_at` = 0 LIMIT 1)
);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '视频列表管理',
    '/video/list',
    'video/VideoList',
    'ele-Document',
    2,
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/video/list' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@main_menu_id, '视频列表管理 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '视频列表管理 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '视频列表管理 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 新增按钮' AND `deleted_at`=0 LIMIT 1);
SET @update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 编辑按钮' AND `deleted_at`=0 LIMIT 1);
SET @delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='视频列表管理 删除按钮' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频列表管理列表', 'video:list', '查看视频列表管理列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理新增', 'video:create', '新增视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理编辑', 'video:update', '编辑视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理删除', 'video:delete', '删除视频列表管理', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:list' AND `deleted_at`=0 LIMIT 1);
SET @create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:create' AND `deleted_at`=0 LIMIT 1);
SET @update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:update' AND `deleted_at`=0 LIMIT 1);
SET @delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='video:delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('视频列表管理列表', 'GET', '/api/v1/videos', '获取视频列表管理列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理新增', 'POST', '/api/v1/videos', '新增视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理编辑', 'PUT', '/api/v1/videos', '编辑视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频列表管理删除', 'DELETE', '/api/v1/videos', '删除视频列表管理', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频播放器', 'GET', '/api/v1/videos/proxy', '代理m3u8等视频流请求', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('视频采集', 'POST', '/api/v1/videos/collect', '采集视频接口（仅记录操作日志）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('m3u8代理', 'GET', '/api/v1/m3u8/proxy', 'm3u8代理服务（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('m3u8代理OPTIONS', 'OPTIONS', '/api/v1/m3u8/proxy', 'm3u8代理CORS预检请求（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公开视频列表', 'GET', '/api/v1/public/videos/list', '公开视频列表（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('公开视频详情', 'GET', '/api/v1/public/videos/info', '公开视频详情（无权限控制）', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='GET' AND `deleted_at`=0 LIMIT 1);
SET @create_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='POST' AND `deleted_at`=0 LIMIT 1);
SET @update_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='PUT' AND `deleted_at`=0 LIMIT 1);
SET @delete_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/videos' AND `method`='DELETE' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '视频播放器',
    '/video/player',
    'video/VideoPlayer',
    'ele-VideoPlay',
    2,
    1,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

