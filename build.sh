#!/bin/sh

PKG=github.com/bobisme/discoverable-api

docker run -v $PWD:/go/src/$PKG -w /go/src/$PKG discobuilder $@
