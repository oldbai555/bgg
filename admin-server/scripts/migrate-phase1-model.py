#!/usr/bin/env python3
"""Phase 1: migrate internal/model/*.go to internal/model/<domain>/"""
import os
import re
import shutil

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
MODEL_DIR = os.path.join(ROOT, "internal", "model")

DOMAIN_MAP = {
    "adminusermodel": "iam",
    "adminrolemodel": "iam",
    "adminpermissionmodel": "iam",
    "adminmenumodel": "iam",
    "admindepartmentmodel": "iam",
    "adminuserrolemodel": "iam",
    "adminrolepermissionmodel": "iam",
    "adminapimodel": "iam",
    "adminpermissionmenumodel": "iam",
    "adminpermissionapimodel": "iam",
    "blogtagmodel": "blog",
    "blogarticlemodel": "blog",
    "blogarticletagmodel": "blog",
    "blogarticleauditmodel": "blog",
    "blogfriendlinkmodel": "blog",
    "blogsocialinfomodel": "blog",
    "videomodel": "video",
    "chatmodel": "chat",
    "chatusermodel": "chat",
    "chatmessagemodel": "chat",
    "sdkkeymodel": "sdk",
    "sdkinterfacemodel": "sdk",
    "sdkkeyapimodel": "sdk",
    "sdkcalllogmodel": "sdk",
    "admintaskmodel": "task",
    "adminoperationlogmodel": "monitoring",
    "adminloginlogmodel": "monitoring",
    "auditlogmodel": "monitoring",
    "adminperformancelogmodel": "monitoring",
    "adminconfigmodel": "system",
    "admindicttypemodel": "system",
    "admindictitemmodel": "system",
    "adminfilemodel": "system",
    "adminnoticemodel": "system",
    "adminnotificationmodel": "system",
    "filemodel": "system",
    "demomodel": "misc",
    "dailyshortsentencemodel": "misc",
}

# Type prefix -> domain (for import rewriting)
TYPE_DOMAIN = {
    "AdminUser": "iam",
    "AdminRole": "iam",
    "AdminPermission": "iam",
    "AdminMenu": "iam",
    "AdminDepartment": "iam",
    "AdminUserRole": "iam",
    "AdminRolePermission": "iam",
    "AdminApi": "iam",
    "AdminPermissionMenu": "iam",
    "AdminPermissionApi": "iam",
    "BlogTag": "blog",
    "BlogArticle": "blog",
    "BlogArticleTag": "blog",
    "BlogArticleAudit": "blog",
    "BlogFriendLink": "blog",
    "BlogSocialInfo": "blog",
    "Video": "video",
    "Chat": "chat",
    "ChatUser": "chat",
    "ChatMessage": "chat",
    "SdkKey": "sdk",
    "SdkInterface": "sdk",
    "SdkKeyApi": "sdk",
    "SdkCallLog": "sdk",
    "AdminTask": "task",
    "AdminOperationLog": "monitoring",
    "AdminLoginLog": "monitoring",
    "AuditLog": "monitoring",
    "AdminPerformanceLog": "monitoring",
    "AdminConfig": "system",
    "AdminDictType": "system",
    "AdminDictItem": "system",
    "AdminFile": "system",
    "AdminNotice": "system",
    "AdminNotification": "system",
    "File": "system",
    "Demo": "misc",
    "DailyShortSentence": "misc",
}

def migrate_files():
    for domain in set(DOMAIN_MAP.values()):
        os.makedirs(os.path.join(MODEL_DIR, domain), exist_ok=True)

    for name in os.listdir(MODEL_DIR):
        if not name.endswith(".go"):
            continue
        if name == "vars.go":
            continue
        path = os.path.join(MODEL_DIR, name)
        if not os.path.isfile(path):
            continue
        base = name[:-3]
        lookup = base.replace("_gen", "")
        domain = DOMAIN_MAP.get(lookup)
        if not domain:
            raise SystemExit(f"no domain for {name}")
        dest_dir = os.path.join(MODEL_DIR, domain)
        dest = os.path.join(dest_dir, name)
        with open(path, "r", encoding="utf-8") as f:
            content = f.read()
        content = re.sub(r"^package model$", f"package {domain}", content, count=1, flags=re.MULTILINE)
        with open(dest, "w", encoding="utf-8") as f:
            f.write(content)
        os.remove(path)
        print(f"MOVED {name} -> {domain}/")

def rewrite_imports():
    """Rewrite model.XXX references across admin-server."""
    skip_dirs = {".scratch-goctl", ".git", "vendor"}
    module_root = ROOT

    for dirpath, dirnames, filenames in os.walk(module_root):
        dirnames[:] = [d for d in dirnames if d not in skip_dirs]
        for fn in filenames:
            if not fn.endswith(".go"):
                continue
            fpath = os.path.join(dirpath, fn)
            with open(fpath, "r", encoding="utf-8") as f:
                content = f.read()
            if "internal/model" not in content and "model." not in content:
                continue
            if "/internal/model/" in fpath.replace("\\", "/") and "/internal/model/" in fpath.replace("\\", "/"):
                # skip files inside model subdirs except vars
                rel = os.path.relpath(fpath, MODEL_DIR)
                if rel != "vars.go" and not rel.startswith(".."):
                    continue

            original = content
            domains_used = set()

            # Find model.Type references
            for m in re.finditer(r"\bmodel\.(New)?([A-Z][A-Za-z0-9]+)", content):
                prefix = m.group(2)
                domain = None
                for plen in range(len(prefix), 0, -1):
                    p = prefix[:plen]
                    if p in TYPE_DOMAIN:
                        domain = TYPE_DOMAIN[p]
                        break
                if domain:
                    domains_used.add(domain)

            if not domains_used and 'postapocgame/admin-server/internal/model"' in content:
                # might only use model.ErrNotFound
                if "model.ErrNotFound" in content:
                    pass  # keep model import for vars
                elif '"postapocgame/admin-server/internal/model"' in content:
                    # remove empty model import later
                    pass

            if not domains_used:
                continue

            # Replace model.XXX with domain.XXX
            def repl(m):
                full = m.group(0)
                new_part = m.group(1) or ""
                type_name = m.group(2)
                domain = None
                for plen in range(len(type_name), 0, -1):
                    p = type_name[:plen]
                    if p in TYPE_DOMAIN:
                        domain = TYPE_DOMAIN[p]
                        break
                if domain:
                    return f"{domain}.{new_part}{type_name}"
                return full

            content = re.sub(r"\bmodel\.(New)?([A-Z][A-Za-z0-9]+)", repl, content)

            # Fix imports
            import_line = '"postapocgame/admin-server/internal/model"'
            if import_line in content:
                if "model.ErrNotFound" in content:
                    # keep model import for vars
                    new_imports = []
                    for d in sorted(domains_used):
                        new_imports.append(f'\t"postapocgame/admin-server/internal/model/{d}"')
                    # replace single model import with model + domains
                    lines = content.split("\n")
                    out = []
                    in_import = False
                    replaced = False
                    for line in lines:
                        if line.strip() == "import (":
                            in_import = True
                            out.append(line)
                            continue
                        if in_import and line.strip() == ")":
                            if not replaced:
                                out.append('\t"postapocgame/admin-server/internal/model"')
                                for d in sorted(domains_used):
                                    out.append(f'\t"postapocgame/admin-server/internal/model/{d}"')
                                replaced = True
                            in_import = False
                            out.append(line)
                            continue
                        if in_import and import_line in line:
                            continue  # skip old model import, add later
                        if in_import and any(f'/model/{d}"' in line for d in domains_used):
                            continue
                        out.append(line)
                    content = "\n".join(out)
                else:
                    # remove model import, add domain imports
                    lines = content.split("\n")
                    out = []
                    in_import = False
                    added = False
                    for line in lines:
                        if line.strip() == "import (":
                            in_import = True
                            out.append(line)
                            continue
                        if in_import and line.strip() == ")":
                            if not added:
                                for d in sorted(domains_used):
                                    out.append(f'\t"postapocgame/admin-server/internal/model/{d}"')
                                added = True
                            in_import = False
                            out.append(line)
                            continue
                        if in_import and import_line in line:
                            continue
                        if in_import and any(f'/model/{d}"' in line for d in domains_used):
                            continue
                        out.append(line)
                    content = "\n".join(out)

            # Use package alias matching domain name - imports are without alias, package name IS domain
            if content != original:
                with open(fpath, "w", encoding="utf-8") as f:
                    f.write(content)
                print(f"UPDATED {os.path.relpath(fpath, module_root)}")

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1 and sys.argv[1] == "imports-only":
        rewrite_imports()
    else:
        migrate_files()
        rewrite_imports()
        print("Done.")
