# 设置基础镜像
FROM golang:1.19
# 在容器中创建一个工作目录
WORKDIR /app
# 复制go.mod和go.sum文件以下载依赖项（如果存在）
COPY . .
# 将应用程序代码复制到容器中
RUN go build -o main
COPY conf/config.toml ./conf
RUN chmod +x main
EXPOSE 8000
# 设置容器启动命令
ENTRYPOINT ["./main"]
