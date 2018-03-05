#!/bin/sh

set -x

TARGET='anycast-operator'

GOPATH='/build'
export GOPATH

if [[ ! -d "${GOPATH}" ]]; then
    echo "[E] ${GOPATH} not found, exiting ..."
    exit 1
fi

go get -v github.com/r3boot/${TARGET}/...
cd ${GOPATH}/src/github.com/r3boot/${TARGET}
go build -v -o ${GOPATH}/${TARGET} cmd/${TARGET}/main.go
strip ${GOPATH}/${TARGET}
ls -lah ${GOPATH}/${TARGET}
chown -R 1000:1000 ${GOPATH}/*