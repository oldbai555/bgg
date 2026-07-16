-- iam/user 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

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
