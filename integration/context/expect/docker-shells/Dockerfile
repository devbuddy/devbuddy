FROM ubuntu:focal

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get -y install --no-install-recommends ca-certificates zsh && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
