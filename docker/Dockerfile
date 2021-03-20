# syntax = docker/dockerfile:1.0-experimental

FROM ubuntu:20.04

RUN \
  --mount=type=cache,target=/var/cache/apt \
  --mount=type=cache,target=/var/lib/apt \
  apt update && \
  apt install -y tzdata && \
  apt install -y --fix-missing \
  curl \
  wget \
  git \
  python3 \
  python3-pip \
  libtiff5-dev \
  libjpeg8-dev \
  libopenjp2-7-dev \
  zlib1g-dev \
  libfreetype6-dev \
  liblcms2-dev \
  libwebp-dev \
  tcl8.6-dev \
  tk8.6-dev \
  python-tk \
  libharfbuzz-dev \
  libfribidi-dev \
  libssl-dev \
  zlib1g-dev \
  libbz2-dev \
  libreadline-dev \
  libsqlite3-dev \
  llvm \
  libncurses5-dev \
  xz-utils \
  tk-dev \
  libxml2-dev \
  libxmlsec1-dev \
  libffi-dev \
  liblzma-dev \
  gcc \
  make \
  build-essential \
  binutils \
  vim \
  emacs

RUN pip3 install flask

RUN wget https://github.com/longld/peda/archive/refs/tags/v1.2.tar.gz && \
  mkdir ~/peda && \
  tar xvzf v1.2.tar.gz -C ~/peda && \
  rm v1.2.tar.gz
RUN echo "source ~/peda/peda.py" >> ~/.gdbinit
