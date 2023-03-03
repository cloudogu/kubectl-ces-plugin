FROM golang:1.20

ENV YQ_CHECKSUM=71ef4141dbd9aec3f7fb45963b92460568d044245c945a7390831a5a470623f7 \
    YQ_VERSION=v4.31.2

RUN set -eux pipefail && \
  apt-get update && \
  apt-get install -y make zip tar wget && \
  wget https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 && \
  echo "${YQ_CHECKSUM} *yq_linux_amd64" | sha256sum -c - && \
  chmod +x yq_linux_amd64 && \
  mv yq_linux_amd64 /usr/bin/yq

