#!/usr/bin/env python3
"""Phase 2: migrate internal/repository/*_repository.go to internal/repository/<domain>/"""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
REPO_DIR = os.path.join(ROOT, "internal", "repository")

DOMAIN_MAP = {
    "user_repository": "iam",
    "role_repository": "iam",
    "permission_repository": "iam",
    "menu_repository": "iam",
    "department_repository": "iam",
    "user_role_repository": "iam",
    "role_permission_repository": "iam",
    "api_repository": "iam",
    "permission_menu_repository": "iam",
    "permission_api_repository": "iam",
    "token_blacklist_repository": "iam",
    "blog_tag_repository": "blog",
    "blog_article_repository": "blog",
    "blog_article_tag_repository": "blog",
    "blog_article_audit_repository": "blog",
    "blog_friend_link_repository": "blog",
    "blog_social_info_repository": "blog",
    "video_repository": "video",
    "chat_repository": "chat",
    "chat_user_repository": "chat",
    "chat_message_repository": "chat",
    "sdk_repository": "sdk",
    "sdk_admin_repository": "sdk",
    "task_repository": "task",
    "operation_log_repository": "monitoring",
    "login_log_repository": "monitoring",
    "audit_log_repository": "monitoring",
    "performance_log_repository": "monitoring",
    "metric_repository": "monitoring",
    "config_repository": "system",
    "dict_type_repository": "system",
    "dict_item_repository": "system",
    "file_repository": "system",
    "notice_repository": "system",
    "notification_repository": "system",
    "demo_repository": "misc",
    "daily_short_sentence_repository": "misc",
}

# Constructor -> domain
CTOR_DOMAIN = {
    "NewUserRepository": "iam",
    "NewRoleRepository": "iam",
    "NewPermissionRepository": "iam",
    "NewMenuRepository": "iam",
    "NewDepartmentRepository": "iam",
    "NewUserRoleRepository": "iam",
    "NewRolePermissionRepository": "iam",
    "NewApiRepository": "iam",
    "NewPermissionMenuRepository": "iam",
    "NewPermissionApiRepository": "iam",
    "NewTokenBlacklistRepository": "iam",
    "NewBlogTagRepository": "blog",
    "NewBlogArticleRepository": "blog",
    "NewBlogArticleTagRepository": "blog",
    "NewBlogArticleAuditRepository": "blog",
    "NewBlogFriendLinkRepository": "blog",
    "NewBlogSocialInfoRepository": "blog",
    "NewVideoRepository": "video",
    "NewChatRepository": "chat",
    "NewChatUserRepository": "chat",
    "NewChatMessageRepository": "chat",
    "NewSdkRepository": "sdk",
    "NewSdkAdminRepository": "sdk",
    "NewTaskRepository": "task",
    "NewOperationLogRepository": "monitoring",
    "NewLoginLogRepository": "monitoring",
    "NewAuditLogRepository": "monitoring",
    "NewPerformanceLogRepository": "monitoring",
    "NewMetricRepository": "monitoring",
    "NewConfigRepository": "system",
    "NewDictTypeRepository": "system",
    "NewDictItemRepository": "system",
    "NewFileRepository": "system",
    "NewNoticeRepository": "system",
    "NewNotificationRepository": "system",
    "NewDemoRepository": "misc",
    "NewDailyShortSentenceRepository": "misc",
}

# Interface type -> domain (for ServiceContext fields etc.)
TYPE_DOMAIN = {
    "UserRepository": "iam",
    "RoleRepository": "iam",
    "PermissionRepository": "iam",
    "MenuRepository": "iam",
    "DepartmentRepository": "iam",
    "UserRoleRepository": "iam",
    "RolePermissionRepository": "iam",
    "ApiRepository": "iam",
    "PermissionMenuRepository": "iam",
    "PermissionApiRepository": "iam",
    "TokenBlacklistRepository": "iam",
    "BlogTagRepository": "blog",
    "BlogArticleRepository": "blog",
    "BlogArticleTagRepository": "blog",
    "BlogArticleAuditRepository": "blog",
    "BlogFriendLinkRepository": "blog",
    "BlogSocialInfoRepository": "blog",
    "VideoRepository": "video",
    "ChatRepository": "chat",
    "ChatUserRepository": "chat",
    "ChatMessageRepository": "chat",
    "SdkRepository": "sdk",
    "SdkAdminRepository": "sdk",
    "TaskRepository": "task",
    "OperationLogRepository": "monitoring",
    "LoginLogRepository": "monitoring",
    "AuditLogRepository": "monitoring",
    "PerformanceLogRepository": "monitoring",
    "MetricRepository": "monitoring",
    "ConfigRepository": "system",
    "DictTypeRepository": "system",
    "DictItemRepository": "system",
    "FileRepository": "system",
    "NoticeRepository": "system",
    "NotificationRepository": "system",
    "DemoRepository": "misc",
    "DailyShortSentenceRepository": "misc",
}

KEEP_TOP = {"repository.go", "sql_conn.go", "cache_conf.go"}


def migrate_files():
    for domain in set(DOMAIN_MAP.values()):
        os.makedirs(os.path.join(REPO_DIR, domain), exist_ok=True)

    for name in os.listdir(REPO_DIR):
        if not name.endswith("_repository.go"):
            continue
        key = name[:-3]  # strip .go
        domain = DOMAIN_MAP.get(key)
        if not domain:
            raise SystemExit(f"no domain for {name}")
        path = os.path.join(REPO_DIR, name)
        dest = os.path.join(REPO_DIR, domain, name)
        with open(path, "r", encoding="utf-8") as f:
            content = f.read()
        content = re.sub(r"^package repository$", f"package {domain}", content, count=1, flags=re.MULTILINE)
        with open(dest, "w", encoding="utf-8") as f:
            f.write(content)
        os.remove(path)
        print(f"MOVED {name} -> {domain}/")


def rewrite_imports():
    skip_dirs = {".scratch-goctl", ".git", "vendor"}
    for dirpath, dirnames, filenames in os.walk(ROOT):
        dirnames[:] = [d for d in dirnames if d not in skip_dirs]
        for fn in filenames:
            if not fn.endswith(".go"):
                continue
            fpath = os.path.join(dirpath, fn)
            rel = os.path.relpath(fpath, REPO_DIR)
            if not rel.startswith("..") and rel != "repository.go" and rel != "sql_conn.go" and rel != "cache_conf.go":
                if "/repository/" in fpath.replace("\\", "/") or rel.endswith("_repository.go"):
                    if rel.count(os.sep) == 0 and fn in KEEP_TOP:
                        pass
                    elif "/repository/" in fpath.replace("\\", "/"):
                        continue  # skip domain repo files for ctor rewrite

            with open(fpath, "r", encoding="utf-8") as f:
                content = f.read()

            original = content
            domains_used = set()

            for ctor, domain in CTOR_DOMAIN.items():
                if f"repository.{ctor}" in content:
                    domains_used.add(domain)
                    content = content.replace(f"repository.{ctor}", f"{domain}.{ctor}")

            for typ, domain in TYPE_DOMAIN.items():
                if re.search(rf"\brepository\.{typ}\b", content):
                    domains_used.add(domain)
                    content = re.sub(rf"\brepository\.{typ}\b", f"{domain}.{typ}", content)

            if not domains_used:
                continue

            import_line = '"postapocgame/admin-server/internal/repository"'
            if import_line not in content:
                # add domain imports only
                lines = content.split("\n")
                out, in_import, added = [], False, False
                for line in lines:
                    if line.strip() == "import (":
                        in_import = True
                        out.append(line)
                        continue
                    if in_import and line.strip() == ")":
                        if not added:
                            for d in sorted(domains_used):
                                if f'/repository/{d}"' not in content:
                                    out.append(f'\t"postapocgame/admin-server/internal/repository/{d}"')
                            added = True
                        in_import = False
                        out.append(line)
                        continue
                    if in_import and any(f'/repository/{d}"' in line for d in domains_used):
                        continue
                    out.append(line)
                content = "\n".join(out)
            else:
                lines = content.split("\n")
                out, in_import, replaced = [], False, False
                for line in lines:
                    if line.strip() == "import (":
                        in_import = True
                        out.append(line)
                        continue
                    if in_import and line.strip() == ")":
                        if not replaced:
                            out.append(import_line)
                            for d in sorted(domains_used):
                                out.append(f'\t"postapocgame/admin-server/internal/repository/{d}"')
                            replaced = True
                        in_import = False
                        out.append(line)
                        continue
                    if in_import and import_line in line:
                        continue
                    if in_import and any(f'/repository/{d}"' in line for d in domains_used):
                        continue
                    out.append(line)
                content = "\n".join(out)

            if content != original:
                with open(fpath, "w", encoding="utf-8") as f:
                    f.write(content)
                print(f"UPDATED {os.path.relpath(fpath, ROOT)}")


def update_repository_go():
    path = os.path.join(REPO_DIR, "repository.go")
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()

    domain_imports = sorted(set(DOMAIN_MAP.values()))
    import_block = "\n".join(
        f'\t"postapocgame/admin-server/internal/repository/{d}"' for d in domain_imports
    )
    model_imports = {
        "iam": "postapocgame/admin-server/internal/model/iam",
        "blog": "postapocgame/admin-server/internal/model/blog",
        "video": "postapocgame/admin-server/internal/model/video",
        "chat": "postapocgame/admin-server/internal/model/chat",
        "sdk": "postapocgame/admin-server/internal/model/sdk",
        "task": "postapocgame/admin-server/internal/model/task",
        "monitoring": "postapocgame/admin-server/internal/model/monitoring",
        "system": "postapocgame/admin-server/internal/model/system",
        "misc": "postapocgame/admin-server/internal/model/misc",
    }
    model_import_block = "\n".join(f'\t"{p}"' for p in model_imports.values())

    # Replace import section - read current file and rebuild imports manually
    lines = content.split("\n")
    # find import block end
    new_lines = []
    i = 0
    while i < len(lines):
        if lines[i].strip() == "import (":
            new_lines.append(lines[i])
            new_lines.append('\t"postapocgame/admin-server/internal/config"')
            for p in model_imports.values():
                new_lines.append(f'\t"{p}"')
            for d in domain_imports:
                new_lines.append(f'\t"postapocgame/admin-server/internal/repository/{d}"')
            new_lines.append('\tbusinesscache "postapocgame/admin-server/pkg/cache"')
            new_lines.append("")
            new_lines.append('\t"github.com/pkg/errors"')
            new_lines.append('\t"github.com/zeromicro/go-zero/core/stores/cache"')
            new_lines.append('\t"github.com/zeromicro/go-zero/core/stores/redis"')
            new_lines.append('\t"github.com/zeromicro/go-zero/core/stores/sqlx"')
            new_lines.append(")")
            i += 1
            while i < len(lines) and lines[i].strip() != ")":
                i += 1
            i += 1
            continue
        new_lines.append(lines[i])
        i += 1

    content = "\n".join(new_lines)

    # model refs already use iam., blog., etc. from phase1 - verify
    if content != open(path).read():
        with open(path, "w", encoding="utf-8") as f:
            f.write(content)
        print("UPDATED repository.go imports")


if __name__ == "__main__":
    migrate_files()
    rewrite_imports()
    print("Phase 2 file migration done. Run go build to verify.")
