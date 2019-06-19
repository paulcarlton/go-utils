# (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP
FROM staging-docker.artifactory.zing.hpelabs.net/panormos/go-build:0.0.13-g5e3627d50c34f1 as builder

ARG VERSION
WORKDIR /go/src/github.hpe.com/platform-core/utils
COPY . .
RUN make build

ENV TAG=$TAG \
  GIT_SHA=$GIT_SHA \
  BUILD_DATE=$BUILD_DATE \
  SRC_REPO=$SRC_REPO

LABEL TAG=$TAG \
  GIT_SHA=$GIT_SHA \
  BUILD_DATE=$BUILD_DATE \
  SRC_REPO=$SRC_REPO
