FROM golang:latest

WORKDIR /app

RUN go install \
    golang.org/x/tools/gopls@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    build-essential \
    cmake \
    unzip \
    gettext \
 #   python3 \
  #  python3-pip \
    #nodejs \
    npm \
    ripgrep \
    fd-find \
    && rm -rf /var/lib/apt/lists/*

# Installing nvim
RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-arm64.appimage && \
    chmod u+x nvim-linux-arm64.appimage && \
    ./nvim-linux-arm64.appimage --appimage-extract && \
    mv squashfs-root /usr/local/nvim && \
    ln -s /usr/local/nvim/AppRun /usr/local/bin/nvim && \
    rm nvim-linux-arm64.appimage

# Nvim Config Files
RUN git clone https://github.com/ilkerciblak/nvim-base-config.git ~/.config/nvim 

# lazy.nvim bootstrap 
RUN git clone --filter=blob:none \
    https://github.com/folke/lazy.nvim.git \
    --branch=stable \
    ~/.local/share/nvim/lazy/lazy.nvim

RUN npm install -g tree-sitter-cli

# Installing plugins in headless mode
RUN nvim --headless "+Lazy! sync" +qa

ENV GO111MODULE=on \ 
    GOPATH=/go \
    PATH=$PATH:/go/bin
