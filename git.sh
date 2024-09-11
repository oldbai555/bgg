git reset "$(git merge-base main "$(git branch --show-current)")"
git add -A && git commit -m 'v.0.0.2'
#git pull
#git push --force
# git remote set-url origin https://<你的令牌>@github.com/<你的git用户名>/<要修改的仓库名>.git
# git reset --hard HEAD 回到最后一次提交之前
