#!/bin/bash
VERSION=$(<VERSION)
echo "$value"

echo "Building SpineChain v$VERSION..."

set GOOS=linux
go build -o package/linux/spinechain-$VERSION -ldflags "-X main.version=$VERSION" .