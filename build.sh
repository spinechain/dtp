echo "Building SpineChain..."


set GOOS=linux
go build -o package/linux/spinechain -ldflags "-X main.version=%VERSION%" .