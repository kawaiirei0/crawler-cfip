# 克隆

```bash
git clone https://github.com/kawaiirei0/crawler-cfip.git
```

# 下载gx

```
https://github.com/kawaiirei0/gx
```

## gx下载地址

```
https://github.com/kawaiirei0/gx/releases/download/v1.0.2/gx-windows-amd64.exe
```

# 打包crawler

```
gx cross-build --os linux --arch amd64 -o dist/my-crawler
```

# 或者使用脚本，一键部署到docker

```bash
curl -sSL https://raw.githubusercontent.com/kawaiirei0/crawler-cfip/main/run.sh | bash -s my-crawler crawler-container 9090
```