#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -v -tags 'netgo osusergo static_build debug' -ldflags '-s -w' -gcflags=all='-l=4 -B -C'
