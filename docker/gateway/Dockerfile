# === 1. 镜像从源码阶段编译 ===
FROM golang:1.18-alpine AS build

# 设置环境变量
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin

# 创建工作目录
RUN mkdir -p $GOPATH/src/app
WORKDIR $GOPATH/src/app

# 拷贝代码到工作目录中
COPY .. .

RUN go env

# 配制go环境变量
RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn,direct
RUN go env -w GO111MODULE=auto
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64
RUN go mod tidy
RUN go build -o gateway service/gateway/main.go

# === 2.镜像工作目录 ===
FROM golang:1.18-alpine

# 配置工作目录
WORKDIR /appruntime

# 拷贝
COPY --from=build /go/src/app/gateway .
COPY --from=build /go/src/app/service/gateway/internal/conf/application.yaml /etc/work/
COPY --from=build /go/src/app/pkg/grpc_tool/etcd.json /etc/work/
ENTRYPOINT ["./gateway"]

EXPOSE 20000
