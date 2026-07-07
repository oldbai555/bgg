#!/usr/bin/env python3
"""Fix logic package references in handlers after domain migration."""
import os
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "internal" / "handler"

OLD_TO_PKG = {
    "blog_tag": "tag", "blog_article": "article", "blog_article_audit": "article_audit",
    "blog_friend_link": "friend_link", "blog_social_info": "social_info", "public_blog": "public",
    "chat_group": "group", "chat_message": "message", "public_video": "public",
    "video_collect": "video_collect", "daily_short_sentence": "daily_short_sentence",
    "dict_type": "dict_type", "dict_item": "dict_item", "user_role": "user_role",
    "role_permission": "role_permission", "permission_menu": "permission_menu",
    "permission_api": "permission_api", "operation_log": "operation_log",
    "login_log": "login_log", "audit_log": "audit_log", "performance_log": "performance_log",
    "metric_admin": "metric_admin", "task_public": "public", "sdk_public": "public",
    "video": "video", "auth": "auth", "user": "user", "role": "role", "permission": "permission",
    "department": "department", "menu": "menu", "api": "api", "file": "file", "config": "config",
    "dict": "dict", "notice": "notice", "notification": "notification", "demo": "demo",
    "ping": "ping", "public": "public", "m3u8": "m3u8", "metric": "metric", "monitor": "monitor",
    "chat": "chat", "sdk": "sdk", "task": "task",
}


def fix_file(path: Path) -> bool:
    content = path.read_text(encoding="utf-8")
    orig = content
    for old, new in OLD_TO_PKG.items():
        content = re.sub(rf"\b{old}\.New", f"{new}.New", content)
    if content != orig:
        path.write_text(content, encoding="utf-8")
        return True
    return False


def main():
    n = 0
    for dirpath, _, files in os.walk(ROOT):
        for fn in files:
            if fn.endswith("handler.go"):
                if fix_file(Path(dirpath) / fn):
                    n += 1
    print(f"fixed {n} files")


if __name__ == "__main__":
    main()
