# ---- 使用轻量级 Debian 作为基础镜像 ----
FROM debian:bullseye-slim

# 设置工作目录
WORKDIR /app

# 安装 curl（用于下载二进制文件）和 ca-certificates
RUN apt-get update && \
    apt-get install -y curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 下载你的 Go 二进制文件（请替换成你的真实远程地址）
RUN curl -L -o crawler https://raw.githubusercontent.com/kawaiirei0/crawler-cfip/main/run && \
    chmod +x crawler

# 暴露 Gin 项目端口
EXPOSE 8080

# 启动二进制文件
CMD ["./crawler"]
