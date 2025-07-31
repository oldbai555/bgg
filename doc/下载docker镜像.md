- sudo apt install skopeo
- sudo skopeo copy docker://docker.io/<镜像名称>:<标签> docker-archive:<存储路径>/<镜像名字>.tar:<标签>
- 示例：sudo skopeo copy docker://docker.io/mayanghua/instock:202410 docker-archive:/home/ubuntu/skope/instock.tar:202410


- go install gitee.com/extrame/dget/cmd/dget@latest
- ./dget mayanghua/instock:202410
- sudo docker load -i instock_202410-img.tar.gz
- 换源
```shell
vim /etc/docker/daemon.json
# 写入内容
{
    "registry-mirrors": [
        "https://docker.m.daocloud.io",
        "https://docker.imgdb.de",
        "https://docker-0.unsee.tech",
        "https://docker.hlmirror.com",
        "https://cjie.eu.org"
    ]
  }
# wq 保存
# 更新重启
sudo systemctl daemon-reload && sudo systemctl restart docker  
```

## docker 反向代理网站
- https://github.com/sky22333/hubproxy.git
