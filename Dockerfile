# 第一阶段：构建阶段
FROM golang:1.21-alpine as builder

# 设置工作目录
WORKDIR /app

# 安装git，某些Go模块需要
RUN apk add --no-cache git

# 复制go mod文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# 第二阶段：运行阶段
FROM alpine:latest
WORKDIR /root/

# 如果你的应用需要运行时依赖，可以在这里安装
# RUN apk add --no-cache <some-packages>

# 从构建器阶段复制构建好的二进制文件到当前阶段
COPY --from=builder /app/main .
# 暴露应用程序监听的端口
EXPOSE 8080
# 运行应用
CMD ["./main"]