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

RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension cnbcool.cnb-welcome \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension orta.vscode-jest \
  && code-server --install-extension dbaeumer.vscode-eslint \
  && code-server --install-extension waderyan.gitblame \
  && code-server --install-extension mhutchie.git-graph \
  && code-server --install-extension donjayamanne.githistory

ENV LANG=C.UTF-8 \
    LANGUAGE=C.UTF-8 \
    GOPATH=/go \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# acli
RUN echo -e "\n认证完成后，请执行：acli rovodev run"
CMD [".ide/acli.sh"]