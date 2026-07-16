-- 路由命名空间彻底分离：后台管理 /bgg/admin/* 与公共页面 /bgg/front/* 不再共享任何路径前缀段
-- 背景：admin-frontend 之前用同一个 createWebHistory('/admin/') 同时承载后台管理页面和公共
-- 展示页面，两者在路由 path 层面共用同一套字面量前缀（如 /blog、/videos）。公共页是启动时就注册
-- 的静态路由，后台管理页是登录后异步拉取菜单才 addRoute() 注册的动态路由——F5 硬刷新时如果地址栏
-- 恰好落在和公共页字面量相同的后台路径（如「博客管理」目录 path=/blog），会在动态路由注册完成前
-- 被公共静态路由抢先匹配，表现为莫名其妙跳到公共页。详见 admin-frontend/docs/10-route-namespace-migration.md。
--
-- 前端路由方案改为 /admin/* + /front/* 两个不相交命名空间后，后台菜单的 path 需要统一加 /admin 前缀
-- 才能和前端新路由定义对上；公共页面（/front/*）不是菜单驱动的，不需要改。
--
-- 幂等性：WHERE path NOT LIKE '/admin%'，重复执行天然跳过已加过前缀的记录
UPDATE admin_menu
SET path = CONCAT('/admin', path), updated_at = UNIX_TIMESTAMP()
WHERE path IS NOT NULL AND path != '' AND path NOT LIKE '/admin%' AND deleted_at = 0;

-- 同一次迁移顺带修正字典项「在线聊天页面路径」（chat_config 下 label='在线聊天页面路径'），
-- 这是 MessageNotification.vue 里"点击聊天消息跳转"用的路径，同样需要加 /admin 前缀。
-- 不是新增字典项，是既有数据值修正，不走 dict_{module}_YYYYMMDD.sql 的新增字典命名模式。
UPDATE admin_dict_item di
JOIN admin_dict_type dt ON dt.id = di.type_id
SET di.value = CONCAT('/admin', di.value), di.updated_at = UNIX_TIMESTAMP()
WHERE dt.code = 'chat_config' AND di.label = '在线聊天页面路径' AND di.value NOT LIKE '/admin%' AND di.deleted_at = 0;
