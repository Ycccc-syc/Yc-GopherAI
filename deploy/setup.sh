#!/bin/bash
# GopherAI 一键部署脚本
# 使用方法: bash deploy/setup.sh your-domain.com
# 在全新的 Ubuntu 云服务器上运行

set -e  # 任何一步出错就停止

DOMAIN=$1
if [ -z "$DOMAIN" ]; then
    echo "用法: bash deploy/setup.sh your-domain.com"
    exit 1
fi

echo "========================================"
echo " GopherAI 部署脚本"
echo " 域名: $DOMAIN"
echo "========================================"

# ---- 1. 安装依赖 ----
echo "[1/7] 安装依赖..."
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx \
                    golang-go nodejs npm \
                    mysql-server redis-server \
                    git build-essential

# 安装 ONNX Runtime（图片识别需要）
# 如果不需要图片识别功能，可以跳过这段
echo "[1.1/7] 安装 ONNX Runtime..."
wget -q https://github.com/microsoft/onnxruntime/releases/download/v1.15.1/onnxruntime-linux-x64-1.15.1.tgz
tar -xzf onnxruntime-linux-x64-1.15.1.tgz
sudo cp onnxruntime-linux-x64-1.15.1/lib/libonnxruntime* /usr/local/lib/
sudo ldconfig
rm -rf onnxruntime-linux-x64-1.15.1*

# ---- 2. 创建目录 ----
echo "[2/7] 创建目录..."
sudo mkdir -p /opt/gopherai
sudo mkdir -p /var/www/gopherai

# ---- 3. 配置数据库 ----
echo "[2.5/7] 配置数据库..."
sudo systemctl start mysql
sudo mysql -e "CREATE DATABASE IF NOT EXISTS GopherAI CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# ---- 4. 构建后端 ----
echo "[3/7] 构建后端..."
export CGO_ENABLED=1
export CGO_CFLAGS="-I$HOME/onnxruntime/include"
export CGO_LDFLAGS="-L/usr/local/lib -lonnxruntime"
go build -ldflags="-s -w" -o gopherai .
sudo cp gopherai /opt/gopherai/
sudo cp -r config /opt/gopherai/

# ---- 5. 构建前端 ----
echo "[4/7] 构建前端..."
cd vue-frontend
npm install
npm run build
sudo mkdir -p /var/www/gopherai/dist
sudo cp -r dist/* /var/www/gopherai/dist/
cd ..

# ---- 6. 配置 Systemd 服务 ----
echo "[5/7] 配置 Systemd 服务..."
sudo cp deploy/gopherai.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable gopherai

# ---- 7. 配置 Nginx ----
echo "[6/7] 配置 Nginx..."
sed "s/your-domain.com/$DOMAIN/g" deploy/nginx.conf | sudo tee /etc/nginx/sites-available/gopherai > /dev/null
sudo ln -sf /etc/nginx/sites-available/gopherai /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx

# ---- 8. HTTPS 证书（Let's Encrypt） ----
echo "[7/7] 申请 HTTPS 证书..."
sudo certbot --nginx -d $DOMAIN --non-interactive --agree-tos -m admin@$DOMAIN || echo "证书申请跳过（可能已存在或需要交互）"

# ---- 启动服务 ----
echo ""
echo "========================================"
echo " 部署完成！启动服务..."
echo "========================================"
sudo systemctl start gopherai
sudo systemctl restart nginx

echo ""
echo "服务状态:"
sudo systemctl status gopherai --no-pager
echo ""
echo "访问地址: https://$DOMAIN"
echo "健康检查: https://$DOMAIN/healthz"
echo "查看后端日志: sudo journalctl -u gopherai -f"
