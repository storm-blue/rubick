#!/bin/bash

build_image() {
  IMAGE_TAG=$1
  DIR_PATH=editor

  IMAGE_NAME=whale-registry.meetwhale.com/whale-dev/rubick:$IMAGE_TAG

  docker -D build --build-arg DIR="${DIR_PATH}" -f ./Dockerfile -t "${IMAGE_NAME}" ../../..

  if [ $? -ne 0 ]; then
    echo "failed"
    exit 1
  else
    echo "succeed"
  fi

  docker push "${IMAGE_NAME}"
}

IMAGE_TAG=$1

if [ -z "$IMAGE_TAG" ]; then
  IMAGE_TAG=latest
fi

build_image "$IMAGE_TAG"
