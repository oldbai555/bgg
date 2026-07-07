#!/usr/bin/env python3
import os, re
ROOT = os.path.join(os.path.dirname(__file__), "..", "internal", "logic")
REPL = [
    ("l.svcCtx.BlogTagRepository", "blogrepo.NewBlogTagRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.BlogArticleRepository", "blogrepo.NewBlogArticleRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.BlogArticleTagRepository", "blogrepo.NewBlogArticleTagRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.BlogArticleAuditRepository", "blogrepo.NewBlogArticleAuditRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.BlogFriendLinkRepository", "blogrepo.NewBlogFriendLinkRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.BlogSocialInfoRepository", "blogrepo.NewBlogSocialInfoRepository(l.svcCtx.Repository)"),
    ("l.svcCtx.UserRepository", "iamrepo.NewUserRepository(l.svcCtx.Repository)"),
]
for dirpath, _, files in os.walk(ROOT):
    for fn in files:
        if not fn.endswith('.go'): continue
        path = os.path.join(dirpath, fn)
        with open(path) as f: c = f.read()
        orig = c
        for old, new in REPL:
            c = c.replace(old, new)
        if c == orig: continue
        if 'blogrepo "' not in c and 'blogrepo.New' in c:
            c = c.replace('import (\n', 'import (\n\tblogrepo "postapocgame/admin-server/internal/repository/blog"\n', 1)
        if 'iamrepo "' not in c and 'iamrepo.New' in c:
            c = c.replace('import (\n', 'import (\n\tiamrepo "postapocgame/admin-server/internal/repository/iam"\n', 1)
        with open(path, 'w') as f: f.write(c)
        print(path)
