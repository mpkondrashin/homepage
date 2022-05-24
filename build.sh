#!/bin/bash
#
# HomePage (c) 2022 by Michael Kondrashin
#
# build.sh - script to build for all platforms
#

set -e

NAME=homepage

build() {
    echo "Build ${NAME} for ${1} ${2}"
    GOOS=$1 GOARCH=$2 go build -o ${NAME}-${1}-${2}${3}
}

build darwin arm64
build darwin amd64
build linux amd64
build windows amd64 .exe

