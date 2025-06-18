# 多阶段构建 - 支持多平台
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

# 构建参数
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG BUILDTIME
ARG VERSION
ARG REVISION
ARG CACHE_BUSTER=default  # 新增缓存破坏参数

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码 - 使用缓存破坏参数
COPY . .

# 构建应用 - 支持交叉编译
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT#v} \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILDTIME} -X main.Revision=${REVISION}" \
    -o main .

# 最终镜像 - 使用多平台基础镜像
FROM --platform=$TARGETPLATFORM alpine:latest

# 镜像元数据
LABEL org.opencontainers.image.title="Antimg" \
      org.opencontainers.image.description="Advanced image watermark attack processing tool" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${REVISION}" \
      org.opencontainers.image.created="${BUILDTIME}" \
      org.opencontainers.image.source="https://github.com/Neurocoda/Antimg" \
      org.opencontainers.image.licenses="MIT"

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# 创建非root用户
RUN addgroup -g 1001 -S antimg && \
    adduser -u 1001 -S antimg -G antimg

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /app/main ./
COPY --from=builder /app/templates ./templates/
COPY --from=builder /app/static ./static/

# 创建必要目录并设置权限
RUN mkdir -p uploads logs && \
    chown -R antimg:antimg /app

# 切换到非root用户
USER antimg

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

# 设置环境变量
ENV GIN_MODE=release \
    PORT=8080 \
    TZ=UTC

# 运行应用
ENTRYPOINT ["./main"]
