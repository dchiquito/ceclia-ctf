#!/bin/sh

mkdir -p build/phase3

echo "Building Phase 3: Logging in..."

echo "\tCopying passwordRecovery file from Phase 4..."
cp build/phase4/passwordRecovery phase3/static/passwordRecovery

cd phase3

echo "\tgo getting go-bindata..."
go get -u github.com/jteeuwen/go-bindata/...

echo "\tUsing go-bindata to build assets package..."
go-bindata -o=assets/bindata.go --nocompress --nometadata --pkg=assets templates/... js/... json/... static/...

echo "\tgo building phase3..."
go build -o ../build/phase3/ceclia-ctf-go

