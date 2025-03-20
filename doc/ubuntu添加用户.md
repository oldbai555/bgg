- 创建新用户 oldbai
```shell
sudo adduser -m oldbai  # 创建用户, 输入账号密码. -m 自动创建用户主目录
sudo usermod -aG sudo oldbai # 将用户添加到 sudo 组
sudo -i -u oldbai  # 输入当前用户的密码（需有 sudo 权限）
sudo whoami  # 验证用户权限 应输出 "root"
```

- 创建新用户 workuser, 拥有 /home/work/package 目录的所有权限且无需使用 sudo 即可远程读写
```shell
sudo useradd -m workuser  # -m 自动创建用户主目录
sudo passwd workuser      # 设置密码
sudo mkdir -p /home/work/package # 确保目录存在
sudo chown -R workuser:workuser /home/work/package # 将目录的所有者和所属组设置为新用户 workuser 及其同名用户组
sudo chmod -R 770 /home/work/package # 赋予所有者和所属组完全权限（读写执行），其他用户无权限
# 770 表示：
# 所有者（workuser）：7（读+写+执行）；
# 所属组（workuser 组）：7；
# 其他用户：0（无权限）。

sudo usermod -aG workuser oldbai # 将用户加入组
groups oldbai # 验证组是否生效
ls -ld /home/work/package # 确保目录组权限正确
# 生效权限（需重新登录或切换会话）
```
