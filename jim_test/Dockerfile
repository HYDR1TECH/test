## 设置基础镜像
#FROM golang:1.19 AS build
## 复制go.mod和go.sum文件以下载依赖项（如果存在）
#WORKDIR /app/
#COPY . ./
## 将应用程序代码复制到容器中
#RUN cd /app && go build -o main main.go
#
FROM alpine:latest
WORKDIR /app
COPY ./main .
COPY conf/config.toml ./conf/
RUN chmod +x /app/main
EXPOSE 8000
# 设置容器启动命令
CMD ["./main"]
