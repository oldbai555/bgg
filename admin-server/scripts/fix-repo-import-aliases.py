#!/usr/bin/env python3
"""Fix repository domain import aliases to avoid conflicts with model packages."""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOMAINS = ["iam", "blog", "video", "chat", "sdk", "task", "monitoring", "system", "misc"]

CTOR_NAMES = [
    "NewUserRepository", "NewRoleRepository", "NewPermissionRepository", "NewMenuRepository",
    "NewDepartmentRepository", "NewUserRoleRepository", "NewRolePermissionRepository",
    "NewApiRepository", "NewPermissionMenuRepository", "NewPermissionApiRepository",
    "NewTokenBlacklistRepository", "NewBlogTagRepository", "NewBlogArticleRepository",
    "NewBlogArticleTagRepository", "NewBlogArticleAuditRepository", "NewBlogFriendLinkRepository",
    "NewBlogSocialInfoRepository", "NewVideoRepository", "NewChatRepository", "NewChatUserRepository",
    "NewChatMessageRepository", "NewSdkRepository", "NewSdkAdminRepository", "NewTaskRepository",
    "NewOperationLogRepository", "NewLoginLogRepository", "NewAuditLogRepository",
    "NewPerformanceLogRepository", "NewMetricRepository", "NewConfigRepository",
    "NewDictTypeRepository", "NewDictItemRepository", "NewFileRepository", "NewNoticeRepository",
    "NewNotificationRepository", "NewDemoRepository", "NewDailyShortSentenceRepository",
]

REPO_TYPES = [
    "UserRepository", "RoleRepository", "PermissionRepository", "MenuRepository",
    "DepartmentRepository", "UserRoleRepository", "RolePermissionRepository", "ApiRepository",
    "PermissionMenuRepository", "PermissionApiRepository", "TokenBlacklistRepository",
    "BlogTagRepository", "BlogArticleRepository", "BlogArticleTagRepository",
    "BlogArticleAuditRepository", "BlogFriendLinkRepository", "BlogSocialInfoRepository",
    "VideoRepository", "ChatRepository", "ChatUserRepository", "ChatMessageRepository",
    "SdkRepository", "SdkAdminRepository", "TaskRepository", "OperationLogRepository",
    "LoginLogRepository", "AuditLogRepository", "PerformanceLogRepository", "MetricRepository",
    "ConfigRepository", "DictTypeRepository", "DictItemRepository", "FileRepository",
    "NoticeRepository", "NotificationRepository", "DemoRepository", "DailyShortSentenceRepository",
    "TaskQueryFilter",
]


def fix_file(path: str) -> bool:
    with open(path, encoding="utf-8") as f:
        content = f.read()
    original = content

    domains_in_repo = []
    for d in DOMAINS:
        imp = f'"postapocgame/admin-server/internal/repository/{d}"'
        if imp in content:
            domains_in_repo.append(d)

    if not domains_in_repo:
        return False

    for d in domains_in_repo:
        alias = f"{d}repo"
        bare = f'"postapocgame/admin-server/internal/repository/{d}"'
        aliased = f'{alias} "postapocgame/admin-server/internal/repository/{d}"'
        if bare in content and aliased not in content:
            content = content.replace(f"\t{bare}", f"\t{aliased}")

    for d in domains_in_repo:
        alias = f"{d}repo"
        for ctor in CTOR_NAMES:
            content = content.replace(f"{d}.{ctor}", f"{alias}.{ctor}")
        for typ in REPO_TYPES:
            content = re.sub(rf"\b{d}\.{typ}\b", f"{alias}.{typ}", content)
        content = content.replace(f"repository.{typ}", f"taskrepo.{typ}")  # wrong, skip

    # repository.TaskQueryFilter -> taskrepo.TaskQueryFilter
    if "repository.TaskQueryFilter" in content:
        if "taskrepo" not in content and "repository/task" in content:
            content = content.replace("repository.TaskQueryFilter", "taskrepo.TaskQueryFilter")

    # Remove standalone repository import if unused
    lines = content.split("\n")
    if '"postapocgame/admin-server/internal/repository"' in content:
        uses_repo_type = bool(re.search(r"\brepository\.(Repository|BuildSources|NewRepository)\b", content))
        if not uses_repo_type:
            new_lines = []
            in_import = False
            for line in lines:
                if in_import and '"postapocgame/admin-server/internal/repository"' in line and "repository/" not in line:
                    continue
                new_lines.append(line)
                if line.strip() == "import (":
                    in_import = True
                elif in_import and line.strip() == ")":
                    in_import = False
            content = "\n".join(new_lines)

    if content != original:
        with open(path, "w", encoding="utf-8") as f:
            f.write(content)
        return True
    return False


def main():
    skip = {".scratch-goctl", ".git", "vendor"}
    count = 0
    for dirpath, dirnames, filenames in os.walk(ROOT):
        dirnames[:] = [d for d in dirnames if d not in skip]
        for fn in filenames:
            if not fn.endswith(".go"):
                continue
            path = os.path.join(dirpath, fn)
            rel = os.path.relpath(path, os.path.join(ROOT, "internal", "repository"))
            if not rel.startswith("..") and rel != "repository.go":
                if "/" not in rel.replace("\\", "/") or rel.endswith("_repository.go"):
                    if rel not in ("repository.go", "sql_conn.go", "cache_conf.go") and "/" not in rel:
                        pass
            # skip domain repo implementation files
            if "/repository/" in path.replace("\\", "/"):
                parts = path.replace("\\", "/").split("/repository/")
                if len(parts) > 1 and "/" in parts[1]:
                    continue
            if fix_file(path):
                count += 1
                print("fixed", os.path.relpath(path, ROOT))
    print(f"done {count} files")


if __name__ == "__main__":
    main()
