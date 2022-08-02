#!/bin/sh

build_push(){
  docker buildx build  --platform ${ARCHS} -t ${REGISTRY}/${NAME}:latest   --push .
}

helm_build_push(){
  FN=${NAME}-${VER}.tgz
  # rm ${FN}
  helm package ./install --version ${VER}
  curl --data-binary "@${FN}" http://helm.alexstorm.solenopsys.org/api/charts
}

REGISTRY=registry.alexstorm.solenopsys.org
NAME=alexstorm-public-api
ARCHS="linux/amd64,linux/arm64"
VER=0.1.0


helm_build_push
build_push





