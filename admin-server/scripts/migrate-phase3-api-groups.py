#!/usr/bin/env python3
"""Update admin.api group names to domain/module format."""
import re
from pathlib import Path

API = Path(__file__).resolve().parent.parent / "api" / "admin.api"

REPLACEMENTS = [
    (r"group:\s+ping\b", "group:      misc/ping"),
    (r"group:\s+monitor\b", "group:      monitoring/monitor"),
    (r"group:\s+auth\b", "group:      iam/auth"),
    (r"group:\s+user\b", "group:      iam/user"),
    (r"group:\s+role\b", "group:      iam/role"),
    (r"group:\s+permission\b", "group:      iam/permission"),
    (r"group:\s+department\b", "group:      iam/department"),
    (r"group:\s+menu\b", "group:      iam/menu"),
    (r"group:\s+api\b", "group:      iam/api"),
    (r"group:\s+user_role\b", "group:      iam/user_role"),
    (r"group:\s+role_permission\b", "group:      iam/role_permission"),
    (r"group:\s+permission_menu\b", "group:      iam/permission_menu"),
    (r"group:\s+permission_api\b", "group:      iam/permission_api"),
    (r"group:\s+file\b", "group:      system/file"),
    (r"group:\s+config\b", "group:      system/config"),
    (r"group:\s+dict_type\b", "group:      system/dict_type"),
    (r"group:\s+dict_item\b", "group:      system/dict_item"),
    (r"group:\s+dict\b", "group:      system/dict"),
    (r"group:\s+demo\b", "group:      misc/demo"),
    (r"group:\s+daily_short_sentence\b", "group:      misc/daily_short_sentence"),
    (r"group:\s+video\b", "group:      video/video"),
    (r"group:\s+m3u8\b", "group:      video/m3u8"),
    (r"group:\s+video_collect\b", "group:      video/video_collect"),
    (r"group:\s+public_video\b", "group:      video/public"),
    (r"group:\s+blog_tag\b", "group:      blog/tag"),
    (r"group:\s+blog_article\b", "group:      blog/article"),
    (r"group:\s+blog_article_audit\b", "group:      blog/article_audit"),
    (r"group:\s+public_blog\b", "group:      blog/public"),
    (r"group:\s+blog_friend_link\b", "group:      blog/friend_link"),
    (r"group:\s+blog_social_info\b", "group:      blog/social_info"),
    (r"group:\s+metric\b", "group:      monitoring/metric"),
    (r"group:\s+metric_admin\b", "group:      monitoring/metric_admin"),
    (r"group:\s+chat\b", "group:      chat/chat"),
    (r"group:\s+chat_group\b", "group:      chat/group"),
    (r"group:\s+chat_message\b", "group:      chat/message"),
    (r"group:\s+operation_log\b", "group:      monitoring/operation_log"),
    (r"group:\s+login_log\b", "group:      monitoring/login_log"),
    (r"group:\s+audit_log\b", "group:      monitoring/audit_log"),
    (r"group:\s+performance_log\b", "group:      monitoring/performance_log"),
    (r"group:\s+notice\b", "group:      system/notice"),
    (r"group:\s+notification\b", "group:      system/notification"),
    (r"group:\s+sdk\b", "group:      sdk/sdk"),
    (r"group:\s+task\b", "group:      task/task"),
    (r"group:\s+task_public\b", "group:      task/public"),
    (r"group:\s+sdk_public\b", "group:      sdk/public"),
    (r"group:\s+public\b", "group:      misc/public"),
]

def main():
    text = API.read_text(encoding="utf-8")
    for pattern, repl in REPLACEMENTS:
        text = re.sub(pattern, repl, text)
    API.write_text(text, encoding="utf-8")
    print("Updated", API)

if __name__ == "__main__":
    main()
