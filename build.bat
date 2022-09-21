echo "Building SpineChain..."

set GOOS=windows
go build -o package/spinechain.exe -ldflags "-X main.version=%VERSION%" .
