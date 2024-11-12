- sudo apt install skopeo
- sudo skopeo copy docker://docker.io/<镜像名称>:<标签> docker-archive:<存储路径>/<镜像名字>.tar:<标签>
- 示例：sudo skopeo copy docker://docker.io/mayanghua/instock:202410 docker-archive:/home/ubuntu/skope/instock.tar:202410


- go install gitee.com/extrame/dget/cmd/dget@latest
- ./dget mayanghua/instock:202410
- sudo docker load -i instock_202410-img.tar.gz
