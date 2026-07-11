-- iam/api 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

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
  -- 群组管理接口（需要权限控制；路径不使用 :id，与 admin.api 一致）
  (59, '群组列表', 'GET', '/api/v1/chats/groups', '获取群组列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (60, '群组创建', 'POST', '/api/v1/chats/groups', '创建群组', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (61, '群组更新', 'PUT', '/api/v1/chats/groups', '更新群组信息', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (62, '群组删除', 'DELETE', '/api/v1/chats/groups', '删除群组', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (63, '群组详情', 'GET', '/api/v1/chats/groups/detail', '获取群组详情', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (64, '群组成员列表', 'GET', '/api/v1/chats/groups/members', '获取群组成员列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (65, '群组成员添加', 'POST', '/api/v1/chats/groups/members', '添加群组成员', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (66, '群组成员移除', 'DELETE', '/api/v1/chats/groups/members', '移除群组成员', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `deleted_at`=0;
