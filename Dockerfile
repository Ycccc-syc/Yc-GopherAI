# ===阶段1：编译=====
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

#CGO_ENABLED=0:禁用CGO，跳过了image识别模块，不需要gcc
RUN GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 go build -o gopherai .

# =====阶段2：运行 ======
FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/gopherai .
COPY --from=builder /app/config/config.toml ./config/config.toml

EXPOSE 9090

CMD ["./gopherai"]

