echo "Building SpineChain..."


set GOOS=linux
go build -o build/spinechain -ldflags "-X main.version=%VERSION%" .