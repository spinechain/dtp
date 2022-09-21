echo "Building SpineChain..."

set VERSION=0.0.1

set GOOS=windows
go build -o package/windows/spinechain-%VERSION%.exe -ldflags "-X main.version=%VERSION%" .

set GOOS=linux

go build -o package/linux/spinechain -ldflags "-X main.version=%VERSION%" .