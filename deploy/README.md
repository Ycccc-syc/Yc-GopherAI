# GopherAI 生产部署指南

本文档教你从零开始，把 GopherAI 部署到云服务器上，让公网用户能通过 HTTPS 访问。

## 目录

1. [整体架构](#1-整体架构)
2. [环境准备](#2-环境准备)
3. [数据库初始化](#3-数据库初始化)
4. [构建与部署](#4-构建与部署)
5. [Nginx 反向代理](#5-nginx-反向代理)
6. [HTTPS 证书](#6-https-证书)
7. [Systemd 服务管理](#7-systemd-服务管理)
8. [验证上线](#8-验证上线)
9. [日常维护](#9-日常维护)
10. [常见问题](#10-常见问题)

---

## 1. 整体架构

```text
用户浏览器
     │
     ▼
  HTTPS (443)         ← 加密传输，安全
     │
     ▼
  Nginx (反向代理)     ← 入口网关
     │
     ├── /api/*  ──▶  Go 后端 (127.0.0.1:9090)
     │                    ├── MySQL (数据库)
     │                    ├── Redis (缓存)
     │                    └── RabbitMQ (消息队列)
     │
     └── /*  ──▶  Vue 前端静态文件
```

**关键理解**:
- 用户只访问 Nginx 的 443 端口（HTTPS）
- Nginx 把 `/api/` 请求转发给 Go 后端，把其他请求当作前端页面返回
- Go 后端不直接暴露给公网，只监听 `127.0.0.1`（本地），安全隔离
- 你只需要在云服务器的**防火墙/安全组**开放 `443` (HTTPS) 和 `80` (HTTP 重定向) 端口

---

## 2. 环境准备

### 2.1 购买云服务器

推荐配置（最低）:
- CPU: 2 核
- 内存: 4GB
- 系统: Ubuntu 22.04 或 CentOS 7+
- 带宽: 5Mbps

云服务商: 阿里云、腾讯云、华为云、AWS Lightsail 等均可。

### 2.2 安全组/防火墙放行端口

在你的云服务商控制台，找到**安全组规则**，放行:

| 端口 | 协议 | 说明 |
|------|------|------|
| 22 | TCP | SSH 连接（默认已开）|
| 80 | TCP | HTTP（用于 Let's Encrypt 验证和重定向）|
| 443 | TCP | HTTPS（生产访问）|

### 2.3 SSH 登录服务器

```bash
ssh root@你的服务器IP
```

### 2.4 安装基础软件

```bash
# Ubuntu/Debian
apt update
apt install -y git golang-go nodejs npm mysql-server redis-server nginx certbot python3-certbot-nginx build-essential

# 安装 ONNX Runtime（图片识别需要）
# 如果不需要图片识别功能，可以跳过这步
wget https://github.com/microsoft/onnxruntime/releases/download/v1.15.1/onnxruntime-linux-x64-1.15.1.tgz
tar -xzf onnxruntime-linux-x64-1.15.1.tgz
sudo cp -r onnxruntime-linux-x64-1.15.1/lib/* /usr/local/lib/
sudo ldconfig
export CGO_CFLAGS="-I/path/to/onnxruntime/include"
export CGO_LDFLAGS="-L/path/to/onnxruntime/lib -lonnxruntime"

# 验证安装
go version    # 需要 1.21+
node --version
npm --version
mysql --version
nginx -v
redis-server --version
```

> **为什么需要这些？**
> - Go: 编译后端二进制
> - Node.js: 构建前端页面
> - MySQL: 存用户、会话、聊天记录
> - Redis: 缓存（验证码等）
> - Nginx: 反向代理 + 静态文件托管 + HTTPS
> - Certbot: 自动申请 Let's Encrypt 免费 SSL 证书

---

## 3. 数据库初始化

### 3.1 启动 MySQL

```bash
# 启动 MySQL
systemctl start mysql
systemctl enable mysql   # 开机自启

# 查看状态
systemctl status mysql
```

### 3.2 创建数据库和用户

```bash
# 登录 MySQL
mysql -u root -p

# 如果刚安装没有密码，直接回车
```

在 MySQL 命令行中执行:

```sql
-- 创建数据库（和我们 config/config.toml 配置一致）
CREATE DATABASE GopherAI CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建专门的应用用户（安全考虑，不要直接用 root）
CREATE USER 'gopherai'@'localhost' IDENTIFIED BY '你的密码';

-- 授权
GRANT ALL PRIVILEGES ON GopherAI.* TO 'gopherai'@'localhost';
FLUSH PRIVILEGES;

-- 退出
EXIT;
```

### 3.3 修改配置文件

用你刚创建的数据库用户和密码，更新 `config/config.toml`:

```toml
[mysqlConfig]
host = "127.0.0.1"
port = 3306
user = "gopherai"
password = "你的密码"
databaseName = "GopherAI"
charset = "utf8mb4"
```

> **注意**: Go 后端会自动建表（`AutoMigrate`），你不需要手动建表。

### 3.4 启动 Redis 和 RabbitMQ（如果需要）

```bash
systemctl start redis
systemctl enable redis

# RabbitMQ（如果配置了）
apt install -y rabbitmq-server
systemctl start rabbitmq-server
systemctl enable rabbitmq-server
```

---

## 4. 构建与部署

### 方式一：在服务器上直接构建（推荐新手）

```bash
# 1. 拉取代码
cd /opt
git clone https://github.com/你的仓库/GopherAI.git
cd GopherAI

# 2. 配置（修改数据库密码等）
vim config/config.toml

# 3. 一键构建后端 + 前端
make deploy
```

### 方式二：本地构建后上传

在你的开发机器上:

```bash
# 编译 Linux 二进制（Windows/Mac 都能编译出 Linux 可执行文件）
make build

# 构建前端
make frontend-build
```

然后把 `gopherai` 二进制文件和 `vue-frontend/dist/` 目录上传到服务器。

### 4.1 复制文件到正确位置

```bash
# 创建应用目录
mkdir -p /opt/gopherai

# 复制后端
cp gopherai /opt/gopherai/
cp -r config /opt/gopherai/

# 复制前端
mkdir -p /var/www/gopherai
cp -r vue-frontend/dist/* /var/www/gopherai/dist/
```

---

## 5. Nginx 反向代理

### 5.1 什么是反向代理？

想象 Nginx 是一个**前台接待员**:
- 用户来了 → 接待员判断他要什么
- 要 API 服务 → 转给 Go 后端（内部分机）
- 要网页 → 直接拿静态文件给他

这样 Go 后端不用直接面对公网，更安全。

### 5.2 部署配置

```bash
# 1. 复制 Nginx 配置（记得先修改域名）
cp deploy/nginx.conf /etc/nginx/sites-available/gopherai

# 2. 编辑配置文件，把 your-domain.com 替换成你的域名
vim /etc/nginx/sites-available/gopherai

# 3. 启用站点
ln -s /etc/nginx/sites-available/gopherai /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default   # 删除默认站点

# 4. 测试配置是否正确
nginx -t

# 5. 重新加载 Nginx
systemctl reload nginx
```

### 5.3 Nginx 配置详解

看 `deploy/nginx.conf`，关键部分:

```nginx
# API 反向代理 —— 把 /api/ 请求转发给 Go 后端
location /api/ {
    proxy_pass http://127.0.0.1:9090/api/v1/;
    # 这些是在传递真实客户端 IP 给后端
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

# 前端静态文件 —— Vue 打包后的页面
location / {
    try_files $uri $uri/ /index.html;
    # 这里 /index.html 是为了支持 Vue 路由
    # 比如用户直接访问 /chat，Nginx 找不到这个文件，就返回 index.html，让 Vue 路由处理
}
```

---

## 6. HTTPS 证书

### 6.1 为什么需要 HTTPS？

- HTTP 传输是明文的，密码、聊天内容容易被截获
- HTTPS 加密传输，安全
- 现代浏览器会标记 HTTP 网站为"不安全"
- 很多 AI API 要求请求来源是 HTTPS

### 6.2 用 Let's Encrypt 申请免费证书

```bash
# 确保域名已经解析到你的服务器 IP
# 在你的 DNS 管理后台，添加 A 记录指向服务器公网 IP

# 申请证书（Certbot 会自动修改 Nginx 配置）
certbot --nginx -d your-domain.com

# 按照提示输入邮箱，同意协议即可

# 测试自动续期（Let's Encrypt 证书 90 天有效）
certbot renew --dry-run
```

> **证书会自动续期**，系统会每天检查一次，到期前 30 天自动续。

### 6.3 验证 HTTPS

打开浏览器访问 `https://your-domain.com`，能看到小锁图标就成功了。

---

## 7. Systemd 服务管理

### 7.1 为什么用 Systemd？

你的 Go 后端是一个需要在后台持续运行的程序。Systemd 负责:
- 服务器开机时自动启动 Go 后端
- 程序崩溃时自动重启
- 记录日志供你排查问题

### 7.2 配置服务

```bash
# 复制服务配置文件
cp deploy/gopherai.service /etc/systemd/system/

# 重新加载 Systemd
systemctl daemon-reload

# 启用开机自启
systemctl enable gopherai

# 启动服务
systemctl start gopherai

# 查看状态
systemctl status gopherai
```

### 7.3 常用管理命令

```bash
# 查看服务状态
systemctl status gopherai

# 查看实时日志
journalctl -u gopherai -f

# 重启服务（更新二进制后）
systemctl restart gopherai

# 停止服务
systemctl stop gopherai

# 查看最近 50 行日志
journalctl -u gopherai -n 50 --no-pager
```

---

## 8. 验证上线

### 8.1 测试健康检查

```bash
curl https://your-domain.com/healthz
```

期望返回:
```json
{"status":"ok","service":"GopherAI","mysql":true,"redis":true}
```

### 8.2 测试注册登录

```bash
curl -X POST https://your-domain.com/api/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","email":"test@example.com"}'
```

### 8.3 测试前端

在浏览器打开 `https://your-domain.com`，应该能看到登录页面。
登录后能正常使用 AI 对话功能。

---

## 9. 日常维护

### 9.1 更新版本

```bash
cd /opt/GopherAI

# 拉取最新代码
git pull

# 重新构建
make deploy

# 复制新二进制
cp gopherai /opt/gopherai/
cp -r vue-frontend/dist/* /var/www/gopherai/dist/

# 重启服务
systemctl restart gopherai
```

### 9.2 查看日志

```bash
# 后端日志
journalctl -u gopherai -f

# Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

### 9.3 备份数据库

```bash
# 备份
mysqldump -u gopherai -p GopherAI > backup-$(date +%Y%m%d).sql

# 恢复
mysql -u gopherai -p GopherAI < backup-20250101.sql
```

### 9.4 监控服务器

```bash
# 查看系统资源
htop
# 或
top

# 查看磁盘
df -h

# 查看内存
free -h
```

---

## 10. 常见问题

### Q: 502 Bad Gateway
Nginx 连不上 Go 后端。原因可能是:
- Go 后端没启动: `systemctl start gopherai`
- Go 后端崩溃: `journalctl -u gopherai -f` 查看错误
- 端口不匹配: 检查 Nginx 配置和 config.toml 端口是否一致

### Q: 403 Forbidden
- 前端 dist 目录权限不对: `chown -R www-data:www-data /var/www/gopherai`
- SELinux 问题: `setenforce 0` 临时测试

### Q: 前端空白页/路由错误
- Nginx 没有配置 `try_files $uri $uri/ /index.html;`
- Vue 是 SPA，所有非 API 请求都要返回 index.html

### Q: 数据库连接失败
- MySQL 没启动: `systemctl start mysql`
- 密码错误: 检查 config.toml
- 数据库不存在: `CREATE DATABASE GopherAI;`

### Q: HTTPS 证书过期
- Let's Encrypt 自动续期失败: `certbot renew`
- 检查域名 DNS 是否还指向你的服务器
