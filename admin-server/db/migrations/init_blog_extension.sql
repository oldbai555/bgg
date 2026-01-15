-- ============================================
-- 博客扩展功能初始化 SQL
-- 功能组: blog
-- 功能名称: 博客扩展功能（友情链接、社交信息、文章置顶）
-- 创建时间：2026-01-16
-- ============================================

-- ============================================
-- 1. 获取博客管理主菜单 ID（path = '/blog'）
-- ============================================
SET @blog_root_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/blog' AND `deleted_at` = 0
  LIMIT 1
);

-- ============================================
-- 2. 获取文章管理菜单 ID（path = '/blog/article'）
-- ============================================
SET @blog_article_menu_id = (
  SELECT `id` FROM `admin_menu`
  WHERE `path` = '/blog/article' AND `deleted_at` = 0
  LIMIT 1
);

-- ============================================
-- 3. 插入菜单数据
-- ============================================

-- 3.1 友情链接管理菜单
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
  @blog_root_menu_id,
  '友情链接管理',
  '/blog/friend-link',
  'blog/BlogFriendLinkList',
  'ele-Link',
  2, -- 菜单
  4,
  1,
  1,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
);

SET @blog_friend_link_menu_id = LAST_INSERT_ID();

-- 3.2 社交信息管理菜单
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
  @blog_root_menu_id,
  '社交信息管理',
  '/blog/social-info',
  'blog/BlogSocialInfoList',
  'ele-Share',
  2, -- 菜单
  5,
  1,
  1,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP(),
  0
);

SET @blog_social_info_menu_id = LAST_INSERT_ID();

-- 3.3 友情链接管理按钮菜单（新增/编辑/删除）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_friend_link_menu_id, '友情链接新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_friend_link_menu_id, '友情链接编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_friend_link_menu_id, '友情链接删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_friend_link_create_menu_id = LAST_INSERT_ID() - 2;
SET @blog_friend_link_update_menu_id = LAST_INSERT_ID() - 1;
SET @blog_friend_link_delete_menu_id = LAST_INSERT_ID();

-- 3.4 社交信息管理按钮菜单（新增/编辑/删除）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_social_info_menu_id, '社交信息新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_social_info_menu_id, '社交信息编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_social_info_menu_id, '社交信息删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_social_info_create_menu_id = LAST_INSERT_ID() - 2;
SET @blog_social_info_update_menu_id = LAST_INSERT_ID() - 1;
SET @blog_social_info_delete_menu_id = LAST_INSERT_ID();

-- 3.5 文章置顶按钮菜单（置顶/取消置顶）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_article_menu_id, '文章置顶', '', '', '', 3, 7, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章取消置顶', '', '', '', 3, 8, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_article_top_menu_id = LAST_INSERT_ID() - 1;
SET @blog_article_untop_menu_id = LAST_INSERT_ID();

-- ============================================
-- 4. 插入权限数据
-- ============================================

-- 友情链接管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('友情链接列表','blog_friend_link:list','查看友情链接列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接新增','blog_friend_link:create','新增友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接编辑','blog_friend_link:update','编辑友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接删除','blog_friend_link:delete','删除友情链接',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_friend_link_list_permission_id   = LAST_INSERT_ID() - 3;
SET @blog_friend_link_create_permission_id = LAST_INSERT_ID() - 2;
SET @blog_friend_link_update_permission_id = LAST_INSERT_ID() - 1;
SET @blog_friend_link_delete_permission_id = LAST_INSERT_ID();

-- 社交信息管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('社交信息列表','blog_social_info:list','查看社交信息列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息新增','blog_social_info:create','新增社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息编辑','blog_social_info:update','编辑社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息删除','blog_social_info:delete','删除社交信息',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_social_info_list_permission_id   = LAST_INSERT_ID() - 3;
SET @blog_social_info_create_permission_id = LAST_INSERT_ID() - 2;
SET @blog_social_info_update_permission_id = LAST_INSERT_ID() - 1;
SET @blog_social_info_delete_permission_id = LAST_INSERT_ID();

-- 文章置顶权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('文章置顶','blog_article:top','设置文章置顶',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('文章取消置顶','blog_article:untop','取消文章置顶',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_top_permission_id   = LAST_INSERT_ID() - 1;
SET @blog_article_untop_permission_id = LAST_INSERT_ID();

-- ============================================
-- 5. 插入接口数据
-- ============================================

-- 友情链接管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('友情链接列表接口','GET','/api/v1/blog/friend-links','友情链接列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接新增接口','POST','/api/v1/blog/friend-links','新增友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接编辑接口','PUT','/api/v1/blog/friend-links','编辑友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('友情链接删除接口','DELETE','/api/v1/blog/friend-links','删除友情链接',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_friend_link_list_api_id   = LAST_INSERT_ID() - 3;
SET @blog_friend_link_create_api_id = LAST_INSERT_ID() - 2;
SET @blog_friend_link_update_api_id = LAST_INSERT_ID() - 1;
SET @blog_friend_link_delete_api_id = LAST_INSERT_ID();

-- 社交信息管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('社交信息列表接口','GET','/api/v1/blog/social-infos','社交信息列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息新增接口','POST','/api/v1/blog/social-infos','新增社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息编辑接口','PUT','/api/v1/blog/social-infos','编辑社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('社交信息删除接口','DELETE','/api/v1/blog/social-infos','删除社交信息',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_social_info_list_api_id   = LAST_INSERT_ID() - 3;
SET @blog_social_info_create_api_id = LAST_INSERT_ID() - 2;
SET @blog_social_info_update_api_id = LAST_INSERT_ID() - 1;
SET @blog_social_info_delete_api_id = LAST_INSERT_ID();

-- 文章置顶接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('文章置顶接口','POST','/api/v1/blog/articles/top','设置文章置顶',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('文章取消置顶接口','POST','/api/v1/blog/articles/untop','取消文章置顶',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_top_api_id   = LAST_INSERT_ID() - 1;
SET @blog_article_untop_api_id = LAST_INSERT_ID();

-- 公共接口（友情链接列表、社交信息列表、标签列表）
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('公共友情链接列表接口','GET','/api/v1/public/blog/friend-links','获取启用的友情链接列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('公共社交信息列表接口','GET','/api/v1/public/blog/social-infos','获取启用的社交信息列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('公共博客标签列表接口','GET','/api/v1/public/blog/tags','获取启用的博客标签列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @public_blog_friend_link_list_api_id = LAST_INSERT_ID() - 2;
SET @public_blog_social_info_list_api_id = LAST_INSERT_ID() - 1;
SET @public_blog_tag_list_api_id          = LAST_INSERT_ID();

-- ============================================
-- 6. 插入权限-菜单关联数据
-- ============================================

-- 友情链接：列表权限 -> 友情链接菜单；C/U/D -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_friend_link_list_permission_id,   @blog_friend_link_menu_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_create_permission_id, @blog_friend_link_create_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_update_permission_id, @blog_friend_link_update_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_delete_permission_id, @blog_friend_link_delete_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 社交信息：列表权限 -> 社交信息菜单；C/U/D -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_social_info_list_permission_id,   @blog_social_info_menu_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_create_permission_id, @blog_social_info_create_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_update_permission_id, @blog_social_info_update_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_delete_permission_id, @blog_social_info_delete_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 文章置顶：置顶权限 -> 文章置顶按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_top_permission_id,   @blog_article_top_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_untop_permission_id, @blog_article_untop_menu_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- ============================================
-- 7. 插入权限-接口关联数据
-- ============================================

-- 友情链接接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_friend_link_list_permission_id,   @blog_friend_link_list_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_create_permission_id, @blog_friend_link_create_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_update_permission_id, @blog_friend_link_update_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_friend_link_delete_permission_id, @blog_friend_link_delete_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 社交信息接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_social_info_list_permission_id,   @blog_social_info_list_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_create_permission_id, @blog_social_info_create_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_update_permission_id, @blog_social_info_update_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_social_info_delete_permission_id, @blog_social_info_delete_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 文章置顶接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_top_permission_id,   @blog_article_top_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_untop_permission_id, @blog_article_untop_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());
