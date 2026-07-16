# Phase 0: goctl 嵌套 group Spike 结论

**日期**: 2026-07-07  
**goctl 版本**: 1.10.1（项目 routes.go 标注 1.9.2，嵌套 group 行为一致）

## 验证方法

使用 [`admin-server/api/scratch.api`](../admin-server/api/scratch.api)，包含三个嵌套 group：
- `misc/ping`
- `iam/auth`
- `iam/api`（用于检测跨模块 import 别名冲突）

```bash
goctl api go -api api/scratch.api -dir .scratch-goctl --style go_zero
```

## 结论：**采用嵌套 group 方案**

| 检查项 | 结果 |
|--------|------|
| 生成 `handler/misc/ping/` | ✅ |
| 生成 `logic/misc/ping/` | ✅ |
| 生成 `handler/iam/auth/`、`handler/iam/api/` | ✅ |
| routes.go import 别名 | `miscping`、`iamauth`、`iamapi` — 域名+模块拼接，**无冲突** |
| 同域不同模块（iam/auth vs iam/api） | 别名 `iamauth` vs `iamapi`，可区分 |

## 最终方案

**全部 44 个 group 改为 `group: <domain>/<module>`**，利用 goctl 原生嵌套目录生成。

### group 映射（旧 → 新）

| 域 | 旧 group | 新 group |
|----|----------|----------|
| misc | ping | misc/ping |
| monitoring | monitor | monitoring/monitor |
| iam | auth | iam/auth |
| iam | user | iam/user |
| iam | role | iam/role |
| iam | permission | iam/permission |
| iam | department | iam/department |
| iam | menu | iam/menu |
| iam | api | iam/api |
| iam | user_role | iam/user_role |
| iam | role_permission | iam/role_permission |
| iam | permission_menu | iam/permission_menu |
| iam | permission_api | iam/permission_api |
| system | file | system/file |
| system | config | system/config |
| system | dict_type | system/dict_type |
| system | dict_item | system/dict_item |
| system | dict | system/dict |
| system | notice | system/notice |
| system | notification | system/notification |
| misc | demo | misc/demo |
| misc | daily_short_sentence | misc/daily_short_sentence |
| video | video | video/video |
| video | m3u8 | video/m3u8 |
| video | video_collect | video/video_collect |
| video | public_video | video/public |
| blog | blog_tag | blog/tag |
| blog | blog_article | blog/article |
| blog | blog_article_audit | blog/article_audit |
| blog | public_blog | blog/public |
| blog | blog_friend_link | blog/friend_link |
| blog | blog_social_info | blog/social_info |
| monitoring | metric | monitoring/metric |
| monitoring | metric_admin | monitoring/metric_admin |
| monitoring | operation_log | monitoring/operation_log |
| monitoring | login_log | monitoring/login_log |
| monitoring | audit_log | monitoring/audit_log |
| monitoring | performance_log | monitoring/performance_log |
| chat | chat | chat/chat |
| chat | chat_group | chat/group |
| chat | chat_message | chat/message |
| sdk | sdk | sdk/sdk |
| sdk | sdk_public | sdk/public |
| task | task | task/task |
| task | task_public | task/public |
| misc | public | misc/public |

**注意**：`blog_article` → `blog/article`（去重域名前缀）；`chat_group` → `chat/group`。

## 未采用的 fallback

- ~~扁平前缀 `iam_auth`~~：可读但无嵌套目录收益
- ~~生成后手工挪目录~~：44 模块手工成本高，spike 已验证原生嵌套可行
