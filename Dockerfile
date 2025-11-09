# ---- 使用轻量级 Debian 作为基础镜像 ----
FROM debian:bullseye-slim

# 设置工作目录
WORKDIR /app

# 安装 Chromium + unzip + 依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        chromium \
        unzip \
        ca-certificates \
        curl \
        fonts-liberation \
        libnss3 \
        libxss1 \
        libatk-bridge2.0-0 \
        libgtk-3-0 \
        libasound2 \
        libxshmfence1 && \
    rm -rf /var/lib/apt/lists/*

# 下载编译好的二进制 zip 文件并解压
RUN curl -L -o crawler.zip https://store.ilin.eu.org/crawler/crawler.zip && \
    unzip crawler.zip && \
    rm crawler.zip && \
    chmod +x crawler

# 暴露 Gin 项目端口
EXPOSE 8080

# 设置 Chromium 执行路径环境变量（chromedp 能自动识别）
ENV CHROME_PATH=/usr/bin/chromium

# 启动二进制文件
CMD ["./crawler"]
