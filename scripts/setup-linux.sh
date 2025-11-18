#!/bin/bash

set -e

dpkg --add-architecture arm64

apt-get update 
apt-get install -y \
    libwebp-dev \
    libwebp-dev:arm64 \
    crossbuild-essential-arm64