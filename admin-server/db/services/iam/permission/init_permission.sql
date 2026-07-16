-- iam/permission 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

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
