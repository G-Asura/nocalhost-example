#! /bin/sh
export GOPROXY=https://goproxy.cn
dlv --headless --log --listen :8181 --api-version 2 --accept-multiclient debug main.go