-- ============================================
-- 博客扩展功能初始化 SQL
-- 功能组: blog
-- 功能名称: 博客扩展功能（友情链接、社交信息、文章置顶）
-- 创建时间：2026-01-16
--
-- 幂等性说明：admin_menu 没有唯一键约束，菜单类 INSERT 统一改用
-- INSERT ... SELECT ... WHERE NOT EXISTS（按 path 或 parent_id+name 判重），
-- 不使用对它不生效的 ON DUPLICATE KEY UPDATE；权限/接口/关联表本身有唯一键，
-- 继续用 ON DUPLICATE KEY UPDATE。也不再用 LAST_INSERT_ID() 偏移量推算 ID，
-- 详见 docs/changelog/archive-backend.md 第 23 节。
-- ============================================

-- ============================================
-- 1. 获取博客管理主菜单 ID（path = '/admin/blog'）
-- ============================================
SET @blog_root_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/admin/blog' AND `deleted_at` = 0
  LIMIT 1
);

-- ============================================
-- 2. 获取文章管理菜单 ID（path = '/admin/blog/article'）
-- ============================================
SET @blog_article_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/admin/blog/article' AND `deleted_at` = 0
  LIMIT 1
);

-- ============================================
-- 3. 插入菜单数据
-- ============================================

-- 3.1 友情链接管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_root_menu_id, '友情链接管理', '/admin/blog/friend-link', 'content/BlogFriendLinkList', 'ele-Link', 2, 4, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `path` = '/admin/blog/friend-link' AND `deleted_at` = 0);

SET @blog_friend_link_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/blog/friend-link' AND `deleted_at` = 0 LIMIT 1);

-- 3.2 社交信息管理菜单
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_root_menu_id, '社交信息管理', '/admin/blog/social-info', 'content/BlogSocialInfoList', 'ele-Share', 2, 5, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `path` = '/admin/blog/social-info' AND `deleted_at` = 0);

SET @blog_social_info_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/blog/social-info' AND `deleted_at` = 0 LIMIT 1);

-- 3.3 友情链接管理按钮菜单（新增/编辑/删除，无 path，按 parent_id+name 判重）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_friend_link_menu_id, '友情链接新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接新增' AND `deleted_at` = 0);
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_friend_link_menu_id, '友情链接编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接编辑' AND `deleted_at` = 0);
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_friend_link_menu_id, '友情链接删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接删除' AND `deleted_at` = 0);

SET @blog_friend_link_create_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接新增' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_update_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接编辑' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_delete_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_friend_link_menu_id AND `name` = '友情链接删除' AND `deleted_at` = 0 LIMIT 1);

-- 3.4 社交信息管理按钮菜单（新增/编辑/删除）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_social_info_menu_id, '社交信息新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息新增' AND `deleted_at` = 0);
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_social_info_menu_id, '社交信息编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息编辑' AND `deleted_at` = 0);
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_social_info_menu_id, '社交信息删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息删除' AND `deleted_at` = 0);

SET @blog_social_info_create_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息新增' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_update_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息编辑' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_delete_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_social_info_menu_id AND `name` = '社交信息删除' AND `deleted_at` = 0 LIMIT 1);

-- 3.5 文章置顶按钮菜单（置顶/取消置顶）
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_article_menu_id, '文章置顶', '', '', '', 3, 7, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_article_menu_id AND `name` = '文章置顶' AND `deleted_at` = 0);
INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
SELECT @blog_article_menu_id, '文章取消置顶', '', '', '', 3, 8, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admin_menu` WHERE `parent_id` = @blog_article_menu_id AND `name` = '文章取消置顶' AND `deleted_at` = 0);

SET @blog_article_top_menu_id   = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_article_menu_id AND `name` = '文章置顶' AND `deleted_at` = 0 LIMIT 1);
SET @blog_article_untop_menu_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id` = @blog_article_menu_id AND `name` = '文章取消置顶' AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 4. 插入权限数据
-- ============================================

-- 友情链接管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('友情链接列表','blog_friend_link:list','查看友情链接列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接新增','blog_friend_link:create','新增友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接编辑','blog_friend_link:update','编辑友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接删除','blog_friend_link:delete','删除友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_friend_link_list_permission_id   = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_friend_link:list' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_friend_link:create' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_friend_link:update' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_friend_link:delete' AND `deleted_at` = 0 LIMIT 1);

-- 社交信息管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('社交信息列表','blog_social_info:list','查看社交信息列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息新增','blog_social_info:create','新增社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息编辑','blog_social_info:update','编辑社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息删除','blog_social_info:delete','删除社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_social_info_list_permission_id   = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_social_info:list' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_social_info:create' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_social_info:update' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_social_info:delete' AND `deleted_at` = 0 LIMIT 1);

-- 文章置顶权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('文章置顶','blog_article:top','设置文章置顶',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('文章取消置顶','blog_article:untop','取消文章置顶',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_article_top_permission_id   = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_article:top' AND `deleted_at` = 0 LIMIT 1);
SET @blog_article_untop_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code` = 'blog_article:untop' AND `deleted_at` = 0 LIMIT 1);

-- ============================================
-- 5. 插入接口数据
-- ============================================

-- 友情链接管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('友情链接列表接口','GET','/api/v1/blog/friend-links','友情链接列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接新增接口','POST','/api/v1/blog/friend-links','新增友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接编辑接口','PUT','/api/v1/blog/friend-links','编辑友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接删除接口','DELETE','/api/v1/blog/friend-links','删除友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_friend_link_list_api_id   = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/blog/friend-links' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_create_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/blog/friend-links' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_update_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'PUT' AND `path` = '/api/v1/blog/friend-links' AND `deleted_at` = 0 LIMIT 1);
SET @blog_friend_link_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/blog/friend-links' AND `deleted_at` = 0 LIMIT 1);

-- 社交信息管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('社交信息列表接口','GET','/api/v1/blog/social-infos','社交信息列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息新增接口','POST','/api/v1/blog/social-infos','新增社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息编辑接口','PUT','/api/v1/blog/social-infos','编辑社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息删除接口','DELETE','/api/v1/blog/social-infos','删除社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_social_info_list_api_id   = (SELECT `id` FROM `admin_api` WHERE `method` = 'GET' AND `path` = '/api/v1/blog/social-infos' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_create_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/blog/social-infos' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_update_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'PUT' AND `path` = '/api/v1/blog/social-infos' AND `deleted_at` = 0 LIMIT 1);
SET @blog_social_info_delete_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'DELETE' AND `path` = '/api/v1/blog/social-infos' AND `deleted_at` = 0 LIMIT 1);

-- 文章置顶接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('文章置顶接口','POST','/api/v1/blog/articles/top','设置文章置顶',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('文章取消置顶接口','POST','/api/v1/blog/articles/untop','取消文章置顶',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @blog_article_top_api_id   = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/blog/articles/top' AND `deleted_at` = 0 LIMIT 1);
SET @blog_article_untop_api_id = (SELECT `id` FROM `admin_api` WHERE `method` = 'POST' AND `path` = '/api/v1/blog/articles/untop' AND `deleted_at` = 0 LIMIT 1);

-- 公共接口（友情链接列表、社交信息列表、标签列表）
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('公共友情链接列表接口','GET','/api/v1/public/blog/friend-links','获取启用的友情链接列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('公共社交信息列表接口','GET','/api/v1/public/blog/social-infos','获取启用的社交信息列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('公共博客标签列表接口','GET','/api/v1/public/blog/tags','获取启用的博客标签列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

-- ============================================
-- 6. 插入权限-菜单关联数据
-- ============================================

-- 友情链接：列表权限 -> 友情链接菜单；C/U/D -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_friend_link_list_permission_id,   @blog_friend_link_menu_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_create_permission_id, @blog_friend_link_create_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_update_permission_id, @blog_friend_link_update_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_delete_permission_id, @blog_friend_link_delete_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 社交信息：列表权限 -> 社交信息菜单；C/U/D -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_social_info_list_permission_id,   @blog_social_info_menu_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_create_permission_id, @blog_social_info_create_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_update_permission_id, @blog_social_info_update_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_delete_permission_id, @blog_social_info_delete_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 文章置顶：置顶权限 -> 文章置顶按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_top_permission_id,   @blog_article_top_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_untop_permission_id, @blog_article_untop_menu_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 7. 插入权限-接口关联数据
-- ============================================

-- 友情链接接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_friend_link_list_permission_id,   @blog_friend_link_list_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_create_permission_id, @blog_friend_link_create_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_update_permission_id, @blog_friend_link_update_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_delete_permission_id, @blog_friend_link_delete_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 社交信息接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_social_info_list_permission_id,   @blog_social_info_list_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_create_permission_id, @blog_social_info_create_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_update_permission_id, @blog_social_info_update_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_delete_permission_id, @blog_social_info_delete_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 文章置顶接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_top_permission_id,   @blog_article_top_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_untop_permission_id, @blog_article_untop_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();
