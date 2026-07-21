-- AI 知识库问答模块初始化 SQL（权限/接口数据）
-- 功能组: ai/knowledge_qa
-- 详见 admin-server/docs/ai-knowledge-qa-spec.md

-- ============================================
-- 1. 本模块没有独立菜单
-- ============================================
-- Phase 1 只做后端原型，还没有对应的前端页面（不套用 generate-sql.sh 标准 CRUD 脚手架，
-- 见 spec 里"偏离单表 CRUD 手写"的决策），照 db/services/iam/demo/init_demo.sql 的先例，
-- 不建 admin_menu 记录，避免产生指向不存在页面的孤儿菜单；后续做前端页面时再补菜单。

-- ============================================
-- 2. 插入权限数据
-- ============================================
INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('AI知识库问答', 'ai_knowledge_qa:ask', '向 AI 知识库提问', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('AI知识库重建索引', 'ai_knowledge_qa:reindex', '重建 AI 知识库向量索引', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @ask_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='ai_knowledge_qa:ask' AND `deleted_at`=0 LIMIT 1);
SET @reindex_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='ai_knowledge_qa:reindex' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 3. 插入接口数据
-- ============================================
-- path 必须与 admin-server/api/admin.api 里 ai/knowledge_qa 服务块实际路由一致
INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('AI知识库问答', 'POST', '/api/v1/ai/knowledge-qa/ask', '向 AI 知识库提问', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('AI知识库重建索引', 'POST', '/api/v1/ai/knowledge-qa/reindex', '重建 AI 知识库向量索引', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `description`=VALUES(`description`), `status`=VALUES(`status`), `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @ask_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='POST' AND `path`='/api/v1/ai/knowledge-qa/ask' AND `deleted_at`=0 LIMIT 1);
SET @reindex_api_id = (SELECT `id` FROM `admin_api` WHERE `method`='POST' AND `path`='/api/v1/ai/knowledge-qa/reindex' AND `deleted_at`=0 LIMIT 1);

-- ============================================
-- 4. 插入权限-接口关联数据
-- ============================================
INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@ask_permission_id, @ask_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@reindex_permission_id, @reindex_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- ============================================
-- 5. 给超级管理员角色（role_id=1，沿用仓库既有约定）补上这两个权限
-- ============================================
INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
VALUES
  (1, @ask_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (1, @reindex_permission_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();
