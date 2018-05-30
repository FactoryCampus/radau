#!/bin/bash

if [ ! -z $DEBIAN_ARCH ]; then
    DEBIAN_ARCH=amd64
fi

mkdir -p build/$DEBIAN_ARCH

GOARCH=$DEBIAN_ARCH go build -o build/$DEBIAN_ARCH/wifi-login-backend

go get -v -u github.com/mh-cbon/go-bin-deb
if [ ! -z $TRAVIS_TAG ]; then
    go-bin-deb generate --version $TRAVIS_TAG -a $DEBIAN_ARCH
else
    go-bin-deb generate --version "0.0.0-${TRAVIS_COMMIT}" -a $DEBIAN_ARCH
fi
