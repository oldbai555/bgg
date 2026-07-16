# 使用tmpfs（内存文件系统）缓存静态文件
```shell
# 1. 创建并挂载tmpfs
sudo mkdir -p /tmp/admin-cache
sudo mount -t tmpfs -o size=200M tmpfs /tmp/admin-cache

# 2. 复制文件
sudo cp -r /home/work/service/admin-frontend/dist/* /tmp/admin-cache/

# 3. 检查文件是否复制成功
ls -la /tmp/admin-cache/

## 将这一行
#alias /home/work/service/admin-frontend/dist;
## 改为
#alias /tmp/admin-cache;
```

# 使用shell+vim在Ubuntu下配置systemd服务
- 创建启动脚本
```shell
sudo vim /usr/local/bin/mount-admin-cache.sh
```
- 写入以下内容
```shell
#!/bin/bash

# 创建目录
mkdir -p /tmp/admin-cache

# 挂载tmpfs
if ! mountpoint -q /tmp/admin-cache; then
    mount -t tmpfs -o size=200M tmpfs /tmp/admin-cache
fi

# 复制文件
cp -r /home/work/service/admin-frontend/dist/* /tmp/admin-cache/

# 设置权限
chown -R www-data:www-data /tmp/admin-cache
chmod -R 755 /tmp/admin-cache

echo "Admin cache mounted and files copied successfully"
```
- 添加执行权限
```shell
sudo chmod +x /usr/local/bin/mount-admin-cache.sh
```
- 测试脚本
```shell
sudo /usr/local/bin/mount-admin-cache.sh
```
- 创建systemd服务
```shell
sudo /usr/local/bin/mount-admin-cache.sh
```
- 写入以下内容
```shell
[Unit]
Description=Mount and populate admin frontend cache
After=network.target
Before=nginx.service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/mount-admin-cache.sh
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
```
- 启用并启动服务
```shell
# 重载systemd配置
sudo systemctl daemon-reload

# 启用服务（开机自动运行）
sudo systemctl enable admin-cache.service

# 启动服务
sudo systemctl start admin-cache.service

# 查看状态
sudo systemctl status admin-cache.service
```
- 创建更新脚本（用于以后更新前端）
```shell
sudo vim /usr/local/bin/update-admin-cache.sh
```
- 写入以下内容
```shell
#!/bin/bash

echo "Updating admin frontend cache..."

# 删除旧文件
rm -rf /tmp/admin-cache/*

# 复制新文件
cp -r /home/work/service/admin-frontend/dist/* /tmp/admin-cache/

# 设置权限
chown -R www-data:www-data /tmp/admin-cache
chmod -R 755 /tmp/admin-cache

echo "Admin cache updated successfully"
```
- 添加执行权限
```shell
sudo chmod +x /usr/local/bin/update-admin-cache.sh
```
- 验证
```shell
# 检查挂载
df -h | grep admin-cache

# 检查文件
ls -la /tmp/admin-cache/

# 检查服务状态
sudo systemctl status admin-cache.service
```
