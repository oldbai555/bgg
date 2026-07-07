#!/usr/bin/env python3
"""Merge handlers from backup, fixing package and logic import paths."""
import os
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BACKUP = ROOT / ".handler-backup"
HANDLER = ROOT / "internal" / "handler"
DOMAIN_DIRS = {"iam", "blog", "video", "chat", "sdk", "task", "monitoring", "system", "misc"}

LOGIC_MAP = {
    "auth": "iam/auth", "user": "iam/user", "role": "iam/role", "permission": "iam/permission",
    "department": "iam/department", "menu": "iam/menu", "api": "iam/api",
    "user_role": "iam/user_role", "role_permission": "iam/role_permission",
    "permission_menu": "iam/permission_menu", "permission_api": "iam/permission_api",
    "file": "system/file", "config": "system/config", "dict_type": "system/dict_type",
    "dict_item": "system/dict_item", "dict": "system/dict", "notice": "system/notice",
    "notification": "system/notification", "demo": "misc/demo",
    "daily_short_sentence": "misc/daily_short_sentence", "ping": "misc/ping", "public": "misc/public",
    "video": "video/video", "m3u8": "video/m3u8", "video_collect": "video/video_collect",
    "public_video": "video/public", "blog_tag": "blog/tag", "blog_article": "blog/article",
    "blog_article_audit": "blog/article_audit", "public_blog": "blog/public",
    "blog_friend_link": "blog/friend_link", "blog_social_info": "blog/social_info",
    "metric": "monitoring/metric", "metric_admin": "monitoring/metric_admin",
    "monitor": "monitoring/monitor", "operation_log": "monitoring/operation_log",
    "login_log": "monitoring/login_log", "audit_log": "monitoring/audit_log",
    "performance_log": "monitoring/performance_log", "chat": "chat/chat",
    "chat_group": "chat/group", "chat_message": "chat/message",
    "sdk": "sdk/sdk", "sdk_public": "sdk/public", "task": "task/task", "task_public": "task/public",
}


def handler_func_name(content: str):
    m = re.search(r"func (\w+Handler)\(", content)
    return m.group(1) if m else None


def package_name(content: str):
    m = re.search(r"^package (\w+)$", content, re.M)
    return m.group(1) if m else None


def fix_logic_imports(content: str) -> str:
    def repl(m):
        mod = m.group(1)
        new = LOGIC_MAP.get(mod, mod)
        return f"postapocgame/admin-server/internal/logic/{new}"
    return re.sub(r"postapocgame/admin-server/internal/logic/([a-zA-Z0-9_]+)", repl, content)


def build_index():
    idx = {}
    for dirpath, _, files in os.walk(HANDLER):
        rel = Path(dirpath).relative_to(HANDLER)
        if rel.parts and rel.parts[0] not in DOMAIN_DIRS:
            continue
        for fn in files:
            if not fn.endswith("handler.go"):
                continue
            path = Path(dirpath) / fn
            content = path.read_text(encoding="utf-8")
            name = handler_func_name(content)
            if name:
                idx[name] = path
    return idx


def merge():
    index = build_index()
    merged = 0
    for dirpath, _, files in os.walk(BACKUP):
        for fn in files:
            if fn in ("routes.go", "custom_routes.go"):
                continue
            if not fn.endswith("handler.go"):
                continue
            content = (Path(dirpath) / fn).read_text(encoding="utf-8")
            hname = handler_func_name(content)
            if not hname or hname not in index:
                continue
            target = index[hname]
            new_pkg = package_name(target.read_text(encoding="utf-8"))
            old_pkg = package_name(content)
            if new_pkg != old_pkg:
                content = re.sub(rf"^package {old_pkg}$", f"package {new_pkg}", content, count=1, flags=re.M)
            content = fix_logic_imports(content)
            content = re.sub(r"// goctl [\d.]+", "// goctl 1.10.1", content, count=1)
            target.write_text(content, encoding="utf-8")
            merged += 1
    print(f"merged {merged} handlers")


if __name__ == "__main__":
    merge()
