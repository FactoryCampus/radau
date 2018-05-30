#!/bin/bash

docker build --pull --cache-from "$IMAGE_NAME" --tag "$IMAGE_NAME" .

echo "$REGISTRY_PASS" | docker login -u "$REGISTRY_USER" --password-stdin

docker tag "$IMAGE_NAME" "${IMAGE_NAME}:${TRAVIS_COMMIT}"

if [ ! -z $TRAVIS_TAG ]; then
    docker tag "$IMAGE_NAME" "${IMAGE_NAME}:${TRAVIS_TAG}"
fi

if [ ! -z $TRAVIS_BRANCH ]; then
    docker tag "$IMAGE_NAME" "${IMAGE_NAME}:git-${TRAVIS_BRANCH}"
fi

if [ "$TRAVIS_BRANCH" = "develop" ]; then
    docker tag "$IMAGE_NAME" "${IMAGE_NAME}:latest"
fi
