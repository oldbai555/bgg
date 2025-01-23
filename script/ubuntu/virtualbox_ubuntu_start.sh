#!/bin/bash
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
BLUE='\033[34m'
PURPLE='\033[35m'
CYAN='\033[36m'
WHITE='\033[37m'
RESET='\033[0m'

tip() {
    local len=${#1}
    local n=$((len+22))
    local color_1=${RED}
    local color_2=${GREEN}
    echo -e "${color_1}$(printf "%0.s-" $(seq 1 ${n}))${RESET}"
    echo -e "${color_1}++++++++++${RESET} ${color_2}${1}${RESET} ${color_1}++++++++++${RESET}"
    echo -e "${color_1}$(printf "%0.s-" $(seq 1 ${n}))${RESET}"
}

# 首先, 更新系统源
tip "UPDATE SYSTEM"
sudo apt update
sudo apt upgrade -y
# 设置用户使用 sudo 命令免密
# 根据文件 /etc/sudoers 中的注释, 比起直接修改 /etc/sudoers
# 更好的方法是在 /etc/sudoers.d/ 目录中加入配置文件
tip "MAKE YOU BECOME A SUDOER AND NO NEED PASSWORD"
echo "${USER} ALL=(ALL) NOPASSWD: ALL" | sudo tee "/etc/sudoers.d/${USER}"
# 如果安装的是最小化的系统的话, 上面没有 vi 编辑器, 也不能用 visudo 命令
# 只能先安装编辑器, 选择 vim
tip "INSTALL EDITOR [vim]"
sudo apt install vim
# 配置 vim
cd ~
echo -e 'set nu\nset tabstop=4\nset softtabstop=4\nset shiftwidth=4\nset expandtab\nset showmatch\nset ruler\nset autoindent\nset undofile\nset undodir=~/.vim/undodir\nset showmode\nset showcmd\nset hlsearch\nset t_Co=256\nset noerrorbells\nset vb t_vb=\nset laststatus=2\nset statusline=%F%m%r%h%w\ [POS=%l,%v][%p%%]\ %{strftime(\"%d/%m/%Y\ -\ %H:%M\")}\nset backspace=2\nsyntax on\nsyntax enable\nset enc=utf-8\nset fenc=utf-8\nset fencs=ucs-bom,utf-8,cp936,gb18030,gb2312,gbk,big5,euc-jp,euc-kr,shift-jis,latin1\nset termencoding=utf-8\n' > .vimrc
# 配置静态 IP
tip "CONFIG STATIC IP ADDRESS"
# 设置默认的 IP，这里只是给个格式示例，没有什么特别的含义
ipv4="192.168.0.3/24"
gateway="192.168.0.1"
dns="45.14.107.16,180.76.76.76,8.8.8.8"

# 配置静态 IP 地址，最好就是该虚拟机在DHCP下被分配的IP地址，避免IP冲突
read -p "set ipv4 address(default: ${ipv4}): " tmp
if [ -n "${tmp}" ]; then
    echo "use user input."
    ipv4=${tmp}
else
    echo "use default."
fi
echo "ipv4: ${ipv4}"
# 配置网关, 默认是 IP 地址的最后一段取 1，请和宿主机的保持一致
read -p "set gatway(default: ${gateway}): " tmp
if [ -n "${tmp}" ]; then
    echo "use user input."
    gateway=${tmp}
else
    echo "use default."
fi
echo "gateway: ${gateway}"
# 配置 DNS 服务器地址
read -p "set DNS servers, sep by comma(default: ${dns}): " tmp
if [ -n "${tmp}" ]; then
    echo "use user input."
    dns=${tmp}
else
    echo "use default."
fi
echo "dns: ${dns}"
# 选择网卡
echo "next to chose a network interface, you can view them by using command \`ip link show\` or \`ip a\`"
interfaces=($(ip link show | awk '/^[0-9]+:/ {print substr($2, 1, length($2)-1)}'))
# 使用 select 命令生成菜单
select network in "${interfaces[@]}"; do
    if [[ -n "${network}" ]]; then
        echo "use user's choice: ${network}"
        break
    else
        echo "invalid option."
    fi
done
# 输出选中的网卡名称
echo "network interface: ${network}"

# 修改关于网络的配置文件, 注意版本的对应, Ubuntu24.04中配置文件是 /etc/netplan/50-cloud-init.yaml
netplan_path="/etc/netplan/"
netplan_filename="50-cloud-init.yaml"
cd ${netplan_path}
# 首先备份原始的配置文件
sudo cp -f ${netplan_filename} "${netplan_filename}.backup"
# 暂时放开权限
sudo chmod 777 ${netplan_filename}
echo -e "network:\n  renderer: networkd\n  ethernets:\n    # your real network interface name\n    ${network}:\n      # close DHCP\n      dhcp4: false\n      dhcp6: false\n      # static ip and subnet\n      addresses: [${ipv4}]\n      # gateway\n      routes:\n        - to: default\n          via: ${gateway}\n      # DNS \n      nameservers:\n        addresses: [${dns}]\n  version: 2" > ${netplan_filename}
# 恢复原有的权限, 原有权限为 -rw------- 即 600
sudo chmod 600 ${netplan_filename}
# 根据配置文件中的注释说明, 这个配置文件中的配置会在重启之后被重置
# 如果不想被重置的话, 就在 /etc/cloud/cloud.cfg.d/99-disable-network-config.cfg 文件中写入内容 network: {config: disabled}
cfg_path="/etc/cloud/cloud.cfg.d/"
cfg_file="99-disable-network-config.cfg"
# 文件默认是不存在的
cd ${cfg_path}
# 先创建, 默认权限是 644
sudo touch ${cfg_file}
# 暂时放开权限
sudo chmod 777 ${cfg_file}
# 写入内容
echo "network: {config: disabled}" > ${cfg_file}
# 恢复权限
sudo chmod 644 ${cfg_file}
# 调整时区
profile_path="/etc/profile"
# 备份文件
sudo cp -f "${profile_path}" "${profile_path}.backup"
sudo chmod 777 "${profile_path}"
echo -e '\n#Time Zone config\nTZ=Asia/Shanghai\nexport TZ' >> "${profile_path}"
sudo chmod 644 "${profile_path}"
source "${profile_path}"
# 到此为止, 总共做了三个基本配置, 一个是配置用户 sudo 免密, 一个是配置静态 IP, 一个是调整时区
# 下面安装基本的环境软件等
tip "INSTALL COMMON SOFTWARE"
sudo apt install -y build-essential net-tools zip unzip
# 刷新网络，注意这一步不能在安装软件之前，否则网络连接状态是NAT且IP是静态IP，会连不上网
sudo netplan apply
echo -e "${GREEN}The basic system environment has been configured. Now, please change the virtual machine's network connection mode to bridged mode.${RESET}"
# 以上, Ubuntu 基本环境就配置好了
