# 基础镜像
FROM node:14.17.0 as build

# 设置工作目录
WORKDIR /app

# 拷贝文件
COPY package*.json ./

# 安装依赖
RUN npm install

# 拷贝项目文件
COPY . .

# 构建项目
RUN npm run build

# 使用 nginx 作为 Web 服务器
FROM nginx:1.19.10-alpine

# 拷贝配置文件
COPY nginx/frontend.conf /etc/nginx/conf.d/default.conf

# 拷贝构建好的项目文件到 nginx
COPY --from=build /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
