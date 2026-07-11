-- iam/rbac 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 关联：用户-角色、角色-权限
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

-- 权限-菜单关联
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

-- 权限-接口关联
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
