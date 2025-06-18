# .ide/Dockerfile
# 开发环境
FROM golang:1.22-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    build-essential \
    libjpeg62-turbo-dev \
    libpng-dev \
    ca-certificates \
    wget \
    vim \
    unzip \
    openssh-server \
    fonts-wqy-microhei \
    && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8 \
    LANGUAGE=C.UTF-8 \
    GOPATH=/go \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

