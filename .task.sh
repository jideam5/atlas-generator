#!/bin/sh
PROJECT=atlas-generator
PbDIR=src/pb

build() {
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -ldflags '-w -s' \
	-o build/${PROJECT} src/main.go
}

run() {
	cd src && go run main.go ../tmp/packParam.json
}

proto() {
	cd ${PbDIR} && protoc --go_out=. *.proto
	ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
}

mod() {
	go mod tidy
}

case "$1" in
	build)
	$1
	;;
	run)
	$1
	;;
	proto)
	$1
	;;
	mod)
	$1
	;;
	*)
	echo "Usage: $0 {build|run|proto}"
	exit 1
	;;
esac
