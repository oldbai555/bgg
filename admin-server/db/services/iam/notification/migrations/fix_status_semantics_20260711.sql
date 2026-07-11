-- 修复 admin_notification.read_status / admin_login_log.status 两套冲突的取值约定。
--
-- 背景：这两个字段的建表 DDL 注释和实际写入代码历史上一直用的是一套"布尔式"取值
-- （read_status: 0=未读,1=已读；status: 0=失败,1=成功），但对应的字典（read_status/
-- login_status，见 db/data.sql）按项目"字典枚举 value 从 1 开始"的约定写成了另一套
-- （read_status: 1=未读,2=已读；login_status: 1=成功,2=失败）。列表筛选走的是字典的
-- 取值，写入走的是旧的取值，导致标记已读/全部已读/清除已读/登录日志按状态筛选实际上
-- 全部失效。本次统一到字典的取值（代码同步修复，见 admin-server/docs/progress.md）。
--
-- 幂等性说明：本文件不是幂等的（UPDATE 会重复执行两次导致数据再次错位），只应该在一个
-- 库上执行一次。全新部署（先跑 tables.sql + data.sql）不需要执行本文件，data.sql 的
-- 种子数据已经是修复后的取值；只有已经跑过旧版 data.sql、且积累了真实业务数据的库需要。

-- admin_notification.read_status：旧值 1(已读) -> 新值 2；旧值 0(未读) -> 新值 1。
-- 必须先处理 1->2，避免和后续 0->1 的目标值冲突。
UPDATE `admin_notification` SET `read_status` = 2 WHERE `read_status` = 1;
UPDATE `admin_notification` SET `read_status` = 1 WHERE `read_status` = 0;

-- admin_login_log.status：旧值 0(失败) -> 新值 2；旧值 1(成功) 在新旧两套取值里都是 1，
-- 不需要迁移。
UPDATE `admin_login_log` SET `status` = 2 WHERE `status` = 0;
