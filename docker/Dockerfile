# syntax = docker/dockerfile:1.3.0

FROM ubuntu:23.04

RUN useradd -m --uid 1000 --groups sudo ubuntu && \
  echo 'ubuntu:ubuntu' | chpasswd
RUN echo 'ubuntu ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
