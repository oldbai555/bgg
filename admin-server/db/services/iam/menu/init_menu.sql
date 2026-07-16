-- iam/menu 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

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
