#!/bin/bash
# 删除.vimrc
cd ~
rm -rf .vimrc
# 恢复网络配置文件
netplan_path="/etc/netplan/"
netplan_filename="50-cloud-init.yaml"
cd ${netplan_path}
sudo cp -f "${netplan_filename}.backup" ${netplan_filename}
# 删除配置文件
cfg_path="/etc/cloud/cloud.cfg.d/"
cfg_file="99-disable-network-config.cfg"
sudo rm -rf "${cfg_path}${cfg_file}"
# 恢复对profile的修改
profile_path="/etc/profile"
# 恢复文件
sudo cp -f "${profile_path}.backup" "${profile_path}"
source "${profile_path}"
# 刷新网络
sudo netplan apply
# 恢复对sudo免密的设置，只需要删除文件即可
sudo rm -rf "/etc/sudoers.d/${USER}"
# 注意执行完脚本后将网络连接模式修改为NAT
