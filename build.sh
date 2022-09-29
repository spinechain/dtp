#!/bin/bash

VERSION=$(<VERSION)
echo "Building SpineChain v$VERSION..."

set GOOS=linux
go build -o package/linux/spinechain-$VERSION -ldflags "-X main.version=$VERSION" .

gh release delete-asset v$VERSION spinechain-$VERSION
gh release upload v$VERSION ./package/linux/spinechain-$VERSION#"spinechain-$VERSION - Ubuntu Package"