#!/bin/sh

# assumes you have already run:
# go get -u github.com/jteeuwen/go-bindata/...

go-bindata -o=assets/bindata.go --nocompress --nometadata --pkg=assets templates/... js/... json/... static/...
go build

