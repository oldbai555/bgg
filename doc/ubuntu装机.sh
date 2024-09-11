# 装机配制 ubuntu
# 1. 修改镜像源, https://www.runoob.com/linux/linux-comm-apt.html
# 2. 安装 brew
# - rm Homebrew.sh ; wget https://gitee.com/cunkai/HomebrewCN/raw/master/Homebrew.sh ; bash Homebrew.sh
# 3. 安装 golang, curl -O https://studygolang.com/dl/golang/go1.20.6.linux-amd64.tar.gz , https://www.runoob.com/go/go-environment.html
# 4. 安装 clash for linux, https://github.com/Dreamacro/clash/releases/tag/v1.17.0, or https://500ml.name
# 5. 安装 goland, 激活码 https://3.jetbra.in/

# =================================================================
# 其他一些工具安装
# rz sz
apt install lrzsz
# curl
apt install curl
# vim
apt install vim
# docker & docker-compose
curl -fsSL https://test.docker.com -o test-docker.sh && sh test-docker.sh
apt install docker-compose
#==================================================================

# k3s
#curl -sfL https://get.k3s.io | sh -
#mkdir -p $HOME/.kube
#sudo cp /etc/rancher/k3s/k3s.* $HOME/.kube/config
#sudo chown $(id -u):$(id -g) $HOME/.kube/config

# curl –sfL \
#     https://rancher-mirror.oss-cn-beijing.aliyuncs.com/k3s/k3s-install.sh | \
#     INSTALL_K3S_MIRROR=cn sh -s - \
#     --system-default-registry "registry.cn-hangzhou.aliyuncs.com" \
#     --write-kubeconfig ~/.kube/config \
#     --write-kubeconfig-mode 666 \
#     --disable traefik
#
# systemctl status k3s.service
# kubectl get pods -n kube-system

# =================================================================
# supervisor 安装
# apt-get install supervisor

# =================================================================
# nginx 安装
# https://www.runoob.com/linux/nginx-install-setup.html

# =================================================================
# etcd 安装
# https://www.runoob.com/linux/nginx-install-setup.html
