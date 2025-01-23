- 导入ovf镜像时，如果IP不可用
```shell
# 固定静态IP
cd /etc/netplan
# 查看文件
ll
# 我这叫下面这个名字
vim 00-installer-config.yaml
# 查看网关
route -n
# 填写yaml内容 主要是下面有注释的那三行
network:
  ethernets:
    ens33:
      dhcp4: false                          # 禁用dhcp    
      gateway4: 192.168.226.2               # 设置网关地址
      addresses: [192.168.226.5/24]         # 设置静态IP地址和掩码
      nameservers:
              addresses: [8.8.8.8, 114.114.114.114] # dns
  version: 2
# 保存生效  
sudo netplan apply
```
