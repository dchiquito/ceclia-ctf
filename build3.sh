#!/bin/sh

mkdir -p build/phase3

echo "Building Phase 3: Logging in..."

echo "\tCopying passwordRecovery file from Phase 4..."
cp build/phase4/passwordRecovery phase3/static/passwordRecovery

cd phase3

echo "\tgo getting go dependencies..."
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/elazarl/go-bindata-assetfs/...
go get -u github.com/dgrijalva/jwt-go/...

echo "\tUsing go-bindata to package all resources..."
go-bindata -o=assets/templates.go --nocompress --nometadata --pkg=assets json/... static/... templates/...

#echo "\tUsing go-bindata-assets to package all resources..."
#go-bindata-assetfs -o=assets/bindata.go --nocompress --nometadata --pkg=assets json/... static/... templates/...
#rm bindata.go # this is created by go-bindata-assetfs for some reason

echo "\tgo building phase3..."
go build -o ../build/phase3/ceclia-ctf-go

