# Copyright 2021 CloudWeGo Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# === 1. 镜像从源码阶段编译 ===
FROM golang:1.18.9-alpine3.17 AS build

# 换alpine镜像源
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 更新软件包和安装git
RUN apk update && apk add git

# 配置工作目录
WORKDIR /source

# 拷贝
COPY . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go env -w GO111MODULE=on
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64
RUN go build -o lbapp-blog ./lbblog.proto/cmd/main.go

# === 2.镜像工作目录 ===
FROM golang:1.18.9-alpine3.17
# 配置工作目录
WORKDIR /appruntime
# 拷贝
COPY --from=build /source/lbapp-blog .
ENTRYPOINT ["./lbapp-blog"]

EXPOSE 8003