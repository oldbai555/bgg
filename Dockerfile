# === 1. 镜像从源码阶段编译 ===
FROM ubuntu:18.04 AS build

# 安装必要的工具和库
RUN apt-get update && \
    apt-get install -y curl git build-essential

# 安装 Go 1.18
RUN curl -O https://storage.googleapis.com/golang/go1.18.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz && \
    rm -rf go1.18.linux-amd64.tar.gz

# 设置环境变量
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin

# 创建工作目录
RUN mkdir -p $GOPATH/src/app
WORKDIR $GOPATH/src/app

# 拷贝代码到工作目录中
COPY . .

# 配制go环境变量
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go env -w GO111MODULE=auto
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64
RUN go mod download
RUN go build -o lb lbserver/cmd/main.go

# === 2.镜像工作目录 ===
FROM ubuntu:18.04

# 安装必要的工具和库
RUN apt-get update && \
    apt-get install -y curl git build-essential

# 配置工作目录
WORKDIR /appruntime

# 拷贝
COPY --from=build /go/src/app/lb .
COPY --from=build /go/src/app/lbserver/resource/application.yaml .
ENTRYPOINT ["./lb"]

EXPOSE 8003
