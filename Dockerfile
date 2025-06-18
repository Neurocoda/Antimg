# 使用多阶段构建（兼容多平台）
ARG GO_VERSION=1.22-alpine

# 构建阶段（显式指定与构建平台一致的基础镜像加速编译）
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -installsuffix cgo -o main .

# 最终镜像（自动选择对应架构的 Alpine）
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
RUN mkdir -p uploads
EXPOSE 8080
ENV GIN_MODE=release PORT=8080
CMD ["./main"]
