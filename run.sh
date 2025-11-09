#!/bin/bash

set -e

# ---------------- 配置 ----------------
REMOTE_DOCKERFILE_URL="https://raw.githubusercontent.com/kawaiirei0/crawler-cfip/main/Dockerfile"  # 替换成你的 Dockerfile 真实地址
IMAGE_NAME="${1:-my-crawler}"       # 镜像名字，可通过第一个参数自定义
CONTAINER_NAME="${2:-my-crawler}"   # 容器名字，可通过第二个参数自定义
HOST_PORT="${3:-8080}"              # 宿主机端口，可通过第三个参数自定义
CONTAINER_PORT=8080                  # 容器内部端口，固定为 Dockerfile 里暴露的端口
# --------------------------------------

echo "===== 构建 Docker 镜像 ====="
docker build -t "$IMAGE_NAME" "$REMOTE_DOCKERFILE_URL"

echo "===== 检查并删除已有同名容器 ====="
EXISTING_CONTAINER=$(docker ps -aq -f name="$CONTAINER_NAME")
if [ "$EXISTING_CONTAINER" ]; then
    echo "停止并删除已有容器: $CONTAINER_NAME"
    docker stop "$CONTAINER_NAME"
    docker rm "$CONTAINER_NAME"
fi

echo "===== 启动容器 ====="
docker run -d \
    --name "$CONTAINER_NAME" \
    -p "$HOST_PORT":"$CONTAINER_PORT" \
    --restart=always \
    "$IMAGE_NAME"

echo "===== 完成 ====="
echo "容器 '$CONTAINER_NAME' 已启动，宿主机端口 $HOST_PORT 映射到容器端口 $CONTAINER_PORT"
echo "容器设置了自动重启策略，如果挂掉会自动重启"
