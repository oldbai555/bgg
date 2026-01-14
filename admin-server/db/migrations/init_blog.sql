-- 博客模块初始化 SQL
-- 功能组: blog
-- 功能名称: 博客管理（标签、文章、审核）

-- ============================================
-- 2. 插入菜单数据
-- ============================================

-- 2.1 博客管理主菜单
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
    0,
  '博客管理',
  '/blog',
  'blog/BlogLayout',
    'ele-Document',
  1, -- 目录
  0,
  1,
  1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @blog_root_menu_id = LAST_INSERT_ID();

-- 2.2 标签管理菜单
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
  @blog_root_menu_id,
  '标签管理',
  '/blog/tag',
  'blog/BlogTagList',
  'ele-PriceTag',
  2, -- 菜单
  1,
  1,
  1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @blog_tag_menu_id = LAST_INSERT_ID();

-- 2.3 文章管理菜单
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
  @blog_root_menu_id,
  '文章管理',
  '/blog/article',
  'blog/BlogArticleList',
  'ele-DocumentCopy',
  2, -- 菜单
  2,
  1,
  1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @blog_article_menu_id = LAST_INSERT_ID();

-- 2.4 文章审核菜单（可选，如不单独展示可仅用于权限控制）
INSERT INTO `admin_menu` (
  `parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`
)
VALUES (
  @blog_root_menu_id,
  '文章审核',
  '/blog/article-audit',
  'blog/BlogArticleAuditList',
  'ele-Finished',
  2, -- 菜单
  3,
  1,
  1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
);

SET @blog_article_audit_menu_id = LAST_INSERT_ID();

-- 2.5 标签管理按钮菜单（新增/编辑/删除）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_tag_menu_id, '标签新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_tag_menu_id, '标签编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_tag_menu_id, '标签删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_tag_create_menu_id = LAST_INSERT_ID() - 2;
SET @blog_tag_update_menu_id = LAST_INSERT_ID() - 1;
SET @blog_tag_delete_menu_id = LAST_INSERT_ID();

-- 2.6 文章管理按钮菜单（新增/编辑/删除/提交审核/上架/下架）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_article_menu_id, '文章新增', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章编辑', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章删除', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章提交审核', '', '', '', 3, 4, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章上架', '', '', '', 3, 5, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_menu_id, '文章下架', '', '', '', 3, 6, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_article_create_menu_id = LAST_INSERT_ID() - 5;
SET @blog_article_update_menu_id = LAST_INSERT_ID() - 4;
SET @blog_article_delete_menu_id = LAST_INSERT_ID() - 3;
SET @blog_article_submit_menu_id = LAST_INSERT_ID() - 2;
SET @blog_article_publish_menu_id = LAST_INSERT_ID() - 1;
SET @blog_article_unpublish_menu_id = LAST_INSERT_ID();

-- 2.7 文章审核按钮菜单（审核通过/驳回/下架）
INSERT INTO `admin_menu` (`parent_id`,`name`,`path`,`component`,`icon`,`type`,`order_num`,`visible`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  (@blog_article_audit_menu_id, '文章审核通过/驳回', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@blog_article_audit_menu_id, '文章审核下架', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0);

SET @blog_article_audit_menu_btn_id = LAST_INSERT_ID() - 1;
SET @blog_article_audit_unpublish_menu_id = LAST_INSERT_ID();

-- ============================================
-- 3. 插入权限数据
-- ============================================

-- 标签管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客标签列表','blog_tag:list','查看博客标签列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签新增','blog_tag:create','新增博客标签',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签编辑','blog_tag:update','编辑博客标签',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签删除','blog_tag:delete','删除博客标签',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_tag_list_permission_id   = LAST_INSERT_ID() - 3;
SET @blog_tag_create_permission_id = LAST_INSERT_ID() - 2;
SET @blog_tag_update_permission_id = LAST_INSERT_ID() - 1;
SET @blog_tag_delete_permission_id = LAST_INSERT_ID();

-- 文章管理权限
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客文章列表','blog_article:list','查看博客文章列表',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章新增','blog_article:create','新增博客文章',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章编辑','blog_article:update','编辑博客文章',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章删除','blog_article:delete','删除博客文章',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章提交审核','blog_article:submit','提交文章审核',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章上架','blog_article:publish','上架文章',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章下架','blog_article:unpublish','下架文章',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_list_permission_id      = LAST_INSERT_ID() - 6;
SET @blog_article_create_permission_id    = LAST_INSERT_ID() - 5;
SET @blog_article_update_permission_id    = LAST_INSERT_ID() - 4;
SET @blog_article_delete_permission_id    = LAST_INSERT_ID() - 3;
SET @blog_article_submit_permission_id    = LAST_INSERT_ID() - 2;
SET @blog_article_publish_permission_id   = LAST_INSERT_ID() - 1;
SET @blog_article_unpublish_permission_id = LAST_INSERT_ID();

-- 文章审核权限（审核员）
INSERT INTO `admin_permission` (`name`,`code`,`description`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客文章审核','blog_article:audit','审核博客文章（通过/驳回）',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章审核下架','blog_article:audit_unpublish','审核员执行文章下架',UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_audit_permission_id         = LAST_INSERT_ID() - 1;
SET @blog_article_audit_unpublish_permission_id = LAST_INSERT_ID();

-- ============================================
-- 4. 插入接口数据（后台标签/文章管理 + 文章审核）
-- 注意：路径不使用 :id，具体参数以 Query/Body 为准，对应 admin.api 定义
-- ============================================

-- 标签管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客标签列表接口','GET','/api/v1/blog/tags','博客标签列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签新增接口','POST','/api/v1/blog/tags','新增博客标签',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签编辑接口','PUT','/api/v1/blog/tags','编辑博客标签',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客标签删除接口','DELETE','/api/v1/blog/tags','删除博客标签',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_tag_list_api_id   = LAST_INSERT_ID() - 3;
SET @blog_tag_create_api_id = LAST_INSERT_ID() - 2;
SET @blog_tag_update_api_id = LAST_INSERT_ID() - 1;
SET @blog_tag_delete_api_id = LAST_INSERT_ID();

-- 文章管理接口
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客文章列表接口','GET','/api/v1/blog/articles','博客文章列表',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章详情接口','GET','/api/v1/blog/articles/detail','博客文章详情',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章新增接口','POST','/api/v1/blog/articles','新增博客文章',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章编辑接口','PUT','/api/v1/blog/articles','编辑博客文章',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章删除接口','DELETE','/api/v1/blog/articles','删除博客文章',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章提交审核接口','POST','/api/v1/blog/articles/submit','提交文章审核',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章上架接口','POST','/api/v1/blog/articles/publish','上架文章',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章下架接口','POST','/api/v1/blog/articles/unpublish','下架文章',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_list_api_id      = LAST_INSERT_ID() - 7;
SET @blog_article_detail_api_id    = LAST_INSERT_ID() - 6;
SET @blog_article_create_api_id    = LAST_INSERT_ID() - 5;
SET @blog_article_update_api_id    = LAST_INSERT_ID() - 4;
SET @blog_article_delete_api_id    = LAST_INSERT_ID() - 3;
SET @blog_article_submit_api_id    = LAST_INSERT_ID() - 2;
SET @blog_article_publish_api_id   = LAST_INSERT_ID() - 1;
SET @blog_article_unpublish_api_id = LAST_INSERT_ID();

-- 文章审核接口（审核员）
INSERT INTO `admin_api` (`name`,`method`,`path`,`description`,`status`,`created_at`,`updated_at`,`deleted_at`)
VALUES
  ('博客文章审核接口','POST','/api/v1/blog/articles/audit','审核博客文章（通过/驳回）',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0),
  ('博客文章审核下架接口','POST','/api/v1/blog/articles/audit/unpublish','审核员执行文章下架',1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),0);

SET @blog_article_audit_api_id          = LAST_INSERT_ID() - 1;
SET @blog_article_audit_unpublish_api_id = LAST_INSERT_ID();

-- ============================================
-- 5. 插入权限-菜单关联数据
-- ============================================

-- 标签：列表权限 -> 标签菜单；C/U/D -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_tag_list_permission_id,   @blog_tag_menu_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_create_permission_id, @blog_tag_create_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_update_permission_id, @blog_tag_update_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_delete_permission_id, @blog_tag_delete_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 文章：列表权限 -> 文章菜单；其他权限 -> 按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_list_permission_id,      @blog_article_menu_id,           UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_create_permission_id,    @blog_article_create_menu_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_update_permission_id,    @blog_article_update_menu_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_delete_permission_id,    @blog_article_delete_menu_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_submit_permission_id,    @blog_article_submit_menu_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_publish_permission_id,   @blog_article_publish_menu_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_unpublish_permission_id, @blog_article_unpublish_menu_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 审核：审核权限 -> 文章审核菜单及按钮
INSERT INTO `admin_permission_menu` (`permission_id`,`menu_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_audit_permission_id,            @blog_article_audit_menu_id,           UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_audit_permission_id,            @blog_article_audit_menu_btn_id,       UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_audit_unpublish_permission_id,  @blog_article_audit_unpublish_menu_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- ============================================
-- 6. 插入权限-接口关联数据
-- ============================================

-- 标签接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_tag_list_permission_id,   @blog_tag_list_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_create_permission_id, @blog_tag_create_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_update_permission_id, @blog_tag_update_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_tag_delete_permission_id, @blog_tag_delete_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 文章接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_list_permission_id,      @blog_article_list_api_id,      UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_list_permission_id,      @blog_article_detail_api_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_create_permission_id,    @blog_article_create_api_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_update_permission_id,    @blog_article_update_api_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_delete_permission_id,    @blog_article_delete_api_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_submit_permission_id,    @blog_article_submit_api_id,    UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_publish_permission_id,   @blog_article_publish_api_id,   UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_unpublish_permission_id, @blog_article_unpublish_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

-- 审核接口权限绑定
INSERT INTO `admin_permission_api` (`permission_id`,`api_id`,`created_at`,`updated_at`)
VALUES
  (@blog_article_audit_permission_id,           @blog_article_audit_api_id,          UNIX_TIMESTAMP(),UNIX_TIMESTAMP()),
  (@blog_article_audit_unpublish_permission_id, @blog_article_audit_unpublish_api_id, UNIX_TIMESTAMP(),UNIX_TIMESTAMP());

