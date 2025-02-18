- 下载
```shell
wget https://dl.min.io/server/minio/release/linux-amd64/minio
```
- 添加权限
```shell
chmod +x minio
```
- 数据目录
```shell
mkdir ~/minio-data
```
- 加入systemctl. 文件名: minio.service
```shell
[Unit]
Description=MinIo
After=network.target

[Service]
Type=simple
ExecStart=/home/bgg/minio server ~/minio-data
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```
- 启动脚本
```shell
# 用户放权
sudo systemctl enable minio.service
sudo systemctl start minio.service
sudo systemctl status minio.service
journalctl -u minio.service # 查看systemctl日志
```
- 配制https,使用自签证书
```shell 
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ~/minio.key -out ~/minio.crt
minio server --secure ~/minio-data
```
- 默认账号密码
```shell
minioadmin
minioadmin
```
